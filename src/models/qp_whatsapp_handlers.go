package models

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	whatsapp "github.com/nocodeleaks/quepasa/whatsapp"
	log "github.com/sirupsen/logrus"
)

// Serviço que controla os servidores / bots individuais do whatsapp
type QPWhatsappHandlers struct {
	server *QpWhatsappServer

	messages     map[string]whatsapp.WhatsappMessage
	sync         *sync.Mutex // Objeto de sinaleiro para evitar chamadas simultâneas a este objeto
	syncRegister *sync.Mutex

	// Appended events handler
	aeh []QpWebhookHandlerInterface
}

// get default log entry, never nil
func (source *QPWhatsappHandlers) GetLogger() *log.Entry {
	if source.server != nil {
		return source.server.GetLogger()
	}

	logger := log.StandardLogger()
	logger.SetLevel(log.ErrorLevel)

	return logger.WithContext(context.Background())
}

func (source QPWhatsappHandlers) HandleGroups() bool {
	global := whatsapp.Options

	var local whatsapp.WhatsappBoolean
	if source.server != nil {
		local = source.server.Groups
	}
	return global.HandleGroups(local)
}

func (source QPWhatsappHandlers) HandleBroadcasts() bool {
	global := whatsapp.Options

	var local whatsapp.WhatsappBoolean
	if source.server != nil {
		local = source.server.Broadcasts
	}
	return global.HandleBroadcasts(local)
}

//#region EVENTS FROM WHATSAPP SERVICE

// Process messages received from whatsapp service
func (source *QPWhatsappHandlers) Message(msg *whatsapp.WhatsappMessage) {

	// should skip groups ?
	if !source.HandleGroups() && msg.FromGroup() {
		return
	}

	// should skip broadcast ?
	if !source.HandleBroadcasts() && msg.FromBroadcast() {
		return
	}

	// messages sended with chat title
	if len(msg.Chat.Title) == 0 {
		msg.Chat.Title = source.server.GetChatTitle(msg.Chat.Id)
	}

	if len(msg.InReply) > 0 {
		cached, err := source.GetMessage(msg.InReply)
		if err == nil {
			maxlength := ENV.SynopsisLength() - 4
			if uint64(len(cached.Text)) > maxlength {
				msg.Synopsis = cached.Text[0:maxlength] + " ..."
			} else {
				msg.Synopsis = cached.Text
			}
		}
	}

	logger := source.GetLogger()
	logger.Debugf("appending to cache, received|sended from another app, id: %s, chatid: %s", msg.Id, msg.Chat.Id)
	source.appendMsgToCache(msg)
}

// does not cache msg, only update status and webhook dispatch
func (source *QPWhatsappHandlers) Receipt(msg *whatsapp.WhatsappMessage) {
	ids := strings.Split(msg.Text, ",")
	for _, element := range ids {
		cached, err := source.GetMessage(element)
		if err == nil {
			logger := source.GetLogger()

			// update status
			logger.Tracef("msg recebida/(enviada por outro meio) em models: %s", cached.Id)
		}
	}

	// triggering external publishers
	source.Trigger(msg)
}

/*
<summary>

	Event on:
		* User Logged Out from whatsapp app
		* Maximum numbers of devices reached
		* Banned
		* Token Expired

</summary>
*/
func (source *QPWhatsappHandlers) LoggedOut(reason string) {

	// one step at a time
	if source.server != nil {

		msg := "logged out !"
		if len(reason) > 0 {
			msg += " reason: " + reason
		}

		logger := source.GetLogger()
		logger.Warn(msg)

		// marking unverified and wait for more analyses
		source.server.MarkVerified(false)
	}
}

/*
<summary>

	Event on:
		* When connected to whatsapp servers and authenticated

</summary>
*/
func (source *QPWhatsappHandlers) OnConnected() {

	// one step at a time
	if source.server != nil {

		// marking unverified and wait for more analyses
		err := source.server.MarkVerified(true)
		if err != nil {
			logger := source.server.GetLogger()
			logger.Errorf("error on mark verified after connected: %s", err.Error())
		}
	}
}

/*
<summary>

	Event on:
		* When connected to whatsapp servers and authenticated

</summary>
*/
func (source *QPWhatsappHandlers) OnDisconnected() {

}

//#endregion
//region MESSAGE CONTROL REGION HANDLE A LOCK

// Salva em cache e inicia gatilhos assíncronos
func (handler *QPWhatsappHandlers) appendMsgToCache(msg *whatsapp.WhatsappMessage) {

	handler.sync.Lock() // await for avoid simultaneous calls

	normalizedId := msg.Id
	normalizedId = strings.ToUpper(normalizedId) // ensure that is an uppercase string before save

	// saving on local normalized cache, do not affect remote msgs
	handler.messages[normalizedId] = *msg

	handler.sync.Unlock()

	// should cleanup old messages ?
	length := ENV.CacheLength()
	handler.CleanUp(length)

	// continue to external dispatchers
	handler.Trigger(msg)
}

// Remove old ones, by timestamp, until a maximum length
func (handler *QPWhatsappHandlers) CleanUp(maximum uint64) {
	max := int(maximum)
	if max > 0 {
		handler.sync.Lock() // await for avoid simultaneous calls

		length := len(handler.messages)
		if length > max {
			messages := []whatsapp.WhatsappMessage{}
			for _, item := range handler.messages {
				messages = append(messages, item)
			}

			sort.Sort(whatsapp.WhatsappOrderedMessages(messages))

			amount := length - max
			for i := 0; i < amount; i++ {
				key := strings.ToUpper(messages[i].Id)
				delete(handler.messages, key)
			}
		}

		handler.sync.Unlock()
	}
}

func (handler *QPWhatsappHandlers) GetMessages(timestamp time.Time) (messages []whatsapp.WhatsappMessage) {
	handler.sync.Lock() // await for avoid simultaneous calls

	for _, item := range handler.messages {
		if item.Timestamp.After(timestamp) {
			messages = append(messages, item)
		}
	}

	handler.sync.Unlock()
	return
}

// Returns the first in time message stored in cache, used for resync history with message services like whatsapp
func (handler *QPWhatsappHandlers) GetLeadingMessage() (message *whatsapp.WhatsappMessage) {
	handler.sync.Lock() // await for avoid simultaneous calls

	now := time.Now()
	for _, item := range handler.messages {
		if !item.Timestamp.IsZero() && item.Timestamp.Before(now) {
			now = item.Timestamp
			message = &item
		}
	}

	handler.sync.Unlock()
	return
}

func (handler *QPWhatsappHandlers) GetMessagesByPrefix(id string) (messages []whatsapp.WhatsappMessage) {
	handler.sync.Lock() // await for avoid simultaneous calls

	for _, item := range handler.messages {
		if strings.HasPrefix(item.Id, id) {
			messages = append(messages, item)
		}
	}

	handler.sync.Unlock()
	return
}

// Get a single message if exists
func (handler *QPWhatsappHandlers) GetMessage(id string) (msg whatsapp.WhatsappMessage, err error) {
	handler.sync.Lock() // await for avoid simultaneous calls

	normalizedId := id
	normalizedId = strings.ToUpper(normalizedId) // ensure that is an uppercase string before save

	// getting from local normalized cache, do not affect remote msgs
	msg, ok := handler.messages[normalizedId]
	if !ok {
		err = fmt.Errorf("message not present on handlers (cache) id: %s", normalizedId)
	}

	handler.sync.Unlock()
	return msg, err
}

// endregion
// region EVENT HANDLER TO INTERNAL USE, GENERALLY TO WEBHOOK

// sends the message throw external publishers
func (source *QPWhatsappHandlers) Trigger(payload *whatsapp.WhatsappMessage) {
	if source != nil {
		if source.server != nil {
			go SignalRHub.Dispatch(source.server.Token, payload)
		}

		for _, handler := range source.aeh {
			go handler.HandleWebHook(payload)
		}
	}
}

// Register an event handler that triggers on a new message received on cache
func (handler *QPWhatsappHandlers) Register(evt QpWebhookHandlerInterface) {
	handler.sync.Lock() // await for avoid simultaneous calls

	if !handler.IsRegistered(evt) {
		handler.aeh = append(handler.aeh, evt)
	}

	handler.sync.Unlock()
}

// Removes an specific event handler
func (handler *QPWhatsappHandlers) UnRegister(evt QpWebhookHandlerInterface) {
	handler.sync.Lock() // await for avoid simultaneous calls

	newHandlers := []QpWebhookHandlerInterface{}
	for _, v := range handler.aeh {
		if v != evt {
			newHandlers = append(handler.aeh, evt)
		}
	}

	// updating
	handler.aeh = newHandlers

	handler.sync.Unlock()
}

// Removes an specific event handler
func (handler *QPWhatsappHandlers) Clear() {
	handler.sync.Lock() // await for avoid simultaneous calls

	// updating
	handler.aeh = nil

	handler.sync.Unlock()
}

// Indicates that has any event handler registered
func (handler *QPWhatsappHandlers) IsAttached() bool {
	return len(handler.aeh) > 0
}

// Indicates that if an specific handler is registered
func (handler *QPWhatsappHandlers) IsRegistered(evt interface{}) bool {
	for _, v := range handler.aeh {
		if v == evt {
			return true
		}
	}

	return false
}

//endregion

func (handler *QPWhatsappHandlers) GetTotal() int {
	return len(handler.messages)
}

func (source *QPWhatsappHandlers) IsInterfaceNil() bool {
	return nil == source
}

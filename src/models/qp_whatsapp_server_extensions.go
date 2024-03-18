package models

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"

	whatsapp "github.com/nocodeleaks/quepasa/whatsapp"
)

// handle message deliver to individual webhook distribution
func PostToWebHookFromServer(server *QpWhatsappServer, message *whatsapp.WhatsappMessage) (err error) {
	if server == nil {
		err = fmt.Errorf("server nil")
		return err
	}

	wid := server.GetWId()

	// ignoring ssl issues
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	logentry := server.GetLogger()
	for _, element := range server.Webhooks {
		sublogentry := logentry.WithField("url", element.Url)

		if message.Id == "readreceipt" && element.IsSetReadReceipts() && !element.ReadReceipts.Boolean() {
			sublogentry.Debugf("ignoring read receipt message: %s", message.Text)
			continue
		}

		if message.FromGroup() && element.IsSetGroups() && !element.Groups.Boolean() {
			sublogentry.Debugf("ignoring group message: %s", message.Id)
			continue
		}

		if message.FromBroadcast() && element.IsSetBroadcasts() && !element.Broadcasts.Boolean() {
			sublogentry.Debugf("ignoring broadcast message: %s", message.Id)
			continue
		}

		if message.Type == whatsapp.CallMessageType && element.IsSetCalls() && !element.Calls.Boolean() {
			sublogentry.Debugf("ignoring call message: %s", message.Id)
			continue
		}

		if !message.FromInternal || (element.ForwardInternal && (len(element.TrackId) == 0 || element.TrackId != message.TrackId)) {
			elerr := element.Post(wid, message)
			if elerr != nil {
				sublogentry.Errorf("error on post webhook: %s", elerr.Error())
			}
		}
	}

	return
}

// region FIND|SEARCH WHATSAPP SERVER
var ErrServerNotFound error = errors.New("the requested whatsapp server was not found")

func GetServerFromID(source string) (server *QpWhatsappServer, err error) {
	server, ok := WhatsappService.Servers[source]
	if !ok {
		err = ErrServerNotFound
		return
	}
	return
}

func GetServerFromBot(source QPBot) (server *QpWhatsappServer, err error) {
	return GetServerFromID(source.Wid)
}

func GetServerFromToken(token string) (server *QpWhatsappServer, err error) {
	for _, item := range WhatsappService.Servers {
		if item != nil && strings.EqualFold(item.Token, token) {
			server = item
			break
		}
	}

	if server == nil {
		err = ErrServerNotFound
	}

	return
}

func GetServersForUserID(user string) (servers map[string]*QpWhatsappServer) {
	return WhatsappService.GetServersForUser(user)
}

func GetServersForUser(user *QpUser) (servers map[string]*QpWhatsappServer) {
	return GetServersForUserID(user.Username)
}

//endregion

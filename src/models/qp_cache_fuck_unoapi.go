package models

import (
	"context"
	"reflect"
	"strings"

	whatsapp "github.com/nocodeleaks/quepasa/whatsapp"
	log "github.com/sirupsen/logrus"
	"go.mau.fi/whatsmeow/proto/waE2E"
)

func ValidateItemBecauseUNOAPIConflict(item QpCacheItem, from string, previous any) bool {
	// debugging messages in cache
	if strings.HasPrefix(from, "message") {

		prevItem := previous.(QpCacheItem)

		logentry := log.New().WithContext(context.Background())
		logentry = logentry.WithField(LogFields.MessageId, item.Key)
		logentry = logentry.WithField("from", from)
		logentry.Level = log.DebugLevel

		logentry.Debug("updating cache item ...")
		logentry.Debugf("old type: %s, %v", reflect.TypeOf(prevItem.Value), prevItem.Value)
		logentry.Debugf("new type: %s, %v", reflect.TypeOf(item.Value), item.Value)
		logentry.Debugf("equals: %v, deep equals: %v", item.Value == prevItem.Value, reflect.DeepEqual(item.Value, prevItem.Value))

		var prevContent interface{}
		if prevWaMsg, ok := prevItem.Value.(*whatsapp.WhatsappMessage); ok {
			prevContent = prevWaMsg.Content

			if nee, ok := prevContent.(*waE2E.Message); ok {
				if neeETM, ok := prevContent.(*waE2E.ExtendedTextMessage); ok {
					prevContent = neeETM.Text
					logentry.Debugf("old content from .ExtendedTextMessage as string: %s", prevContent)
				} else {
					prevContent = nee.String()
					logentry.Debugf("old content from .Message as string: %s", prevContent)
				}
			}
		}

		var newContent interface{}
		if newWaMsg, ok := item.Value.(*whatsapp.WhatsappMessage); ok {
			newContent = newWaMsg.Content

			if nee, ok := newContent.(*waE2E.Message); ok {
				conversation := nee.GetConversation()
				if len(conversation) > 0 {
					newContent = conversation
					logentry.Debugf("new content from .Message.Conversation: %s", newContent)
				} else {
					newContent = nee.String()
					logentry.Debugf("new content as string: %s", newContent)
				}
			}
		}

		if prevContent != nil && newContent != nil {
			logentry.Infof("content equals: %v, content deep equals: %v", prevContent == newContent, reflect.DeepEqual(prevContent, newContent))

			// if equals, deny triggers
			return !reflect.DeepEqual(prevContent, newContent)
		}
	}

	return true
}

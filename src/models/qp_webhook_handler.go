package models

import (
	"context"
	"reflect"
	"strings"

	whatsapp "github.com/nocodeleaks/quepasa/whatsapp"
	log "github.com/sirupsen/logrus"
)

type QPWebhookHandler struct {
	server *QpWhatsappServer
}

func (source QPWebhookHandler) GetLogger() *log.Entry {
	if source.server != nil {
		return source.server.GetLogger()
	}

	return log.WithContext(context.Background())
}

func (source *QPWebhookHandler) HandleWebHook(payload *whatsapp.WhatsappMessage) {
	if !source.HasWebhook() {
		return
	}

	logger := source.GetLogger()
	if payload.Type == whatsapp.DiscardMessageType || payload.Type == whatsapp.UnknownMessageType {
		logger.Debugf("ignoring discard|unknown message type on webhook request: %v", reflect.TypeOf(&payload))
		return
	}

	if payload.Type == whatsapp.TextMessageType && len(strings.TrimSpace(payload.Text)) <= 0 {
		logger.Debugf("ignoring empty text message on webhook request: %s", payload.Id)
		return
	}

	err := PostToWebHookFromServer(source.server, payload)
	if err != nil {
		logger.Errorf("error on handle webhook distributions: %s", err.Error())
	}
}

func (source *QPWebhookHandler) HasWebhook() bool {
	if source.server != nil {
		return len(source.server.Webhooks) > 0
	}
	return false
}

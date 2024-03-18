package models

import (
	"context"
	"fmt"

	whatsapp "github.com/nocodeleaks/quepasa/whatsapp"
	log "github.com/sirupsen/logrus"
)

type QpWhatsappServerWebhook struct {
	*QpWebhook

	server *QpWhatsappServer
}

func (source *QpWhatsappServerWebhook) GetLogger() *log.Entry {
	if source != nil && source.server != nil {
		logentry := source.server.GetLogger()
		if source.QpWebhook != nil {
			return logentry.WithField("url", source.QpWebhook.Url)
		}
		return logentry
	}

	logger := log.New()
	return logger.WithContext(context.Background())
}

//#region IMPLEMENTING WHATSAPP OPTIONS INTERFACE

func (source *QpWhatsappServerWebhook) GetOptions() *whatsapp.WhatsappOptions {
	if source == nil {
		return nil
	}

	return &source.WhatsappOptions
}

//#endregion

func (source *QpWhatsappServerWebhook) Save() (err error) {

	if source == nil {
		err = fmt.Errorf("nil webhook source")
		return err
	}

	if source.server == nil {
		err = fmt.Errorf("nil server")
		return err
	}

	if source.QpWebhook == nil {
		err = fmt.Errorf("nil source webhook")
		return err
	}

	logentry := source.GetLogger()
	logentry.Debugf("saving webhook info: %+v", source)

	affected, err := source.server.WebhookAddOrUpdate(source.QpWebhook)
	if err == nil {
		logentry.Infof("saved webhook with %v affected rows", affected)
	}

	return err
}

func (source *QpWhatsappServerWebhook) ToggleForwardInternal() (handle bool, err error) {
	source.ForwardInternal = !source.ForwardInternal
	return source.ForwardInternal, source.Save()
}

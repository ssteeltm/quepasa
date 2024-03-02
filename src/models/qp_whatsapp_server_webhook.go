package models

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
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

func (source *QpWhatsappServerWebhook) ToggleBroadcasts() (handle *bool, err error) {
	if source.Broadcasts == nil {
		source.Broadcasts = proto.Bool(true)
	} else if *source.Broadcasts {
		source.Broadcasts = proto.Bool(false)
	} else {
		source.Broadcasts = nil
	}
	return source.Broadcasts, source.Save()
}

func (source *QpWhatsappServerWebhook) ToggleReadReceipts() (handle *bool, err error) {
	if source.ReadReceipts == nil {
		source.ReadReceipts = proto.Bool(true)
	} else if *source.ReadReceipts {
		source.ReadReceipts = proto.Bool(false)
	} else {
		source.ReadReceipts = nil
	}
	return source.ReadReceipts, source.Save()
}

func (source *QpWhatsappServerWebhook) ToggleGroups() (handle *bool, err error) {
	if source.Groups == nil {
		source.Groups = proto.Bool(true)
	} else if *source.Groups {
		source.Groups = proto.Bool(false)
	} else {
		source.Groups = nil
	}
	return source.Groups, source.Save()
}

package models

import (
	"context"

	"github.com/google/uuid"
	library "github.com/nocodeleaks/quepasa/library"
	whatsapp "github.com/nocodeleaks/quepasa/whatsapp"
	log "github.com/sirupsen/logrus"
)

type QpWhatsappPairing struct {
	// Public token
	Token string `db:"token" json:"token" validate:"max=100"`

	// Whatsapp session id
	Wid string `db:"wid" json:"wid" validate:"max=255"`

	User *QpUser `json:"user,omitempty"`

	conn whatsapp.IWhatsappConnection `json:"-"`
}

func (source *QpWhatsappPairing) GetLogger() *log.Entry {
	if source.conn != nil && !source.conn.IsInterfaceNil() {
		return source.conn.GetLogger()
	}

	logger := log.WithContext(context.Background())

	if len(source.Token) > 0 {
		logger = logger.WithField("token", source.Token)
	}

	if len(source.Wid) > 0 {
		logger = logger.WithField("wid", source.Wid)
	}

	return logger
}

func (source *QpWhatsappPairing) OnPaired(wid string) {
	source.Wid = wid

	// updating token if from user
	if source.User != nil {
		source.Token = source.GetUserToken()
	}

	logger := source.GetLogger()
	if source.conn != nil {
		options := source.conn.GetOptions()
		if options != nil {
			options.EnableAutoReconnect = true
			options.Wid = source.Wid
			options.Logger = logger
		}
	}

	logger.Info("paired whatsapp section")
	server, err := WhatsappService.AppendPaired(source)
	if err != nil {
		logger.Errorf("paired error: %s", err.Error())
		return
	}

	go server.EnsureReady()
}

func (source *QpWhatsappPairing) GetConnection() (whatsapp.IWhatsappConnection, error) {
	if source.conn == nil {
		conn, err := NewEmptyConnection(source.OnPaired)
		if err != nil {
			return nil, err
		}
		source.conn = conn
	}

	return source.conn, nil
}

func (source *QpWhatsappPairing) GetUserToken() string {
	phone := library.GetPhoneByWId(source.Wid)

	logger := source.GetLogger()
	logger.Infof("wid to phone: %s", phone)

	servers := WhatsappService.GetServersForUser(source.User.Username)
	for _, item := range servers {
		if item.GetNumber() == phone {
			return item.Token
		}
	}

	return uuid.New().String()
}

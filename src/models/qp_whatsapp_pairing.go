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

	// if token was not setted
	// remember that user may want a different section for the same whatsapp
	if len(source.Token) == 0 {

		// updating token if from user
		if source.User != nil {
			source.Token = source.GetUserToken()
		}
	}

	if source.conn != nil {
		source.conn.SetReconnect(true)
	}

	logentry := source.GetLogger()
	logentry.Info("paired whatsapp section")
	server, err := WhatsappService.AppendPaired(source)
	if err != nil {
		logentry.Errorf("paired error: %s", err.Error())
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

// gets an existent token for same phone number and user, or create a new token
func (source *QpWhatsappPairing) GetUserToken() string {
	phone := library.GetPhoneByWId(source.Wid)

	logentry := source.GetLogger()
	logentry.Infof("wid to phone: %s", phone)

	servers := WhatsappService.GetServersForUser(source.User.Username)
	for _, item := range servers {
		if item.GetNumber() == phone {
			return item.Token
		}
	}

	return uuid.New().String()
}

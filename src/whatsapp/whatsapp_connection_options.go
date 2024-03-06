package whatsapp

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type WhatsappConnectionOptions struct {

	// whatsapp connection id and session
	Wid string `json:"-"`

	// should emit read receipts
	ReadReceipts *bool `json:"readreceipts,omitempty"`

	// should auto reject calls
	RejectCalls *bool `json:"rejectcalls,omitempty"`

	// should auto reconnect, false for qrcode scanner
	EnableAutoReconnect bool `json:"enableautoreconnect"`

	// log entry
	Logger *log.Entry `json:"-"`
}

// get default log entry, never nil
func (source WhatsappConnectionOptions) GetLogger() *log.Entry {
	if source.Logger == nil {
		logger := log.StandardLogger()
		logger.SetLevel(log.ErrorLevel)

		serverLogEntry := logger.WithContext(context.Background())

		if len(source.Wid) > 0 {
			serverLogEntry = serverLogEntry.WithField("wid", source.Wid)
		}

		source.Logger = serverLogEntry
	}

	return source.Logger
}

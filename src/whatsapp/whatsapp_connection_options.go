package whatsapp

import (
	"context"
	"reflect"

	"github.com/nocodeleaks/quepasa/library"
	log "github.com/sirupsen/logrus"
)

// Used only as parameters for start a new connection, wont propagate
type WhatsappConnectionOptions struct {
	library.LogStruct // logging
	*WhatsappOptions

	Wid       string
	Reconnect bool
}

func (source *WhatsappConnectionOptions) GetWid() string {
	return source.Wid
}

// should auto reconnect, false for qrcode scanner
func (source *WhatsappConnectionOptions) SetReconnect(value bool) {
	source.Reconnect = value
}

func (source *WhatsappConnectionOptions) GetReconnect() bool {
	return source.Reconnect
}

// get default log entry, never nil
func (source *WhatsappConnectionOptions) GetLogger() *log.Entry {
	if source != nil && source.LogEntry != nil {
		return source.LogEntry
	}

	logentry := log.WithContext(context.Background())
	logentry.Level = log.ErrorLevel
	logentry.Infof("generating new log entry for %s, with level: %s", reflect.TypeOf(source), logentry.Level)

	if source != nil {
		logentry = logentry.WithField(library.LogFields.WId, source.Wid)
		source.LogEntry = logentry
	}

	return logentry
}

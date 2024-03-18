package whatsapp

import log "github.com/sirupsen/logrus"

// Used only as parameters for start a new connection, wont propagate
type WhatsappConnectionOptions struct {
	*WhatsappOptions

	Wid       string
	Reconnect bool

	LogEntry *log.Entry
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

func (source *WhatsappConnectionOptions) GetLogger() *log.Entry {
	return source.LogEntry
}

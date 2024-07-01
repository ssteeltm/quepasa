package whatsmeow

import (
	"github.com/sirupsen/logrus"
	types "go.mau.fi/whatsmeow/types"
)

const WhatsmeowLogLevel = logrus.WarnLevel // default log level for whatsmeow
const WhatsmeowClientLogLevel = "INFO"     // default log level for whatsmeow client
const WhatsmeowDBLogLevel = "WARN"         // default log level for whatsmeow database

const WhatsmeowPresence = types.PresenceUnavailable

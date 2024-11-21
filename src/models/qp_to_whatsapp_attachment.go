package models

import whatsapp "github.com/nocodeleaks/quepasa/whatsapp"

type QpToWhatsappAttachment struct {
	Attach *whatsapp.WhatsappAttachment
	Debug  []string `json:"debug,omitempty"`
}

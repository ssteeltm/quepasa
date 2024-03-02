package models

import (
	whatsapp "github.com/nocodeleaks/quepasa/whatsapp"
	whatsmeow "github.com/nocodeleaks/quepasa/whatsmeow"
)

type QPFormAccountData struct {
	PageTitle    string
	ErrorMessage string
	Version      string
	Servers      map[string]*QpWhatsappServer
	User         QpUser
	Options      whatsapp.WhatsappOptions   `json:"options,omitempty"`
	WMOptions    whatsmeow.WhatsmeowOptions `json:"wmoptions,omitempty"`
}

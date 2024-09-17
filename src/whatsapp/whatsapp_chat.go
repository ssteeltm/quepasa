package whatsapp

import "strings"

type WhatsappChat struct {
	Id    string `json:"id"`
	Title string `json:"title,omitempty"`
}

var WASYSTEMCHAT = WhatsappChat{Id: "system", Title: "Internal System Message"}

func (source *WhatsappChat) FormatContact() {
	// removing session id
	if strings.Contains(source.Id, ":") {
		prefix := strings.Split(source.Id, ":")[0]
		suffix := strings.Split(source.Id, "@")[1]
		source.Id = prefix + "@" + suffix
	}
}

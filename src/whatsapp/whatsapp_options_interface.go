package whatsapp

type IWhatsappOptions interface {
	GetOptions() *WhatsappOptions
	Save() error
}

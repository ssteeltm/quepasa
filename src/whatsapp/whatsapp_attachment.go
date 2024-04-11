package whatsapp

import "strings"

type WhatsappAttachment struct {
	content *[]byte `json:"-"`

	// means that it can be downloaded from whatsapp servers
	CanDownload bool `json:"-"`

	Mimetype string `json:"mime"`

	// important to navigate throw content, declared file length
	FileLength uint64 `json:"filelength"`

	// document
	FileName string `json:"filename,omitempty"`

	// video | image | location (base64 image)
	JpegThumbnail string `json:"thumbnail,omitempty"`

	// audio
	Seconds uint32 `json:"seconds,omitempty"`

	// location msgs
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Sequence  int64   `json:"sequence,omitempty"` // live location

	// Public access url helper content
	Url string `json:"url,omitempty"`
}

func (source *WhatsappAttachment) GetContent() *[]byte {
	return source.content
}

func (source *WhatsappAttachment) SetContent(content *[]byte) {
	source.content = content
}

func (source *WhatsappAttachment) HasContent() bool {
	return nil != source.content
}

// used at receive.tmpl view
func (source *WhatsappAttachment) IsValidSize() bool {
	if source.FileLength > 500 {
		return true
	}

	if source.content != nil {
		length := len(*source.content)
		if length > 500 {
			return true
		}

		// there are many simple vcards with low bytes
		if strings.Contains(source.Mimetype, "vcard") && length > 70 {
			return true
		}
	}
	return false
}

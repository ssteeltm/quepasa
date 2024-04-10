package whatsapp

type WhatsappAttachment struct {
	content *[]byte `json:"-"`

	Mimetype string `json:"mime"`

	// important to navigate throw content
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

	if nil != source.content {
		if len(*source.content) > 500 {
			return true
		}
	}
	return false
}

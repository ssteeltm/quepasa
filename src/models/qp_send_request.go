package models

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	library "github.com/nocodeleaks/quepasa/library"
	whatsapp "github.com/nocodeleaks/quepasa/whatsapp"
	log "github.com/sirupsen/logrus"
)

type QpSendRequest struct {
	// (Optional) Used if passed
	Id string `json:"id,omitempty"`

	// Recipient of this message
	ChatId string `json:"chatId"`

	// (Optional) TrackId - less priority (urlparam -> query -> header -> body)
	TrackId string `json:"trackId,omitempty"`

	Text string `json:"text,omitempty"`

	// Msg in reply of another ? Message ID
	InReply string `json:"inreply,omitempty"`

	// (Optional) Sugested filename on user download
	FileName string `json:"fileName,omitempty"`

	// (Optional) important to navigate throw content
	FileLength uint64 `json:"filelength,omitempty"`

	// (Optional) mime type to facilitate correct delivery
	Mimetype string `json:"mime,omitempty"`

	Content []byte
}

func (source *QpSendRequest) EnsureChatId(r *http.Request) (err error) {
	if len(source.ChatId) == 0 {
		source.ChatId = GetChatId(r)
	}

	if len(source.ChatId) == 0 {
		err = fmt.Errorf("chat id missing")
	}
	return
}

func (source *QpSendRequest) EnsureValidChatId(r *http.Request) (err error) {
	err = source.EnsureChatId(r)
	if err != nil {
		return
	}

	chatid, err := whatsapp.FormatEndpoint(source.ChatId)
	if err != nil {
		return
	}

	source.ChatId = chatid
	return
}

func (source *QpSendRequest) ToWhatsappMessage() (msg *whatsapp.WhatsappMessage, err error) {
	chatId, err := whatsapp.FormatEndpoint(source.ChatId)
	if err != nil {
		return
	}

	chat := whatsapp.WhatsappChat{Id: chatId}
	msg = &whatsapp.WhatsappMessage{
		Id:           source.Id,
		TrackId:      source.TrackId,
		InReply:      source.InReply,
		Text:         source.Text,
		Chat:         chat,
		FromMe:       true,
		FromInternal: true,
	}

	// setting default type
	if len(msg.Text) > 0 {
		msg.Type = whatsapp.TextMessageType
	}

	return
}

func (source *QpSendRequest) ToWhatsappAttachment() (attach *whatsapp.WhatsappAttachment, err error) {
	contentMime := library.GetMimeTypeFromContent(source.Content, source.FileName)

	var requestExtension string
	if len(source.FileName) > 0 {
		requestExtension = filepath.Ext(source.FileName)
	} else if len(source.Mimetype) > 0 {
		requestExtension, _ = library.TryGetExtensionFromMimeType(source.Mimetype)
	}

	contentExtension, _ := library.TryGetExtensionFromMimeType(contentMime)

	attach = &whatsapp.WhatsappAttachment{
		Mimetype:   source.Mimetype,
		FileLength: source.FileLength,
		FileName:   source.FileName,
	}

	if len(attach.Mimetype) == 0 {
		attach.Mimetype = contentMime
	}

	// validating mime information
	if requestExtension != contentExtension {
		// invalid attachment
		log.Warnf("invalid mime for attachment, request extension: %s != content extension: %s :: content mime: %s, revalidating for security", requestExtension, contentExtension, contentMime)
		attach.Mimetype = contentMime
		attach.FileName = "invalid-" + library.GenerateFileNameFromMimeType(contentMime)
	}

	// validating content length
	contentLength := uint64(len(source.Content))
	if source.FileLength > 0 && source.FileLength != contentLength {
		log.Warnf("invalid attachment length, request length: %v != content length: %v, revalidating for security", source.FileLength, contentLength)
		attach.FileLength = uint64(len(source.Content))
	}

	// adjusting codec for ptt audio messages
	// inserting a trick for change from wave to ogg ... insecure
	convertFromWav := ENV.ShouldConvertWaveToOgg() && strings.Contains(attach.Mimetype, "wav")
	if (strings.Contains(attach.Mimetype, "ogg") || convertFromWav) && !strings.Contains(attach.Mimetype, "opus") {
		attach.Mimetype = "audio/ogg; codecs=opus"
	}

	log.Debugf("detected mime type: %s, filename: %s", attach.Mimetype, attach.FileName)

	// Defining a filename if not found before
	if len(attach.FileName) == 0 {
		attach.FileName = library.GenerateFileNameFromMimeType(attach.Mimetype)
	}

	attach.SetContent(&source.Content)
	return
}

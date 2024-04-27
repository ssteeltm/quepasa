package models

import (
	"context"
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

// get default log entry, never nil
func (source *QpSendRequest) GetLogger() *log.Entry {
	logentry := log.WithContext(context.Background())
	if len(source.ChatId) > 0 {
		logentry.WithField("chatid", source.ChatId)
	}

	return logentry
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
	contentLength := len(source.Content)
	if contentLength == 0 {
		return
	}

	logentry := source.GetLogger()

	attach = &whatsapp.WhatsappAttachment{
		CanDownload: false,
		Mimetype:    source.Mimetype,
		FileLength:  source.FileLength,
		FileName:    source.FileName,
	}

	// validating content length
	uIntContentLength := uint64(contentLength)
	if attach.FileLength != uIntContentLength {
		logentry.Warnf("invalid attachment length, request length: %v != content length: %v, revalidating for security", attach.FileLength, contentLength)
		attach.FileLength = uIntContentLength
	}

	// end source use and set content
	attach.SetContent(&source.Content)

	SecureAndCustomizeAttach(attach, logentry)
	return
}

func SecureAndCustomizeAttach(attach *whatsapp.WhatsappAttachment, logentry *log.Entry) {
	if attach == nil {
		return
	}

	var contentMime string
	content := attach.GetContent()
	if content != nil {
		contentMime = library.GetMimeTypeFromContent(*content)
		logentry.Debugf("send request, detected mime type from content: %s", contentMime)
	}

	var requestExtension string
	if len(attach.FileName) > 0 {
		requestExtension = filepath.Ext(attach.FileName)
		logentry.Debugf("send request, detected extension from filename: %s", requestExtension)

	} else if len(attach.Mimetype) > 0 {
		requestExtension, _ = library.TryGetExtensionFromMimeType(attach.Mimetype)
		logentry.Debugf("send request, detected extension from mime type: %s", requestExtension)
	}

	if len(contentMime) > 0 {

		if len(attach.Mimetype) == 0 {
			attach.Mimetype = contentMime
			logentry.Debugf("send request, updating empty mime type from content: %s", contentMime)
		}

		contentExtension, success := library.TryGetExtensionFromMimeType(contentMime)
		if success {
			logentry.Debugf("send request, content extension: %s", contentExtension)

			// validating mime information
			if requestExtension != contentExtension {
				// invalid attachment
				logentry.Warnf("send request, invalid mime for attachment, request extension: %s != content extension: %s :: content mime: %s, revalidating for security", requestExtension, contentExtension, contentMime)
				attach.Mimetype = contentMime
				attach.FileName = whatsapp.InvalidFilePrefix + library.GenerateFileNameFromMimeType(contentMime)
			}
		}
	}

	// set compatible audios to be sent as ptt
	ForceCompatiblePTT := ENV.UseCompatibleMIMEsAsAudio()
	if ForceCompatiblePTT && !attach.IsValidAudio() && IsCompatibleWithPTT(attach.Mimetype) {
		logentry.Infof("send request, setting that it should be sent as ptt, regards its incompatible mime type: %s", attach.Mimetype)
		attach.SetPTTCompatible(true)
	}

	// Defining a filename if not found before
	if len(attach.FileName) == 0 {
		attach.FileName = library.GenerateFileNameFromMimeType(attach.Mimetype)
		logentry.Debugf("send request, empty file name, generating a new one based on mime type: %s, file name: %s", attach.Mimetype, attach.FileName)
	}

	logentry.Debugf("send request, resolved mime type: %s, filename: %s", attach.Mimetype, attach.FileName)
}

func IsCompatibleWithPTT(mime string) bool {
	// switch for basic mime type, ignoring suffix
	mimeOnly := strings.Split(mime, ";")[0]

	for _, item := range whatsapp.WhatsappMIMEAudioPTTCompatible {
		if item == mimeOnly {
			return true
		}
	}

	return false
}

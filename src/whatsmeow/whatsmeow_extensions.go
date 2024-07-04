package whatsmeow

import (
	"encoding/base64"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	whatsapp "github.com/nocodeleaks/quepasa/whatsapp"
	whatsmeow "go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	types "go.mau.fi/whatsmeow/types"
)

func GetMediaTypeFromAttachment(source *whatsapp.WhatsappAttachment) whatsmeow.MediaType {
	msgType := whatsapp.GetMessageType(source)
	return GetMediaTypeFromWAMsgType(msgType)
}

// Traz o MediaType para download do whatsapp
func GetMediaTypeFromWAMsgType(msgType whatsapp.WhatsappMessageType) whatsmeow.MediaType {
	switch msgType {
	case whatsapp.ImageMessageType:
		return whatsmeow.MediaImage
	case whatsapp.AudioMessageType:
		return whatsmeow.MediaAudio
	case whatsapp.VideoMessageType:
		return whatsmeow.MediaVideo
	default:
		return whatsmeow.MediaDocument
	}
}

func ToWhatsmeowMessage(source whatsapp.IWhatsappMessage) (msg *waE2E.Message, err error) {
	messageText := source.GetText()

	if !source.HasAttachment() {
		internal := &waE2E.ExtendedTextMessage{Text: &messageText}
		msg = &waE2E.Message{ExtendedTextMessage: internal}
	}

	return
}

func NewWhatsmeowMessageAttachment(response whatsmeow.UploadResponse, waMsg whatsapp.WhatsappMessage, media whatsmeow.MediaType) (msg *waE2E.Message) {
	attach := waMsg.Attachment

	var seconds *uint32
	if attach.Seconds > 0 {
		seconds = proto.Uint32(attach.Seconds)
	}

	var mimetype *string
	if len(attach.Mimetype) > 0 {
		mimetype = proto.String(attach.Mimetype)
	}

	switch media {
	case whatsmeow.MediaImage:
		internal := &waE2E.ImageMessage{
			URL:           proto.String(response.URL),
			DirectPath:    proto.String(response.DirectPath),
			MediaKey:      response.MediaKey,
			FileEncSHA256: response.FileEncSHA256,
			FileSHA256:    response.FileSHA256,
			FileLength:    proto.Uint64(response.FileLength),

			Mimetype: mimetype,
			Caption:  proto.String(waMsg.Text),
		}
		msg = &waE2E.Message{ImageMessage: internal}
		return
	case whatsmeow.MediaAudio:

		var ptt *bool
		if attach.IsValidPTT() {
			ptt = proto.Bool(true)
		} else if attach.IsPTTCompatible() { // trick to send audio as ptt, "technical resource"
			ptt = proto.Bool(true)
			mimetype = proto.String(whatsapp.WhatsappPTTMime)
		}

		internal := &waE2E.AudioMessage{
			URL:           proto.String(response.URL),
			DirectPath:    proto.String(response.DirectPath),
			MediaKey:      response.MediaKey,
			FileEncSHA256: response.FileEncSHA256,
			FileSHA256:    response.FileSHA256,
			FileLength:    proto.Uint64(response.FileLength),
			Seconds:       seconds,
			Mimetype:      mimetype,
			PTT:           ptt,
		}
		msg = &waE2E.Message{AudioMessage: internal}
		return
	case whatsmeow.MediaVideo:
		internal := &waE2E.VideoMessage{
			URL:           proto.String(response.URL),
			DirectPath:    proto.String(response.DirectPath),
			MediaKey:      response.MediaKey,
			FileEncSHA256: response.FileEncSHA256,
			FileSHA256:    response.FileSHA256,
			FileLength:    proto.Uint64(response.FileLength),
			Seconds:       seconds,
			Mimetype:      mimetype,
			Caption:       proto.String(waMsg.Text),
		}
		msg = &waE2E.Message{VideoMessage: internal}
		return
	default:
		internal := &waE2E.DocumentMessage{
			URL:           proto.String(response.URL),
			DirectPath:    proto.String(response.DirectPath),
			MediaKey:      response.MediaKey,
			FileEncSHA256: response.FileEncSHA256,
			FileSHA256:    response.FileSHA256,
			FileLength:    proto.Uint64(response.FileLength),

			Mimetype: mimetype,
			FileName: proto.String(attach.FileName),
			Caption:  proto.String(waMsg.Text),
		}
		msg = &waE2E.Message{DocumentMessage: internal}
		return
	}
}

func GetStringFromBytes(bytes []byte) string {
	if len(bytes) > 0 {
		return base64.StdEncoding.EncodeToString(bytes)
	}
	return ""
}

// returns a valid chat title from local memory store
func GetChatTitle(client *whatsmeow.Client, jid types.JID) string {
	if jid.Server == "g.us" {
		gInfo, _ := client.GetGroupInfo(jid)
		if gInfo != nil {
			return gInfo.Name
		}
	} else {
		cInfo, _ := client.Store.Contacts.GetContact(jid)
		if cInfo.Found {
			if len(cInfo.BusinessName) > 0 {
				return cInfo.BusinessName
			} else if len(cInfo.FullName) > 0 {
				return cInfo.FullName
			} else {
				return cInfo.PushName
			}
		}
	}

	return ""
}

var BUTTONSMSGREGEX regexp.Regexp = *regexp.MustCompile(`(?i)(?P<content>.*)\s?[\$#]buttons:\[(?P<buttons>.*)\]\s?(?P<footer>.*)`)
var BUTTONSREGEXCONTENTINDEX int = BUTTONSMSGREGEX.SubexpIndex("content")
var BUTTONSREGEXFOOTERINDEX int = BUTTONSMSGREGEX.SubexpIndex("footer")
var BUTTONSREGEXBUTTONSINDEX int = BUTTONSMSGREGEX.SubexpIndex("buttons")

var RegexButton regexp.Regexp = *regexp.MustCompile(`\((?P<value>.*)\)(?P<display>.*)`)
var RegexButtonValue int = RegexButton.SubexpIndex("value")
var RegexButtonDisplay int = RegexButton.SubexpIndex("display")

func GenerateButtonsMessage(messageText string) *waE2E.ButtonsMessage {
	var contentText *string
	var footerText *string
	var buttons []*waE2E.ButtonsMessage_Button

	matches := BUTTONSMSGREGEX.FindStringSubmatch(messageText)
	contentMatched := matches[BUTTONSREGEXCONTENTINDEX]
	if len(contentMatched) > 0 {
		contentText = &contentMatched
	}

	footerMatched := matches[BUTTONSREGEXFOOTERINDEX]
	if len(footerMatched) > 0 {
		footerText = &footerMatched
	}

	buttonsText := matches[BUTTONSREGEXBUTTONSINDEX]
	buttonsSplited := strings.Split(buttonsText, ",")
	for _, s := range buttonsSplited {
		normalized := strings.TrimSpace(s)

		buttonText := &waE2E.ButtonsMessage_Button_ButtonText{}
		buttonText.DisplayText = &normalized
		buttonId := buttonText.DisplayText

		matchesButton := RegexButton.FindStringSubmatch(normalized)
		if len(matchesButton) > 0 {
			buttonValueMatched := matchesButton[RegexButtonValue]
			if len(buttonValueMatched) > 0 {
				buttonId = &buttonValueMatched
			}

			buttonDisplayMatched := matchesButton[RegexButtonDisplay]
			if len(buttonDisplayMatched) > 0 {
				buttonText.DisplayText = &buttonDisplayMatched
			}
		}

		buttonType := waE2E.ButtonsMessage_Button_RESPONSE
		buttons = append(buttons, &waE2E.ButtonsMessage_Button{ButtonText: buttonText, ButtonID: buttonId, Type: &buttonType})
	}

	headerType := waE2E.ButtonsMessage_EMPTY
	return &waE2E.ButtonsMessage{HeaderType: &headerType, ContentText: contentText, Buttons: buttons, FooterText: footerText}
}

func IsValidForButtons(text string) bool {
	lowerText := strings.ToLower(text)
	if strings.Contains(lowerText, "buttons:") {
		matches := BUTTONSMSGREGEX.FindStringSubmatch(text)
		if len(matches) > 0 {
			if len(strings.TrimSpace(matches[0])) > 0 {
				return true
			}
		}
	}
	return false
}

/*
<summary>

	Send defined presence when connecting and when the pushname is changed.
	This makes sure that outgoing messages always have the right pushname.

<summary/>
*/
func SendPresence(client *whatsmeow.Client, presence types.Presence, from string, logentry *log.Entry) {
	if len(client.Store.PushName) > 0 {
		err := client.SendPresence(presence)
		if err != nil {
			logentry.Warnf("failed to send presence: '%s', error: %s, from: %s", presence, err.Error(), from)
		} else {
			logentry.Debugf("marked self as '%s', from: %s", presence, from)
		}
	}
}

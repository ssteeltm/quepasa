package whatsmeow

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	slug "github.com/gosimple/slug"
	whatsapp "github.com/nocodeleaks/quepasa/whatsapp"
	log "github.com/sirupsen/logrus"
	proto "go.mau.fi/whatsmeow/binary/proto"
)

func HandleKnowingMessages(handler *WhatsmeowHandlers, out *whatsapp.WhatsappMessage, in *proto.Message) {
	logentry := handler.GetLogger()
	logentry.Tracef("handling knowing message: %v", in)

	switch {
	case in.ImageMessage != nil:
		HandleImageMessage(logentry, out, in.ImageMessage)
	case in.StickerMessage != nil:
		HandleStickerMessage(logentry, out, in.StickerMessage)
	case in.DocumentMessage != nil:
		HandleDocumentMessage(logentry, out, in.DocumentMessage)
	case in.AudioMessage != nil:
		HandleAudioMessage(logentry, out, in.AudioMessage)
	case in.VideoMessage != nil:
		HandleVideoMessage(logentry, out, in.VideoMessage)
	case in.ExtendedTextMessage != nil:
		HandleExtendedTextMessage(logentry, out, in.ExtendedTextMessage)
	case in.ButtonsResponseMessage != nil:
		HandleButtonsResponseMessage(logentry, out, in.ButtonsResponseMessage)
	case in.LocationMessage != nil:
		HandleLocationMessage(logentry, out, in.LocationMessage)
	case in.LiveLocationMessage != nil:
		HandleLiveLocationMessage(logentry, out, in.LiveLocationMessage)
	case in.ContactMessage != nil:
		HandleContactMessage(logentry, out, in.ContactMessage)
	case in.ReactionMessage != nil:
		HandleReactionMessage(logentry, out, in.ReactionMessage)
	case in.EditedMessage != nil:
		HandleEditTextMessage(logentry, out, in.EditedMessage)
	case in.ProtocolMessage != nil:
		HandleProtocolMessage(logentry, out, in.ProtocolMessage)
	case in.SenderKeyDistributionMessage != nil:
		out.Type = whatsapp.DiscardMessageType
	case in.StickerSyncRmrMessage != nil:
		out.Type = whatsapp.DiscardMessageType
	case len(in.GetConversation()) > 0:
		HandleTextMessage(logentry, out, in)
	default:
		out.Type = whatsapp.UnknownMessageType
		logentry.Warnf("message not handled: %v", in)
	}
}

//#region HANDLING TEXT MESSAGES

func HandleTextMessage(log *log.Entry, out *whatsapp.WhatsappMessage, in *proto.Message) {
	log.Debug("Received a text message !")
	out.Type = whatsapp.TextMessageType
	out.Text = in.GetConversation()
}

func HandleEditTextMessage(log *log.Entry, out *whatsapp.WhatsappMessage, in *proto.FutureProofMessage) {
	// never throws , obs !!!!
	// it came as a single text msg
	log.Debug("Received a edited text message !")
	out.Type = whatsapp.TextMessageType
	out.Text = in.String()
}

func HandleProtocolMessage(logentry *log.Entry, out *whatsapp.WhatsappMessage, in *proto.ProtocolMessage) {
	logentry.Trace("Received a protocol message !")

	switch v := in.GetType(); {
	case v == proto.ProtocolMessage_MESSAGE_EDIT:
		out.Type = whatsapp.TextMessageType
		out.Id = in.Key.GetId()
		out.Text = in.EditedMessage.GetConversation()
		out.Edited = true
		return

	case v == proto.ProtocolMessage_REVOKE:
		out.Id = in.Key.GetId()
		out.Type = whatsapp.RevokeMessageType
		return

	default:
		out.Type = whatsapp.UnknownMessageType
		b, err := json.Marshal(in)
		if err != nil {
			logentry.Error(err)
			return
		}

		out.Text = "ProtocolMessage :: " + string(b)
		return
	}
}

// Msg em resposta a outra
func HandleExtendedTextMessage(log *log.Entry, out *whatsapp.WhatsappMessage, in *proto.ExtendedTextMessage) {
	log.Debug("Received a text|extended message !")
	out.Type = whatsapp.TextMessageType

	out.Text = in.GetText()

	info := in.ContextInfo
	if info != nil {
		out.ForwardingScore = info.GetForwardingScore()
		out.InReply = info.GetStanzaId()
	}
}

func HandleReactionMessage(log *log.Entry, out *whatsapp.WhatsappMessage, in *proto.ReactionMessage) {
	log.Debug("Received a Reaction message!")

	out.Type = whatsapp.TextMessageType
	out.Text = in.GetText()
	out.InReply = in.Key.GetId()
}

//#endregion

func HandleButtonsResponseMessage(log *log.Entry, out *whatsapp.WhatsappMessage, in *proto.ButtonsResponseMessage) {
	log.Debug("Received a buttons response message !")
	out.Type = whatsapp.TextMessageType

	/*
		b, err := json.Marshal(in)
		if err != nil {
			log.Error(err)
			return
		}
		log.Debug(string(b))
	*/

	out.Text = in.GetSelectedButtonId()

	info := in.ContextInfo
	if info != nil {
		out.ForwardingScore = info.GetForwardingScore()
		out.InReply = info.GetStanzaId()
	}
}

func HandleImageMessage(log *log.Entry, out *whatsapp.WhatsappMessage, in *proto.ImageMessage) {
	log.Debug("Received an image message !")
	out.Type = whatsapp.ImageMessageType

	// in case of caption passed
	out.Text = in.GetCaption()

	jpeg := GetStringFromBytes(in.JpegThumbnail)
	out.Attachment = &whatsapp.WhatsappAttachment{
		CanDownload:   true,
		Mimetype:      in.GetMimetype(),
		FileLength:    in.GetFileLength(),
		JpegThumbnail: jpeg,
	}

	info := in.ContextInfo
	if info != nil {
		out.ForwardingScore = info.GetForwardingScore()
		out.InReply = info.GetStanzaId()
	}
}

func HandleStickerMessage(log *log.Entry, out *whatsapp.WhatsappMessage, in *proto.StickerMessage) {
	log.Debug("Received a image|sticker message !")

	if in.GetIsAnimated() {
		out.Type = whatsapp.VideoMessageType
	} else {
		out.Type = whatsapp.ImageMessageType
	}

	jpeg := GetStringFromBytes(in.PngThumbnail)
	out.Attachment = &whatsapp.WhatsappAttachment{
		CanDownload: true,
		Mimetype:    in.GetMimetype(),
		FileLength:  in.GetFileLength(),

		JpegThumbnail: jpeg,
	}
}

func HandleVideoMessage(log *log.Entry, out *whatsapp.WhatsappMessage, in *proto.VideoMessage) {
	log.Debug("Received a video message !")
	out.Type = whatsapp.VideoMessageType

	// in case of caption passed
	out.Text = in.GetCaption()

	jpeg := base64.StdEncoding.EncodeToString(in.JpegThumbnail)
	out.Attachment = &whatsapp.WhatsappAttachment{
		CanDownload: true,
		Mimetype:    in.GetMimetype(),
		FileLength:  in.GetFileLength(),

		JpegThumbnail: jpeg,
	}

	info := in.ContextInfo
	if info != nil {
		out.ForwardingScore = info.GetForwardingScore()
		out.InReply = info.GetStanzaId()
	}
}

func HandleDocumentMessage(log *log.Entry, out *whatsapp.WhatsappMessage, in *proto.DocumentMessage) {
	log.Debug("Received a document message !")
	out.Type = whatsapp.DocumentMessageType

	// in case of caption passed
	out.Text = in.GetCaption()

	jpeg := base64.StdEncoding.EncodeToString(in.JpegThumbnail)
	out.Attachment = &whatsapp.WhatsappAttachment{
		CanDownload: true,
		Mimetype:    in.GetMimetype(),
		FileLength:  in.GetFileLength(),

		FileName:      in.GetFileName(),
		JpegThumbnail: jpeg,
	}

	info := in.ContextInfo
	if info != nil {
		out.ForwardingScore = info.GetForwardingScore()
		out.InReply = info.GetStanzaId()
	}
}

func HandleAudioMessage(log *log.Entry, out *whatsapp.WhatsappMessage, in *proto.AudioMessage) {
	log.Debug("Received an audio message !")
	out.Type = whatsapp.AudioMessageType

	out.Attachment = &whatsapp.WhatsappAttachment{
		CanDownload: true,
		Mimetype:    in.GetMimetype(),
		FileLength:  in.GetFileLength(),
		Seconds:     in.GetSeconds(),
	}

	info := in.ContextInfo
	if info != nil {
		out.ForwardingScore = info.GetForwardingScore()
		out.InReply = info.GetStanzaId()
	}
}

func HandleLocationMessage(log *log.Entry, out *whatsapp.WhatsappMessage, in *proto.LocationMessage) {
	log.Debug("Received a Location message !")
	out.Type = whatsapp.LocationMessageType

	// in a near future, create a environment variable for that
	defaultUrl := "https://www.google.com/maps?ll={lat},{lon}&q={lat}+{lon}"

	defaultUrl = strings.Replace(defaultUrl, "{lat}", fmt.Sprintf("%f", *in.DegreesLatitude), -1)
	defaultUrl = strings.Replace(defaultUrl, "{lon}", fmt.Sprintf("%f", *in.DegreesLongitude), -1)

	filename := fmt.Sprintf("%f_%f", in.GetDegreesLatitude(), in.GetDegreesLongitude())
	filename = fmt.Sprintf("%s.url", slug.Make(filename))

	content := []byte("[InternetShortcut]\nURL=" + defaultUrl)
	length := uint64(len(content))
	jpeg := base64.StdEncoding.EncodeToString(in.JpegThumbnail)

	out.Attachment = &whatsapp.WhatsappAttachment{
		CanDownload:   false,
		Mimetype:      "text/x-uri; location",
		Latitude:      in.GetDegreesLatitude(),
		Longitude:     in.GetDegreesLongitude(),
		JpegThumbnail: jpeg,
		Url:           defaultUrl,
		FileName:      filename,
		FileLength:    length,
	}

	out.Attachment.SetContent(&content)
}

func HandleLiveLocationMessage(log *log.Entry, out *whatsapp.WhatsappMessage, in *proto.LiveLocationMessage) {
	log.Debug("Received a Live Location message !")
	out.Type = whatsapp.LocationMessageType

	// in case of caption passed
	out.Text = in.GetCaption()

	// in a near future, create a environment variable for that
	defaultUrl := "https://www.google.com/maps?ll={lat},{lon}&q={lat}+{lon}"

	defaultUrl = strings.Replace(defaultUrl, "{lat}", fmt.Sprintf("%f", *in.DegreesLatitude), -1)
	defaultUrl = strings.Replace(defaultUrl, "{lon}", fmt.Sprintf("%f", *in.DegreesLongitude), -1)

	filename := out.Text
	if len(filename) == 0 {
		filename = fmt.Sprintf("%f_%f", *in.DegreesLatitude, *in.DegreesLongitude)
	}
	filename = fmt.Sprintf("%s.url", slug.Make(filename))

	content := []byte("[InternetShortcut]\nURL=" + defaultUrl)
	length := uint64(len(content))
	jpeg := base64.StdEncoding.EncodeToString(in.JpegThumbnail)

	out.Attachment = &whatsapp.WhatsappAttachment{
		CanDownload:   false,
		Mimetype:      "text/x-uri; live location",
		Latitude:      in.GetDegreesLatitude(),
		Longitude:     in.GetDegreesLongitude(),
		Sequence:      in.GetSequenceNumber(),
		JpegThumbnail: jpeg,
		Url:           defaultUrl,
		FileName:      filename,
		FileLength:    length,
	}

	out.Attachment.SetContent(&content)
}

func HandleContactMessage(log *log.Entry, out *whatsapp.WhatsappMessage, in *proto.ContactMessage) {
	log.Debug("Received a Contact message !")
	out.Type = whatsapp.ContactMessageType

	out.Text = in.GetDisplayName()
	filename := out.Text
	if len(filename) == 0 {
		filename = out.Id
	}
	filename = fmt.Sprintf("%s.vcf", slug.Make(filename))

	content := []byte(in.GetVcard())
	length := uint64(len(content))

	out.Attachment = &whatsapp.WhatsappAttachment{
		CanDownload: false,
		Mimetype:    "text/x-vcard",
		FileName:    filename,
		FileLength:  length,
	}

	out.Attachment.SetContent(&content)
}

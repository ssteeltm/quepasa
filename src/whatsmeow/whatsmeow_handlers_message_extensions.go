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
	logger := handler.GetLogger()

	logger.Tracef("handling knowing message: %v", in)
	if in.ImageMessage != nil {
		HandleImageMessage(logger, out, in.ImageMessage)
	} else if in.StickerMessage != nil {
		HandleStickerMessage(logger, out, in.StickerMessage)
	} else if in.DocumentMessage != nil {
		HandleDocumentMessage(logger, out, in.DocumentMessage)
	} else if in.AudioMessage != nil {
		HandleAudioMessage(logger, out, in.AudioMessage)
	} else if in.VideoMessage != nil {
		HandleVideoMessage(logger, out, in.VideoMessage)
	} else if in.ExtendedTextMessage != nil {
		HandleExtendedTextMessage(logger, out, in.ExtendedTextMessage)
	} else if in.ButtonsResponseMessage != nil {
		HandleButtonsResponseMessage(logger, out, in.ButtonsResponseMessage)
	} else if in.LocationMessage != nil {
		HandleLocationMessage(logger, out, in.LocationMessage)
	} else if in.LiveLocationMessage != nil {
		HandleLiveLocationMessage(logger, out, in.LiveLocationMessage)
	} else if in.ContactMessage != nil {
		HandleContactMessage(logger, out, in.ContactMessage)
	} else if in.ReactionMessage != nil {
		HandleReactionMessage(logger, out, in.ReactionMessage)
	} else if in.EditedMessage != nil {
		HandleEditTextMessage(logger, out, in.EditedMessage)
	} else if in.ProtocolMessage != nil {
		HandleProtocolMessage(logger, out, in.ProtocolMessage)
	} else if in.SenderKeyDistributionMessage != nil {
		out.Type = whatsapp.DiscardMessageType
		b, err := json.Marshal(in)
		if err != nil {
			log.Error(err)
			return
		}

		out.Text = string(b)
	} else if len(in.GetConversation()) > 0 {
		HandleTextMessage(logger, out, in)
	} else {
		logger.Warnf("message not handled: %v", in)
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

func HandleProtocolMessage(log *log.Entry, out *whatsapp.WhatsappMessage, in *proto.ProtocolMessage) {
	log.Trace("Received a protocol message !")

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
		out.Type = whatsapp.DiscardMessageType
		b, err := json.Marshal(in)
		if err != nil {
			log.Error(err)
			return
		}

		out.Text = string(b)
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
	out.Text = in.GetTitle()

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

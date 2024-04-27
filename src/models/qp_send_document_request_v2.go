package models

type QPSendDocumentRequestV2 struct {
	Recipient  string         `json:"recipient,omitempty"`
	Message    string         `json:"message,omitempty"`
	Attachment QPAttachmentV1 `json:"attachment,omitempty"`
}

func (source *QPSendDocumentRequestV2) ToQpSendRequest() *QpSendRequest {
	var request *QpSendRequest
	if len(source.Attachment.Base64) > 0 {
		RequestEncoded := &QpSendRequestEncoded{
			Content: source.Attachment.Base64,
		}

		RequestEncoded.ChatId = source.Recipient
		RequestEncoded.GenerateContent()
		request = &RequestEncoded.QpSendRequest
		request.Mimetype = source.Attachment.MIME
		request.FileName = source.Attachment.FileName

	} else if len(source.Attachment.Url) > 0 {
		RequestUrl := &QpSendRequestUrl{
			Url: source.Attachment.Url,
		}

		RequestUrl.ChatId = source.Recipient
		RequestUrl.GenerateContent()
		request = &RequestUrl.QpSendRequest
	}

	request.Text = source.Message
	return request
}

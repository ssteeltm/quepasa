package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	metrics "github.com/nocodeleaks/quepasa/metrics"
	models "github.com/nocodeleaks/quepasa/models"
	whatsapp "github.com/nocodeleaks/quepasa/whatsapp"
)

// -------------------------- PUBLIC METHODS
//region TYPES OF SENDING

// SendAPIHandler renders route "/v3/bot/{token}/send"
func SendAny(w http.ResponseWriter, r *http.Request) {
	response := &models.QpSendResponse{}

	server, err := GetServer(r)
	if err != nil {
		metrics.MessageSendErrors.Inc()

		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	// Declare a new request struct.
	request := &models.QpSendAnyRequest{}

	if r.ContentLength > 0 && r.Method == http.MethodPost {
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err = json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			jsonErr := fmt.Errorf("invalid json body: %s", err.Error())
			response.ParseError(jsonErr)
			RespondInterface(w, response)
			return
		}
	}

	// Getting ChatId parameter
	err = request.EnsureValidChatId(r)
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	if len(request.Url) == 0 && r.URL.Query().Has("url") {
		request.Url = r.URL.Query().Get("url")
	}

	if len(request.Url) > 0 {

		// download content to byte array
		err = request.GenerateUrlContent()
		if err != nil {
			metrics.MessageSendErrors.Inc()
			response.ParseError(err)
			RespondInterface(w, response)
			return
		}
	} else if len(request.Content) > 0 {

		// BASE64 content to byte array
		err = request.GenerateEmbedContent()
		if err != nil {
			metrics.MessageSendErrors.Inc()
			response.ParseError(err)
			RespondInterface(w, response)
			return
		}
	}

	SendRequest(w, r, &request.QpSendRequest, server)
}

//endregion

// -------------------------- INTERNAL METHODS

// Send a request already validated with chatid and server
func SendRequest(w http.ResponseWriter, r *http.Request, request *models.QpSendRequest, server *models.QpWhatsappServer) {
	response := &models.QpSendResponse{}
	var err error

	attach, err := request.ToWhatsappAttachment()
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	// if not set, try to recover text from url
	if len(request.Text) == 0 && r.URL.Query().Has("text") {
		request.Text = r.URL.Query().Get("text")
	}

	if attach == nil && len(request.Text) == 0 {
		metrics.MessageSendErrors.Inc()
		err = fmt.Errorf("text not found, do not send empty messages")
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	// getting trackid if not passed in request
	if len(request.TrackId) == 0 {
		request.TrackId = GetTrackId(r)
	}

	Send(server, response, request, w, attach)
}

// finally sends to the whatsapp server
func Send(server *models.QpWhatsappServer, response *models.QpSendResponse, request *models.QpSendRequest, w http.ResponseWriter, attach *whatsapp.WhatsappAttachment) {
	waMsg, err := request.ToWhatsappMessage()
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	logentry := server.GetLogger()

	if attach != nil {
		waMsg.Attachment = attach
		waMsg.Type = whatsapp.GetMessageType(attach)
		logentry.Debugf("send attachment of type: %v, mime: %s, length: %v, filename: %s", waMsg.Type, attach.Mimetype, attach.FileLength, attach.FileName)
	} else {
		// test for poll, already set from ToWhatsappMessage
		waMsg.Type = whatsapp.TextMessageType
	}

	if waMsg.Type == whatsapp.UnknownMessageType {
		// correct msg type for texts contents
		if len(waMsg.Text) > 0 {
			waMsg.Type = whatsapp.TextMessageType
		}

		logentry.Errorf("unknown message type: %v", waMsg)
		// *** implement an error here if not found any knowing type
	}

	// Checking for ready state
	status := server.GetStatus()
	if status != whatsapp.Ready {
		err = &ApiServerNotReadyException{Wid: server.GetWId(), Status: status}
		response.ParseError(err)
		RespondInterfaceCode(w, response, http.StatusServiceUnavailable)
		return
	}

	sendResponse, err := server.SendMessage(waMsg)
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	// success
	metrics.MessagesSent.Inc()

	result := &models.QpSendResponseMessage{}
	result.Wid = server.GetWId()
	result.Id = sendResponse.GetId()
	result.ChatId = waMsg.Chat.Id
	result.TrackId = waMsg.TrackId

	response.ParseSuccess(result)
	RespondInterface(w, response)
}

package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	metrics "github.com/nocodeleaks/quepasa/metrics"
	models "github.com/nocodeleaks/quepasa/models"
	whatsapp "github.com/nocodeleaks/quepasa/whatsapp"
)

// ReceiveAPIHandler renders route GET "/{version}/bot/{token}/receive"
func ReceiveAPIHandler(w http.ResponseWriter, r *http.Request) {
	response := &models.QpReceiveResponse{}

	server, err := GetServer(r)
	if err != nil {
		metrics.MessageReceiveErrors.Inc()
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	// Checking for ready state
	status := server.GetStatus()
	if status != whatsapp.Ready {
		metrics.MessageReceiveErrors.Inc()
		err = &ApiServerNotReadyException{Wid: server.GetWId(), Status: status}
		response.ParseError(err)
		RespondInterfaceCode(w, response, http.StatusServiceUnavailable)
		return
	}

	if server.Handler == nil {
		metrics.MessageReceiveErrors.Inc()
		err = fmt.Errorf("handlers not attached")
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	response.Total = uint64(server.Handler.GetTotal())

	timestamp, err := GetTimestamp(r)
	if err != nil {
		metrics.MessageReceiveErrors.Inc()
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	messages := GetOrderedMessages(server, timestamp)
	metrics.MessagesReceived.Add(float64(len(messages)))

	response.Server = server.QpServer
	response.Messages = messages

	if timestamp > 0 {
		searchTime := time.Unix(timestamp, 0)
		msg := fmt.Sprintf("getting with timestamp: %v => %s", timestamp, searchTime)
		response.ParseSuccess(msg)
	} else {
		response.ParseSuccess("getting without filter")
	}

	RespondSuccess(w, response)
}

// SendAPIHandler renders route "/v3/bot/{token}/sendtext"
func SendText(w http.ResponseWriter, r *http.Request) {
	response := &models.QpSendResponse{}

	server, err := GetServer(r)
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	// Declare a new request struct.
	request := &models.QpSendRequest{}

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err = json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	if len(request.Text) == 0 {
		metrics.MessageSendErrors.Inc()
		err = fmt.Errorf("text not found, do not send empty messages")
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	// Getting ChatId parameter
	err = request.EnsureValidChatId(r)
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	// override trackid if passed throw any other way
	trackid := GetTrackId(r)
	if len(trackid) > 0 {
		request.TrackId = trackid
	}

	Send(server, response, request, w, nil)
}

/*
<summary>

	Renders route POST "/{version}/bot/{token}/sendbinary/{chatid}/{filename}/{text}"

	Any of then, at this order of priority
	Path parameters: {chatid}
	Path parameters: {filename}
	Path parameters: {text} only images
	Url parameters: ?chatid={chatid}
	Url parameters: ?filename={filename}
	Url parameters: ?text={text} only images
	Header parameters: X-QUEPASA-CHATID = {chatid}
	Header parameters: X-QUEPASA-FILENAME = {filename}
	Header parameters: X-QUEPASA-TEXT = {text} only images

</summary>
*/
func SendDocumentFromBinary(w http.ResponseWriter, r *http.Request) {
	response := &models.QpSendResponse{}

	server, err := GetServer(r)
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	// Declare a new request struct.
	request := &models.QpSendRequest{}

	// Getting ChatId parameter
	err = request.EnsureValidChatId(r)
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	content, err := io.ReadAll(r.Body)
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(fmt.Errorf("attachment content missing or read error"))
		RespondInterface(w, response)
		return
	}

	request.Content = content
	request.Mimetype = r.Header.Get("Content-Type")

	InformedLength := r.Header.Get("Content-Length")
	if len(InformedLength) > 0 {
		length, err := strconv.ParseUint(InformedLength, 10, 64)
		if err == nil {
			request.FileLength = length
		}
	}

	// Getting FileName
	request.FileName = GetFileName(r)

	// Getting Text Label
	request.Text = GetTextParameter(r)

	SendRequest(w, r, request, server)
}

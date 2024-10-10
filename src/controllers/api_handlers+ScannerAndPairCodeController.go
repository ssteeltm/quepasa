package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	models "github.com/nocodeleaks/quepasa/models"
	log "github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
)

func ScannerController(w http.ResponseWriter, r *http.Request) {
	// setting default response type as json
	w.Header().Set("Content-Type", "application/json")

	response := &models.QpResponse{}

	token := GetToken(r)
	if len(token) == 0 {
		err := fmt.Errorf("token not found")
		RespondBadRequest(w, err)
		return
	}

	username, err := GetUsername(r)
	if err != nil {
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	HSDString := models.GetRequestParameter(r, "historysyncdays")
	historysyncdays, _ := strconv.ParseUint(HSDString, 10, 32)

	pairing := &models.QpWhatsappPairing{
		Token:           token,
		Username:        username,
		HistorySyncDays: uint32(historysyncdays),
	}

	con, err := pairing.GetConnection()
	if err != nil {
		err := fmt.Errorf("cant get connection: %s", err.Error())
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	log.Infof("requesting qrcode for token %s", token)
	result := con.GetWhatsAppQRCode()

	var png []byte
	png, err = qrcode.Encode(result, qrcode.Medium, 256)
	if err != nil {
		err := fmt.Errorf("cant get qrcode: %s", err.Error())
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=qrcode.png")
	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(png))
}

func PairCodeController(w http.ResponseWriter, r *http.Request) {
	// setting default response type as json
	w.Header().Set("Content-Type", "application/json")

	response := &models.QpResponse{}

	token := GetToken(r)
	if len(token) == 0 {
		err := fmt.Errorf("token not found")
		RespondBadRequest(w, err)
		return
	}

	username, err := GetUsername(r)
	if err != nil {
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	pairing := &models.QpWhatsappPairing{Token: token, Username: username}
	con, err := pairing.GetConnection()
	if err != nil {
		err := fmt.Errorf("can't get connection: %s", err.Error())
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	phone := models.GetRequestParameter(r, "phone")
	if len(phone) == 0 {
		err := errors.New("missing phone number")
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	code, err := con.PairPhone(phone)
	if err != nil {
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	response.Success = true
	response.Status = code
	RespondSuccess(w, response)
}

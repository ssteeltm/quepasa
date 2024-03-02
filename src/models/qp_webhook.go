package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	whatsapp "github.com/nocodeleaks/quepasa/whatsapp"
	log "github.com/sirupsen/logrus"
)

type QpWebhook struct {
	Url             string      `db:"url" json:"url,omitempty"`                         // destination
	ForwardInternal bool        `db:"forwardinternal" json:"forwardinternal,omitempty"` // forward internal msg from api
	TrackId         string      `db:"trackid" json:"trackid,omitempty"`                 // identifier of remote system to avoid loop
	ReadReceipts    *bool       `db:"readreceipts" json:"readreceipts,omitempty"`       // should emit read receipts
	Groups          *bool       `db:"groups" json:"groups,omitempty"`                   // should handle groups messages
	Broadcasts      *bool       `db:"broadcasts" json:"broadcasts,omitempty"`           // should handle broadcast messages
	Extra           interface{} `db:"extra" json:"extra,omitempty"`                     // extra info to append on payload
	Failure         *time.Time  `json:"failure,omitempty"`                              // first failure timestamp
	Success         *time.Time  `json:"success,omitempty"`                              // last success timestamp
	Timestamp       *time.Time  `db:"timestamp" json:"timestamp,omitempty"`
}

//#region VIEWS TRICKS

func (source QpWebhook) GetReadReceipts() bool {
	return *source.ReadReceipts
}

func (source QpWebhook) IsSetReadReceipts() bool {
	return source.ReadReceipts != nil
}

func (source QpWebhook) GetGroups() bool {
	return *source.Groups
}

func (source QpWebhook) IsSetGroups() bool {
	return source.Groups != nil
}

func (source QpWebhook) GetBroadcasts() bool {
	return *source.Broadcasts
}

func (source QpWebhook) IsSetBroadcasts() bool {
	return source.Broadcasts != nil
}

func (source QpWebhook) IsSetExtra() bool {
	return source.Extra != nil
}

//#endregion

var ErrInvalidResponse error = errors.New("the requested url do not return 200 status code")

func (source *QpWebhook) Post(wid string, message *whatsapp.WhatsappMessage) (err error) {
	log.Infof("dispatching webhook from: %s, id: %s, to: %s", wid, message.Id, source.Url)

	payload := &QpWebhookPayload{
		WhatsappMessage: message,
		Extra:           source.Extra,
	}

	payloadJson, err := json.Marshal(&payload)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", source.Url, bytes.NewBuffer(payloadJson))
	req.Header.Set("User-Agent", "Quepasa")
	req.Header.Set("X-QUEPASA-WID", wid)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	client.Timeout = time.Second * 10
	resp, err := client.Do(req)
	if err != nil {
		log.Warnf("(%s) error at post webhook: %s", wid, err.Error())
	}

	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			err = ErrInvalidResponse
		}
	}

	time := time.Now().UTC()
	if err != nil {
		if source.Failure == nil {
			source.Failure = &time
		}
	} else {
		source.Failure = nil
		source.Success = &time
	}

	return
}

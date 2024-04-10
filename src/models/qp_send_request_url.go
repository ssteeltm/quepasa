package models

import (
	"fmt"
	"io"
	"net/http"
	"path"
)

type QpSendRequestUrl struct {
	QpSendRequest
	Url string `json:"url"`
}

func (source *QpSendRequestUrl) GenerateContent() (err error) {
	resp, err := http.Get(source.Url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("error on generate content from QpSendRequestUrl, unexpected status code: %v", resp.StatusCode)
		return
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	source.QpSendRequest.Content = content

	// setting filename if empty
	if len(source.QpSendRequest.FileName) == 0 {
		source.QpSendRequest.FileName = path.Base(source.Url)
	}

	return
}

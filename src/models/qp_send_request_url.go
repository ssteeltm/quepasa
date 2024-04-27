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

		logentry := source.GetLogger()
		logentry.Error(err)
		return
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	source.QpSendRequest.Content = content

	if resp.ContentLength > -1 {
		source.QpSendRequest.FileLength = uint64(resp.ContentLength)
	}

	source.QpSendRequest.Mimetype = resp.Header.Get("content-type")

	// setting filename if empty
	if len(source.QpSendRequest.FileName) == 0 {
		source.QpSendRequest.FileName = path.Base(source.Url)
	}

	return
}

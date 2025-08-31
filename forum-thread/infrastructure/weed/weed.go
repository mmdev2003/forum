package weed

import (
	"bytes"
	"github.com/linxGnu/goseaweedfs"
	"io"
	"net/http"
	"time"
)

func New(
	filer string,
	masterUrl string,
) *Weed {
	filers := []string{filer}
	sw, err := goseaweedfs.NewSeaweed(
		masterUrl,
		filers,
		8096,
		&http.Client{Timeout: 5 * time.Minute},
	)
	if err != nil {
		panic("Failed to create seaweed: " + err.Error())
	}
	return &Weed{sw}

}

type Weed struct {
	operation *goseaweedfs.Seaweed
}

func (weed *Weed) Upload(
	file []byte,
	filename string,
	size int64,
	collection string,
) (string, error) {
	fileData, err := weed.operation.Upload(
		bytes.NewReader(file),
		filename,
		size,
		collection,
		"",
	)
	if err != nil {
		return "", err
	}

	return fileData.FileID, nil
}

func (weed *Weed) Download(
	FileID string,
) ([]byte, error) {
	var fileBuffer bytes.Buffer

	_, err := weed.operation.Download(
		FileID,
		nil,
		func(r io.Reader) error {
			_, err := io.Copy(&fileBuffer, r)
			return err
		},
	)
	file := fileBuffer.Bytes()
	if err != nil {
		return file, err
	}

	return file, nil
}

func (weed *Weed) Update(
	file []byte,
	fileID string,
	filename string,
	size int64,
	collection string,
) error {
	err := weed.operation.Replace(
		fileID,
		bytes.NewReader(file),
		filename,
		size,
		collection,
		"",
		true,
	)
	if err != nil {
		return err
	}

	return nil
}

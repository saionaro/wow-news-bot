package imageloader

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

func DownloadImage(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	var data bytes.Buffer
	_, err = io.Copy(&data, res.Body)
	if err != nil {
		return nil, err
	}
	return data.Bytes(), nil
}

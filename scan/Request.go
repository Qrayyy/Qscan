package scan

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	client = &http.Client{
		Timeout: time.Second * 30,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
)

func Get(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36")
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	newBody, err := io.ReadAll(resp.Body)
	return string(newBody), nil
}

func post(Url string, contentType string, data interface{}) (string, error) {
	var (
		payload io.Reader
		err     error
		req     *http.Request
	)

	switch contentType {
	case "application/json":
		var jsonData []byte
		jsonData, err = json.Marshal(data)
		if err != nil {
			return "", err
		}
		payload = bytes.NewBuffer(jsonData)
		req, err = http.NewRequest("POST", Url, payload)
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")
	case "application/x-www-form-urlencoded":
		values := url.Values{}
		formData, ok := data.(map[string]string)
		if !ok {
			return "", errors.New("invalid form data")
		}
		for key, val := range formData {
			values.Add(key, val)
		}
		payload = strings.NewReader(values.Encode())
		req, err = http.NewRequest("POST", Url, payload)
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	case "multipart/form-data":
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		formData, ok := data.(map[string]string)
		if !ok {
			return "", errors.New("invalid form data")
		}
		for key, val := range formData {
			partHeader := textproto.MIMEHeader{}
			partHeader.Set("Content-Disposition", `form-data; name="`+key+`"`)
			part, err := writer.CreatePart(partHeader)
			if err != nil {
				return "", err
			}
			_, err = part.Write([]byte(val))
			if err != nil {
				return "", err
			}
		}
		uploadData, ok := data.(map[string]*os.File)
		if !ok {
			return "", errors.New("invalid file data")
		}
		for key, val := range uploadData {
			part, err := writer.CreateFormFile(key, filepath.Base(val.Name()))
			if err != nil {
				return "", err
			}
			_, err = io.Copy(part, val)
			if err != nil {
				return "", err
			}
		}
		writer.Close()
		req, err = http.NewRequest("POST", Url, body)
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
	default:
		return "", errors.New("unsupported content type")
	}

	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	newBody, err := io.ReadAll(resp.Body)
	return string(newBody), nil
}

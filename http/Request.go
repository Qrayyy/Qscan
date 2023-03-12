package http

import (
	"bytes"
	"crypto/tls"
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

type Response struct {
	Status        string
	StatusCode    int
	Body          string
	Header        http.Header
	ContentLength int
	RequestUrl    string
	Location      string
}

var (
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{
		Transport: tr,
		Timeout:   time.Second * 30,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
)

func Get(url string, headers map[string]string) (*Response, error) {
	var location string
	var respbody string
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36")
	if err != nil {
		return &Response{"999", 999, "", nil, 0, "", ""}, err
	}
	for v, k := range headers {
		req.Header[v] = []string{k}
	}
	resp, err := client.Do(req)
	return &Response{resp.Status, resp.StatusCode, respbody, resp.Header, len(respbody), resp.Request.URL.String(), location}, nil
	//newBody, err := io.ReadAll(resp.Body)
	//return string(newBody), nil
}

func post(Url string, contentType string, data interface{}, headers map[string]string) (*Response, error) {
	var (
		payload io.Reader
		err     error
		req     *http.Request
	)
	var location string
	var respbody string
	switch contentType {
	case "application/json":
		var jsonData []byte
		jsonData, err = json.Marshal(data)
		if err != nil {
			return &Response{"999", 999, "", nil, 0, "", ""}, err
		}
		payload = bytes.NewBuffer(jsonData)
		req, err = http.NewRequest("POST", Url, payload)
		if err != nil {
			return &Response{"999", 999, "", nil, 0, "", ""}, err
		}
		req.Header.Set("Content-Type", "application/json")
	case "application/x-www-form-urlencoded":
		values := url.Values{}
		formData, ok := data.(map[string]string)
		if !ok {
			return &Response{"999", 999, "", nil, 0, "", ""}, errors.New("invalid form data")
		}
		for key, val := range formData {
			values.Add(key, val)
		}
		payload = strings.NewReader(values.Encode())
		req, err = http.NewRequest("POST", Url, payload)
		if err != nil {
			return &Response{"999", 999, "", nil, 0, "", ""}, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	case "multipart/form-data":
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		formData, ok := data.(map[string]string)
		if !ok {
			return &Response{"999", 999, "", nil, 0, "", ""}, errors.New("invalid form data")
		}
		for key, val := range formData {
			partHeader := textproto.MIMEHeader{}
			partHeader.Set("Content-Disposition", `form-data; name="`+key+`"`)
			part, err := writer.CreatePart(partHeader)
			if err != nil {
				return &Response{"999", 999, "", nil, 0, "", ""}, err
			}
			_, err = part.Write([]byte(val))
			if err != nil {
				return &Response{"999", 999, "", nil, 0, "", ""}, err
			}
		}
		uploadData, ok := data.(map[string]*os.File)
		if !ok {
			return &Response{"999", 999, "", nil, 0, "", ""}, errors.New("invalid file data")
		}
		for key, val := range uploadData {
			part, err := writer.CreateFormFile(key, filepath.Base(val.Name()))
			if err != nil {
				return &Response{"999", 999, "", nil, 0, "", ""}, err
			}
			_, err = io.Copy(part, val)
			if err != nil {
				return &Response{"999", 999, "", nil, 0, "", ""}, err
			}
		}
		writer.Close()
		req, err = http.NewRequest("POST", Url, body)
		if err != nil {
			return &Response{"999", 999, "", nil, 0, "", ""}, err
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
	default:
		return &Response{"999", 999, "", nil, 0, "", ""}, errors.New("unsupported content type")
	}

	if err != nil {
		return &Response{"999", 999, "", nil, 0, "", ""}, err
	}
	for v, k := range headers {
		req.Header[v] = []string{k}
	}
	resp, err := client.Do(req)
	//newBody, err := io.ReadAll(resp.Body)
	return &Response{resp.Status, resp.StatusCode, respbody, resp.Header, len(respbody), resp.Request.URL.String(), location}, nil
}

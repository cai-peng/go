package owl

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	POST string = "POST"
	GET  string = "GET"
)

type Header struct {
	Key, Value string
}

type HttpRequest struct {
	Url     string
	Method  string
	Headers []Header
	Data    *url.Values
}

func (h *HttpRequest) Request() (string, error) {
	request, err := http.NewRequest(h.Method, h.Url, bytes.NewBufferString(h.Data.Encode()))
	if err != nil {
		return "", err
	}

	for _, i := range h.Headers {
		request.Header.Add(i.Key, i.Value)
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	body, _ := ioutil.ReadAll(response.Body)

	defer response.Body.Close()
	if response.StatusCode == 200 {
		return string(body), nil
	} else {
		return "", errors.New("StatusCode:" + strconv.Itoa(response.StatusCode))
	}
}

type BodyRequest struct {
	Url      string
	Method   string
	Headers  []Header
	BodyData []byte
}

func (r *BodyRequest) Request() ([]byte, error) {
	request, err := http.NewRequest(r.Method, r.Url, bytes.NewReader(r.BodyData))
	if err != nil {
		return []byte(""), err
	}

	for _, i := range r.Headers {
		request.Header.Add(i.Key, i.Value)
	}
	var respone *http.Response
	respone, err = http.DefaultClient.Do(request)
	if err != nil {
		return []byte(""), err
	}
	defer respone.Body.Close()
	b, err := ioutil.ReadAll(respone.Body)
	return b, err
}

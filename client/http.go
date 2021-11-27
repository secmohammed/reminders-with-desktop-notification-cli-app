package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// reminderBody represents reminder request body
type reminderBody struct {
	ID       string        `json:"id"`
	Title    string        `json:"title"`
	Message  string        `json:"message"`
	Duration time.Duration `json:"duration"`
}

type HTTPClient struct {
	client     *http.Client
	BackendURI string
}

func NewHTTPClient(uri string) *HTTPClient {
	return &HTTPClient{
		client:     &http.Client{},
		BackendURI: uri,
	}
}

func (h HTTPClient) Create(title, message string, duration time.Duration) ([]byte, error) {
	requestBody := reminderBody{
		Title:    title,
		Message:  message,
		Duration: duration,
	}
	return h.apiCall(
		http.MethodPost,
		"/reminders",
		&requestBody,
		http.StatusCreated,
	)
}

func (h HTTPClient) Edit(id, title, message string, duration time.Duration) ([]byte, error) {
	requestBody := reminderBody{
		ID:       id,
		Title:    title,
		Message:  message,
		Duration: duration,
	}
	return h.apiCall(
		http.MethodPatch,
		"/reminders/"+id,
		&requestBody,
		http.StatusOK,
	)
}

func (h HTTPClient) Fetch(ids []string) ([]byte, error) {
	idsSet := strings.Join(ids, ",")
	return h.apiCall(
		http.MethodGet,
		"/reminders/"+idsSet,
		nil,
		http.StatusOK,
	)
}

func (h HTTPClient) Delete(ids []string) error {
	idsSet := strings.Join(ids, ",")
	_, err := h.apiCall(
		http.MethodDelete,
		"/reminders/"+idsSet,
		nil,
		http.StatusNoContent,
	)
	return err
}

func (h HTTPClient) Healthy(host string) bool {
	fmt.Println("Checking health of backend at: ", host)
	res, err := http.Get(host + "/health")
	if err != nil || res.StatusCode != http.StatusOK {
		return false
	}
	return true
}

// apiCall makes a new backend api call
func (c HTTPClient) apiCall(method, path string, body interface{}, resCode int) ([]byte, error) {
	bs, err := json.Marshal(body)
	if err != nil {
		e := wrapError("could not marshal request body", err)
		return nil, e
	}
	req, err := http.NewRequest(method, c.BackendURI+path, bytes.NewReader(bs))
	if err != nil {
		e := wrapError("could not create request", err)
		return []byte{}, e
	}

	res, err := c.client.Do(req)
	if err != nil {
		e := wrapError("could not make http call", err)
		return []byte{}, e
	}

	resBody, err := c.readResBody(res.Body)
	if err != nil {
		return []byte{}, err
	}
	if res.StatusCode != resCode {
		if len(resBody) > 0 {
			fmt.Printf("got this response body:\n%s\n", resBody)
		}
		return []byte{}, fmt.Errorf(
			"expected response code: %d, got: %d",
			resCode,
			res.StatusCode,
		)
	}

	return []byte(resBody), err
}

// readBody reads response body
func (h HTTPClient) readResBody(b io.Reader) (string, error) {
	bs, err := ioutil.ReadAll(b)
	if err != nil {
		return "", wrapError("could not read response body", err)
	}

	if len(bs) == 0 {
		return "", nil
	}

	buff := bytes.NewBuffer(bs)
	if _, err := json.MarshalIndent(buff, "", "\t"); err != nil {
		return "", wrapError("could not indent json", err)
	}

	return buff.String(), nil
}

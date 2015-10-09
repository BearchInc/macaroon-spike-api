package gcm

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var (
	Endpoint = "https://gcm-http.googleapis.com/gcm/send"
)

type Params map[string]interface{}

type Message struct {
	To               string `json:"to"`
	ContentAvailable bool   `json:"content_available"`
	Notification     Params `json:"notification"`
}

type Sender struct {
	key  string
	http *http.Client
}

func NewSender(key string) *Sender {
	return NewSenderWithHttpClient(key, http.DefaultClient)
}

func NewSenderWithHttpClient(key string, client *http.Client) *Sender {
	return &Sender{
		key:  key,
		http: client,
	}
}

func (gcm *Sender) Send(message Message) error {
	b, _ := json.Marshal(message)

	req, _ := http.NewRequest("POST", Endpoint, bytes.NewReader(b))
	req.Header.Set("Content-type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", "key="+gcm.key)
	res, err := gcm.http.Do(req)

	var response Params
	json.NewDecoder(res.Body).Decode(&response)
	if response["failure"].(float64) > 0 {
		log.Printf("Error Sending Push Notification: %+v", response)
		return errors.New("Error: Could not send push notification.")
	}
	return err
}

package acme

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/square/go-jose"
)

// Client represents the ACME client
type Client struct {
	ServerURL *url.URL
	Signer    jose.Signer
}

// NewClient creates a new instance of Client with the specified URL as the remote server
func NewClient(baseURL string) *Client {
	URL, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}

	client := Client{
		ServerURL: URL,
	}

	return &client
}

func (client *Client) request(method string, path string, obj interface{}) *http.Response {
	data, err := json.Marshal(obj)
	if err != nil {
		log.Fatal(err)
	}

	data = client.sign(data)
	reader := bytes.NewReader(data)

	request, err := http.NewRequest(method, "", reader)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.URL = client.ServerURL
	request.URL.Path = path

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	return response
}

func (client *Client) sign(data []byte) []byte {
	if client.Signer == nil {
		return data
	}

	signature, err := client.Signer.Sign(data)
	if err != nil {
		log.Fatal(err)
	}

	payload := signature.FullSerialize()

	return []byte(payload)
}

package acme

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// ClientRegistration is returned from a registration attempt
type ClientRegistration struct {
	Key           string
	RecoveryToken string
	Contact       []string
}

// Register the Client with the remote server
func (client *Client) Register(contact []string) *ClientRegistration {
	payload := map[string]interface{}{
		"contact": contact,
	}

	res := client.request("POST", "/acme/new-registration", payload)
	defer res.Body.Close()

	if res.StatusCode != 201 {
		log.Fatal("Registration failed")
	}

	return newClientRegistrationFromResponse(res)
}

func newClientRegistrationFromResponse(res *http.Response) *ClientRegistration {
	buf := new(bytes.Buffer)

	if _, err := buf.ReadFrom(res.Body); err != nil {
		log.Fatal(err)
	}

	registration := new(ClientRegistration)
	json.Unmarshal(buf.Bytes(), &registration)

	return registration
}

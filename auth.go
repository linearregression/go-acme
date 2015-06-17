package acme

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// AuthorizationRequest describes the request for an ACME identifier
type AuthorizationRequest struct {
	Status       AuthorizationStatus
	Expires      time.Time
	Identifier   Identity
	Challenges   []AuthorizationChallenge
	Combinations []authorizationChallengeCombinations
}

// AuthorizationStatus indicates the status of the authorization request
type AuthorizationStatus string

// AuthorizationStatuses
const (
	AuthStatusUnknown    AuthorizationStatus = "unknown"
	AuthStatusPending                        = "pending"
	AuthStatusProcessing                     = "processing"
	AuthStatusValid                          = "valid"
	AuthStatusInvalid                        = "invalid"
	AuthStatusRevoked                        = "revoked"
)

// Identity describes the identifier that should be authorized
type Identity struct {
	Type  string
	Value string
}

// AuthorizationChallenge describes the method for the client to verify the identity
type AuthorizationChallenge interface {
}

type authorizationChallengeCombinations []int

// Authorize starts the verification process for the given identifier
func (client *Client) Authorize(identifier *Identity) *AuthorizationRequest {
	payload := map[string]interface{}{
		"identifier": identifier,
	}

	res := client.request("POST", "/acme/new-authorization", payload)
	defer res.Body.Close()

	if res.StatusCode != 201 {
		log.Fatal("Authorization request failed")
	}

	return newAuthorizationRequestFromResponse(res)
}

func newAuthorizationRequestFromResponse(res *http.Response) *AuthorizationRequest {
	buf := new(bytes.Buffer)

	if _, err := buf.ReadFrom(res.Body); err != nil {
		log.Fatal(err)
	}

	request := new(AuthorizationRequest)
	json.Unmarshal(buf.Bytes(), &request)

	return request
}

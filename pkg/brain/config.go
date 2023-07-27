package brain

import (
	"net/http"
	"time"
)

// The global accessible initialised client
var Client *ClientWithResponses

func InitClient(serverUrl string) error {
	c := http.Client{Timeout: time.Duration(90) * time.Second}
	newClient, err := NewClientWithResponses(serverUrl, WithHTTPClient(&c))

	if err != nil {
		return err
	}

	Client = newClient
	return nil
}

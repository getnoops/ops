package brain

// The global accessible initialised client
var Client *ClientWithResponses

func InitClient(serverUrl string) error {
	newClient, err := NewClientWithResponses(serverUrl)
	if err != nil {
		return err
	}

	Client = newClient
	return nil
}

package types

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type NuvlaResource struct {
	Uuid         string
	ResourceType string
	Data         map[string]interface{}
}

func NewResourceFromResponse(resp *http.Response) *NuvlaResource {
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading response body: %s", err)
		return nil
	}
	defer resp.Body.Close()

	// Unmarshal response body into NuvlaResource
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Errorf("Error unmarshaling response body: %s", err)
		return nil
	}

	// Return NuvlaResource
	return &NuvlaResource{
		Data: data,
	}

}

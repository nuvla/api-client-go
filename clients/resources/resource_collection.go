package resources

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type NuvlaResourceCollection struct {
	Resources    []map[string]interface{} `json:"resources"`
	Count        int                      `json:"count"`
	ResourceName string                   `json:"id"`
}

// NewCollectionFromResponse creates a NuvlaResourceCollection from a http.Response. It expects the body
// of the response to be a CimiCollection where key "resources" contains the list of resources.
func NewCollectionFromResponse(response *http.Response) (*NuvlaResourceCollection, error) {
	// Create NuvlaResourceCollection from response

	// Unmarshal response body into NuvlaResourceCollection
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		log.Errorf("Error reading response body: %s", err)
		return nil, err
	}

	collection := &NuvlaResourceCollection{}
	err = json.Unmarshal(body, collection)
	if err != nil {
		log.Errorf("Error unmarshaling response body: %s", err)
		return nil, err
	}

	return collection, nil
}

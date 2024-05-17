package clients

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

func MarshalResourceIntoMap(resource interface{}, resMap map[string]interface{}) error {
	var mapRes map[string]interface{}
	b, err := json.Marshal(resource)
	if err != nil {
		log.Error("Error marshaling resource to bytes")
		return err
	}

	err = json.Unmarshal(b, &mapRes)
	if err != nil {
		log.Error("Error unmarshalling resource to map")
		return err
	}

	log.Debug("Successfully marshaled resource into map")
	return nil
}

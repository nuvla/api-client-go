package types

import (
	log "github.com/sirupsen/logrus"
	"strings"
)

type NuvlaID struct {
	Uuid         string
	ResourceType string
	Id           string
}

func (n *NuvlaID) String() string {
	return n.Id
}

func NewNuvlaID(uuid string, resourceType string) *NuvlaID {
	return &NuvlaID{
		Uuid:         uuid,
		ResourceType: resourceType,
		Id:           resourceType + "/" + uuid,
	}
}

func NewNuvlaIDFromId(id string) *NuvlaID {
	d := strings.Split(id, "/")
	if len(d) != 2 {
		log.Errorf("Invalid Nuvla ID: %s", id)
		return nil
	}
	return &NuvlaID{
		Id:           id,
		Uuid:         d[1],
		ResourceType: d[0],
	}
}

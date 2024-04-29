package types

import (
	log "github.com/sirupsen/logrus"
	"strings"
)

type InvalidNuvlaID struct {
	Id string
}

func (i *InvalidNuvlaID) Error() string {
	return "Invalid Nuvla ID: " + i.Id
}

type EmptyNuvlaID struct{}

func (e *EmptyNuvlaID) Error() string {
	return "Empty Nuvla ID"
}

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
	if id == "" {
		log.Warnf("Empty Nuvla ID")
		// If empty string, return an empty NuvlaID to prevent NullPointerExceptions
		return &NuvlaID{}
	}

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

package resourceClients

import (
	"api-client-go/client"
	"api-client-go/client/types"
)

type NuvlaEdge struct {
	apiClient *client.NuvlaClient

	nuvlaEdgeId *types.NuvlaID
}

package api_client_go

type NuvlaEdgeClient struct {
	apiClient *NuvlaClient

	NuvlaEdgeId       *NuvlaID
	NuvlaEdgeStatusId *NuvlaID
	CredentialId      *NuvlaID
}

//func NewNuvlaEdgeClient(nuvlaEdgeId *api_client_go.NuvlaID, endpoint string, insecure bool, debug bool) *NuvlaEdgeClient {
//	return &NuvlaEdgeClient{
//		NuvlaEdgeId: nuvlaEdgeId,
//		apiClient: api_client_go.NewNuvlaClient(&api_client_go.SessionOptions{
//			Endpoint: endpoint,
//			Insecure: insecure,
//			Debug:    debug,
//		}),
//	}
//}

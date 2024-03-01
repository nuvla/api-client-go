package api_client_go

type UserClient struct {
	Client *NuvlaClient
}

func NewUserClient(endpoint string, insecure bool, debug bool) *UserClient {
	return &UserClient{
		Client: NewNuvlaClient(endpoint, insecure, debug),
	}
}

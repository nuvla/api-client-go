package api_client_go

type RequestOpts struct {
	Method   string
	Endpoint string
	Data     interface{}
	JsonData map[string]interface{}
	Params   *RequestParams
	Headers  map[string]string
}

type RequestParams struct {
	Select []string
}

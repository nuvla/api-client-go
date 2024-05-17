package types

type RequestOpts struct {
	Method   string
	Endpoint string
	Data     map[string]string
	JsonData map[string]interface{}
	Params   *RequestParams
	Headers  map[string]string
}

type RequestParams struct {
	Select []string
}

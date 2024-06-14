package types

type RequestOpts struct {
	Method   string
	Endpoint string
	Data     map[string]interface{}
	JsonData interface{}
	Params   *RequestParams
	Headers  map[string]string
	Bulk     bool
}

type RequestParams struct {
	Select []string
}

package api_client_go

type SessionOptions struct {
	Endpoint       string `json:"endpoint"`
	Insecure       bool   `json:"insecure"`
	ReAuthenticate bool   `json:"re-authenticate"`
	PersistCookie  bool   `json:"persist-cookie"`
	CookieFile     string `json:"cookie-file"`
	AuthHeader     string `json:"auth-header"`
	Debug          bool   `json:"debug"`
}

func NewSessionOpts(opts *SessionOptions) *SessionOptions {
	if opts == nil {
		opts = &SessionOptions{}
	}

	if opts.Endpoint == "" {
		opts.Endpoint = DefaultEndpoint
	}
	// Insecure is already false as default since declared boolean variables are initialised to 0
	// The same applies to ReAuthenticate, persist Cookie and Debug.

	if opts.PersistCookie && opts.CookieFile == "" {
		opts.CookieFile = DefaultCookieFile
	}

	return opts
}

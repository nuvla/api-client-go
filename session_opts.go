package api_client_go

import (
	"github.com/nuvla/api-client-go/types"
)

type SessionOptFunc func(*SessionOptions)

type SessionOptions struct {
	Endpoint       string `json:"endpoint"`
	Insecure       bool   `json:"insecure"`
	ReAuthenticate bool   `json:"re-authenticate"`
	PersistCookie  bool   `json:"persist-cookie"`
	CookieFile     string `json:"cookie-file"`
	AuthHeader     string `json:"auth-header"`
	Compress       bool   `json:"compress"`
	Debug          bool   `json:"debug"`
}

func DefaultSessionOpts() *SessionOptions {
	return &SessionOptions{
		Endpoint:       types.DefaultEndpoint,
		Insecure:       false,
		ReAuthenticate: false,
		PersistCookie:  true,
		CookieFile:     types.DefaultCookieFile,
		AuthHeader:     "",
		Compress:       true,
		Debug:          false,
	}
}

func WithEndpoint(endpoint string) SessionOptFunc {
	return func(opts *SessionOptions) {
		opts.Endpoint = endpoint
	}
}

func WithDebugSession(flag bool) SessionOptFunc {
	return func(opts *SessionOptions) {
		opts.Debug = flag
	}
}

func WithInsecureSession(flag bool) SessionOptFunc {
	return func(opts *SessionOptions) {
		opts.Insecure = flag
	}
}

func ReAuthenticateSession(opts *SessionOptions) {
	opts.ReAuthenticate = true
}

func WithoutPersistCookie(opts *SessionOptions) {
	opts.PersistCookie = false
}

func WithCookieFile(cookieFile string) SessionOptFunc {
	return func(opts *SessionOptions) {
		opts.CookieFile = cookieFile
	}
}

func WithAuthHeader(authHeader string) SessionOptFunc {
	return func(opts *SessionOptions) {
		opts.AuthHeader = authHeader
	}
}

func WithOutCompressSession(opts *SessionOptions) {
	opts.Compress = false
}

func NewSessionOpts(opts *SessionOptions) *SessionOptions {
	if opts == nil {
		opts = &SessionOptions{}
	}

	if opts.Endpoint == "" {
		opts.Endpoint = types.DefaultEndpoint
	}
	// Insecure is already false as default since declared boolean variables are initialised to 0
	// The same applies to ReAuthenticate, persist Cookie and Debug.
	if opts.PersistCookie && opts.CookieFile == "" {
		opts.CookieFile = types.DefaultCookieFile
	}

	return opts
}

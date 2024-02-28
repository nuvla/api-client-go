package types

type SessionAttributes struct {
	Endpoint       string `json:"endpoint"`
	Insecure       bool   `json:"insecure"`
	ReAuthenticate bool   `json:"re-authenticate"`
	PersistCookie  bool   `json:"persist-cookie"`
	CookieFile     string `json:"cookie-file"`
	AuthHeader     string `json:"auth-header"`
	Debug          bool   `json:"debug"`
}

func NewSessionOpts() {

}

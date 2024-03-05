package api_client_go

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type NuvlaSession struct {
	endpoint       string
	reauthenticate bool
	persistCookie  bool
	loginParams    map[string]string
	authnHeader    string
	compress       bool
	debug          bool

	session *http.Client

	// Nuvla session data
	cookies *NuvlaCookies
}

func NewNuvlaSession(sessionAttrs *SessionOptions) *NuvlaSession {

	log.Infof("Creating new Nuvla session for endpoint %s", sessionAttrs.Endpoint)
	s := &NuvlaSession{
		endpoint:       sessionAttrs.Endpoint,
		reauthenticate: sessionAttrs.ReAuthenticate,
		persistCookie:  sessionAttrs.PersistCookie,
		authnHeader:    sessionAttrs.AuthHeader,
		debug:          sessionAttrs.Debug,
		session: &http.Client{
			Timeout: time.Second * DefaultTimeout,
			Jar:     nil,
		},
	}

	if sessionAttrs.Insecure {
		s.session.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	// Try import jar
	if sessionAttrs.PersistCookie {
		s.cookies = NewNuvlaCookies(sessionAttrs.CookieFile, sessionAttrs.Endpoint)
		s.session.Jar = s.cookies.jar
	}
	// Probably, check here if jar are GOOD

	return s
}

/****************************************************************************************
************************ Credentials management **********************************************
****************************************************************************************/

func (s *NuvlaSession) NeedToLogin() bool {
	return false
}

func (s *NuvlaSession) login(loginParams LogInParams) error {
	// Build headers for login
	h := make(map[string]string)
	h["Content-Type"] = "application/json"
	h["Accept"] = "application/json"

	// Build parameters from interface (Could be either password or API key)
	p := make(map[string]interface{})
	p["template"] = loginParams.GetParams()

	// Send request
	resp, err := s.Request(&RequestOpts{
		Method:   "POST",
		Endpoint: s.endpoint + SessionEndpoint,
		JsonData: p,
		Headers:  h,
	})
	log.Infof("Login response: %s", resp)

	return err
}

/****************************************************************************************
************************ Request management **********************************************
****************************************************************************************/

func (s *NuvlaSession) request(req *http.Request) (*http.Response, error) {
	if s.authnHeader != "" {
		req.Header.Add("nuvla-authn-info", s.authnHeader)
	}

	resp, err := s.session.Do(req)
	if err != nil {
		log.Errorf("Error executing request: %s", err)
		return nil, err
	}

	return resp, nil
}

func addParamsToQuery(req *http.Request, input *RequestParams) {
	if input.Select != nil {
		q := req.URL.Query()
		for _, f := range input.Select {
			q.Add("select", f)
		}
		req.URL.RawQuery = q.Encode()
	}

}

func compressPayload(payload []byte) *bytes.Buffer {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(payload); err != nil {
		log.Warn("Error compressing payload, returning uncompressed payload")
		return bytes.NewBuffer(payload)
	}
	if err := gz.Close(); err != nil {
		log.Errorf("Error closing gzip writer: %s", err)
	}
	return &buf
}

func (s *NuvlaSession) Request(reqInput *RequestOpts) (*http.Response, error) {
	// Build endpoint
	log.Infof("Requesting %s", reqInput.Endpoint)

	r, err := http.NewRequest(reqInput.Method, reqInput.Endpoint, nil)
	if err != nil {
		log.Errorf("Error creating request: %s", err)
		return nil, err
	}
	// Add payload if needed
	if reqInput.JsonData != nil {
		jsonPayload, err := json.Marshal(reqInput.JsonData)
		if err != nil {
			log.Errorf("Error marshalling payload: %s", err)
			return nil, err
		}
		var buffer *bytes.Buffer
		if s.compress {
			buffer = compressPayload(jsonPayload)
		} else {
			buffer = bytes.NewBuffer(jsonPayload)
		}

		r.Body = io.NopCloser(buffer)
	}

	for k, v := range reqInput.Headers {
		log.Debugf("Adding header %s: %s", k, v)
		r.Header.Set(k, v)
	}

	if reqInput.Params != nil {
		addParamsToQuery(r, reqInput.Params)
	}

	resp, err := s.request(r)
	if err != nil {
		log.Errorf("Error executing request: %s", err)
		return nil, err
	}

	if s.persistCookie && resp.Header.Get("Set-Cookie") != "" {
		// Save new jar
		err := s.cookies.SaveIfNeeded(s.session.Jar)
		if err != nil {
			log.Errorf("Error saving jar: %s", err)
		}
	}
	return resp, nil
}

func (s *NuvlaSession) logout() error {
	log.Infof("Logging out from %s", s.endpoint)
	// TODO: Implement me
	// Delete current session
	// Remove cookie
	return nil
}

/****************************************************************************************
************************ Generic utils **********************************************
****************************************************************************************/

func (s *NuvlaSession) String() string {
	return "Nuvla session for endpoint " + s.endpoint
}

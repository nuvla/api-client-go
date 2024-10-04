package api_client_go

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/nuvla/api-client-go/common"
	"github.com/nuvla/api-client-go/types"
	log "github.com/sirupsen/logrus"
	"github.com/wI2L/jsondiff"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type NuvlaSession struct {
	endpoint       string
	insecure       bool
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

func SanitiseEndpoint(endpoint string) string {
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		return endpoint
	} else {
		log.Infof("Endpoint %s does not have a protocol. Assuming https", endpoint)
		return "https://" + endpoint
	}
}

func NewNuvlaSession(sessionAttrs *SessionOptions) *NuvlaSession {

	log.Debugf("Creating new Nuvla session for endpoint %s", sessionAttrs.Endpoint)
	s := &NuvlaSession{
		endpoint:       SanitiseEndpoint(sessionAttrs.Endpoint),
		reauthenticate: sessionAttrs.ReAuthenticate,
		persistCookie:  sessionAttrs.PersistCookie,
		authnHeader:    sessionAttrs.AuthHeader,
		debug:          sessionAttrs.Debug,
		session: &http.Client{
			Timeout: time.Second * types.DefaultTimeout,
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
	} else {
		j, _ := cookiejar.New(nil)
		s.session.Jar = j
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

func (s *NuvlaSession) login(loginParams types.LogInParams) error {
	// Build headers for login
	h := make(map[string]string)
	h["Content-Type"] = "application/json"
	h["Accept"] = "application/json"

	// Build parameters from interface (Could be either password or API key)
	p := make(map[string]interface{})
	p["template"] = loginParams.GetParams()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(common.DefaultRequestTimeout)*time.Second)
	defer cancel()

	// Send request
	log.Debug("Sending login request...")
	res, err := s.Request(ctx, &types.RequestOpts{
		Method:   "POST",
		Endpoint: s.endpoint + types.SessionEndpoint,
		JsonData: p,
		Headers:  h,
	})

	if err != nil {
		log.Errorf("Error logging in: %s", err)
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		log.Errorf("Error logging in: %s", res.Status)
		return fmt.Errorf("error logging in: %s", res.Status)
	}

	return nil
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

func addParamsToQuery(req *http.Request, input *types.RequestParams) {
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

// bodyTypeCompatible checks if the body content is compatible with the body type. Currently supported types are:
// - map[string]interface{}
// - []map[string]interface{}
func bodyTypeCompatible(bodyContent interface{}) bool {
	switch bodyContent.(type) {
	case map[string]interface{}:
		return true
	case []map[string]interface{}:
		return true
	case jsondiff.Patch:
		return true
	default:
		return false
	}
}

func encodeBody(request *http.Request, reqInput *types.RequestOpts, compress bool) error {
	if reqInput.JsonData == nil && reqInput.Data == nil {
		return nil
	}

	if reqInput.JsonData != nil && reqInput.Data != nil {
		log.Warn("Both Data and JsonData provided, this could lead to unexpected behavior. Using JsonData")
	}

	if reqInput.JsonData != nil {
		if !bodyTypeCompatible(reqInput.JsonData) {
			log.Warnf("Unknown type %T for json payload", reqInput.JsonData)
			return nil
		}

		jsonPayload, err := json.Marshal(reqInput.JsonData)
		if err != nil {
			log.Errorf("Error marshalling json payload: %s", err)
			return err
		}

		var buffer *bytes.Buffer
		if compress {
			buffer = compressPayload(jsonPayload)
		} else {
			buffer = bytes.NewBuffer(jsonPayload)
		}
		request.Header.Set("Content-Type", "application/json")
		request.Body = io.NopCloser(buffer)
		request.ContentLength = int64(buffer.Len())
	}

	if reqInput.Data != nil {
		log.Debug("Encoding data payload")
		data := url.Values{}
		for k, value := range reqInput.Data {
			switch v := value.(type) {
			case []string:
				for _, s := range v {
					data.Add(k, s)
				}
			case string:
				data.Add(k, v)
			case int:
				data.Add(k, strconv.Itoa(v))
			default:
				log.Warnf("Unknown type %T for key %s", v, k)
				data.Add(k, fmt.Sprintf("%v", v))
			}
		}
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		request.Body = io.NopCloser(bytes.NewBufferString(data.Encode()))
	}
	return nil
}

func (s *NuvlaSession) Request(ctx context.Context, reqInput *types.RequestOpts) (*http.Response, error) {
	// Build endpoint
	log.Debugf("Sending [%s] request to endpoint: %s", reqInput.Method, reqInput.Endpoint)

	r, err := http.NewRequestWithContext(ctx, reqInput.Method, reqInput.Endpoint, nil)
	if err != nil {
		log.Errorf("Error creating request: %s", err)
		return nil, err
	}

	// Encode body asserting from json or data encoded as URL
	err = encodeBody(r, reqInput, s.compress)
	if err != nil {
		log.Errorf("Error encoding body: %s", err)
		return nil, err
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
	// TODO: Remove cookies if present
	// For the moment, lets release unused connections...
	s.session.CloseIdleConnections()
	return nil
}

/****************************************************************************************
************************ Generic utils **********************************************
****************************************************************************************/

func (s *NuvlaSession) String() string {
	return "Nuvla session for endpoint " + s.endpoint
}

func (s *NuvlaSession) GetSessionOpts() SessionOptions {
	// Fill all session opts
	opts := SessionOptions{
		Endpoint:       s.endpoint,
		Insecure:       s.insecure,
		ReAuthenticate: s.reauthenticate,
		AuthHeader:     s.authnHeader,
		Debug:          s.debug,
		Compress:       s.compress,
	}
	if s.persistCookie && s.cookies != nil {
		opts.PersistCookie = s.persistCookie
		opts.CookieFile = s.cookies.cookieFile
	}

	return opts
}

package types

import "encoding/json"

const (
	HrefSessionTemplateApiKey   = "session-template/api-key"
	HrefSessionTemplatePassword = "session-template/password"
)

type LogInParams interface {
	GetParams() map[string]string
}

type ApiKeyLogInParams struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
	Href   string `json:"href"`
}

func NewApiKeyLogInParams(key, secret string) *ApiKeyLogInParams {
	return &ApiKeyLogInParams{
		Key:    key,
		Secret: secret,
		Href:   HrefSessionTemplateApiKey,
	}
}

func (p *ApiKeyLogInParams) GetParams() map[string]string {
	params := map[string]string{
		"href":   HrefSessionTemplateApiKey,
		"key":    p.Key,
		"secret": p.Secret,
	}

	return params
}

type UserLogInParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Href     string `json:"href"`
}

func NewUserLogInParams(username, password string) *UserLogInParams {
	return &UserLogInParams{
		Username: username,
		Password: password,
		Href:     HrefSessionTemplatePassword,
	}
}

func (p *UserLogInParams) GetParams() map[string]string {
	var params map[string]string
	jsonParams, _ := json.Marshal(p)

	err := json.Unmarshal(jsonParams, &params)
	if err != nil {
		return nil
	}
	return params
}

package types

const DefaultEndpoint = "https://nuvla.io"
const DefaultInsecure = false

// Default path locations
const (
	DefaultConfigPath = "/tmp/.nuvla/"

	DefaultCookieFile = DefaultConfigPath + ".jar"
	SessionPath       = DefaultConfigPath + ".session"
)

// SessionEndpoint Default endpoints
const (
	SessionEndpoint = "/api/session"
)

// DefaultTimeout
// Network defaults
const (
	DefaultTimeout = 10
)

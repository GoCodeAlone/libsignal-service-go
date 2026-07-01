package sandbox

import (
	"errors"
	"strings"

	"github.com/GoCodeAlone/libsignal-service-go/fake"
	"github.com/GoCodeAlone/libsignal-service-go/service"
)

var (
	ErrMissingEndpoint  = errors.New("missing sandbox endpoint")
	ErrOfficialEndpoint = errors.New("official-looking endpoint forbidden in sandbox")
)

type Config struct {
	Endpoint                     string
	AllowOfficialEndpointForTest bool
}

func NewAdapter(cfg Config) (*service.Adapter, error) {
	if cfg.Endpoint == "" {
		return nil, ErrMissingEndpoint
	}
	if !cfg.AllowOfficialEndpointForTest && officialLookingEndpoint(cfg.Endpoint) {
		return nil, ErrOfficialEndpoint
	}
	return service.NewAdapter(fake.NewAdapter(), service.AdapterConfig{
		Mode:     service.AdapterModeSandbox,
		Endpoint: cfg.Endpoint,
	})
}

func officialLookingEndpoint(endpoint string) bool {
	normalized := strings.ToLower(endpoint)
	for _, pattern := range []string{
		"signal" + ".org",
		"text" + "secure",
		"whisper" + "systems",
		"chat." + "signal",
	} {
		if strings.Contains(normalized, pattern) {
			return true
		}
	}
	return false
}

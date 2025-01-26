package validators

import (
	"errors"
	"net/http"
	"net/url"
)

// URLValidator defines an interface for URL validation
type URLValidator interface {
	Validate(targetURL string) error
}

// DefaultURLValidator implements URL validation logic
type DefaultURLValidator struct{}

// NewURLValidator creates a new instance of DefaultURLValidator
func NewURLValidator() URLValidator {
	return &DefaultURLValidator{}
}

var (
	// ErrMissingURL indicates that the URL parameter is required
	ErrMissingURL = errors.New("URL parameter is required")

	// ErrInvalidURLFormat indicates that the URL is not in a valid format
	ErrInvalidURLFormat = errors.New("invalid URL format, please provide a valid URL")

	// ErrURLNotReachable indicates that the URL could not be reached
	ErrURLNotReachable = errors.New("URL is not reachable, please provide a valid URL")

	// ErrNon200StatusCode indicates that the URL returned a non-200 HTTP status code
	ErrNon200StatusCode = errors.New("URL returned non-200 status")
)

// Validate checks if the given URL is valid and reachable
func (v *DefaultURLValidator) Validate(targetURL string) error {
	_, err := url.ParseRequestURI(targetURL)
	if err != nil {
		return ErrInvalidURLFormat
	}

	resp, err := http.Get(targetURL)
	if err != nil {
		return ErrURLNotReachable
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrNon200StatusCode
	}

	return nil
}
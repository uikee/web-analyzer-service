package validators

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestValidate_SuccessfulURL(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com").
		Reply(200)

	validator := NewURLValidator()
	err := validator.Validate("http://example.com")

	assert.NoError(t, err, "Expected no error for a reachable valid URL")
}

func TestValidate_InvalidURLFormat(t *testing.T) {
	validator := NewURLValidator()
	err := validator.Validate("invalid-url")

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidURLFormat, err, "Expected ErrInvalidURLFormat for malformed URL")
}

func TestValidate_URLNotReachable(t *testing.T) {
	defer gock.Off()

	gock.New("http://unreachable.com").
		ReplyError(ErrURLNotReachable)

	validator := NewURLValidator()
	err := validator.Validate("http://unreachable.com")

	assert.Error(t, err)
	assert.Equal(t, ErrURLNotReachable, err, "Expected ErrURLNotReachable for an unreachable URL")
}

func TestValidate_Non200StatusCode(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com").
		Reply(404)

	validator := NewURLValidator()
	err := validator.Validate("http://example.com")

	assert.Error(t, err)
	assert.Equal(t, ErrNon200StatusCode, err, "Expected ErrNon200StatusCode for non-200 response")
}
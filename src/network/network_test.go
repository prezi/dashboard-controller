package network

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestAddProtocolAndPortToIP(t *testing.T) {
	assert.Equal(t, "http://10.0.0.126:1234", AddProtocolAndPortToIP("10.0.0.126", "1234"))
}

func TestErrorHandler(t *testing.T) {
	assert.Equal(t, false, ErrorHandler(nil, "This is an error message."))
	err := errors.New("This is an error message.")
	assert.Equal(t, true, ErrorHandler(err, "%v"))
}

func TestGetOS(t *testing.T) {
	OS := GetOS()
	assert.IsType(t, "Some OS Name", OS)
}

func TestCreateFormWithInitialValues(t *testing.T) {
	urlToDisplay := "some valid url"
	form := CreateFormWithInitialValues(map[string]string{"url": urlToDisplay})
	assert.Equal(t, form, url.Values{"url": []string{"some valid url"}})
}

func TestGetRelativeFilePath(t *testing.T) {
	filepath := GetRelativeFilePath("assets/images")
	assert.IsType(t, "some/filepath", filepath)
}

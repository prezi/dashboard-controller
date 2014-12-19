package network

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestErrorHandler(t *testing.T) {
	assert.Equal(t, false, ErrorHandler(nil, "This is an error message."))
	err := errors.New("This is an error message.")
	assert.Equal(t, true, ErrorHandler(err, "%v"))
}

func TestCreateFormWithInitialValues(t *testing.T) {
	urlToDisplay := "some valid url"
	form := CreateFormWithInitialValues(map[string]string{"url": urlToDisplay})
	assert.Equal(t, form, url.Values{"url": []string{"some valid url"}})
}

func TestGetOS(t *testing.T) {
	assert.IsType(t, "Some OS Name", GetOS())
}

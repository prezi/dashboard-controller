package network

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"strings"
)

const (
	TEST_URL = "http://httpstatusdogs.com"
)

func TestErrorHandler(t *testing.T) {
	assert.Equal(t, false, ErrorHandler(nil, "This is an error message."))
	err := errors.New("This is an error message.")
	assert.Equal(t, true, ErrorHandler(err, "%v"))
}

func TestCreateFormWithInitialValues(t *testing.T) {
	form := CreateFormWithInitialValues(map[string]string{"url": TEST_URL})
	assert.Equal(t, form, url.Values{"url": []string{TEST_URL}})
}

func TestGetOS(t *testing.T) {
	assert.IsType(t, "Some OS Name", GetOS())
}

func TestGetProjectRoot(t *testing.T) {
	projectPath := getProjectRoot()
	projectName := "dashboard-controller"
	assert.True(t, strings.Contains(projectPath, projectName))
}

package network

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"regexp"
)

func TestGetIPAddress(t *testing.T) {
	IPAddress := getIPAddress("4003")
	IPAddressRegexpPattern := "([0-9]*\\.){3}[0-9]*:[0-9]*"
	re := regexp.MustCompile(IPAddressRegexpPattern)

	assert.Equal(t, true, re.MatchString(IPAddress))
}

func TestGetIPAddressFromCmdLine(t *testing.T) {
	IPAddress := getIPAddressBytesFromCmdLine()
	temp_string := ""

	assert.IsType(t, temp_string, IPAddress)
}

func TestParseIpAddress(t *testing.T) {
	parsedIPAddress := parseIpAddress(getIPAddressBytesFromCmdLine())
	temp_string := ""

	assert.IsType(t, temp_string, parsedIPAddress)
}

func TestAddProtocolAndPortToIp(t *testing.T) {
	assert.Equal(t, "http://10.0.0.126:1234", addProtocolAndPortToIp("10.0.0.126", "1234"))
}

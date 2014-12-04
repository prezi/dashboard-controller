package hash

import (
	// "fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestIsHashMatchInUserAuthenticationMap(t *testing.T) {
	hashMap := map[string][16]byte{
		"lilo": CreateHashFromString("secret"),
	}
	result := IsHashMatchInUserAuthenticationMap("lilo", "secret", hashMap)
	assert.Equal(t, true, result)
}

func TestInitializeUserAuthenticationMap(t *testing.T) {
	result := InitializeUserAuthenticationMap()
	expectedResult := map[string][16]byte{"lilo": CreateHashFromString("poke")}
	assert.Equal(t, true, reflect.DeepEqual(result, expectedResult))
}

func TestGetUserNameAndPasswordFromFile(t *testing.T) {
	content := "lilo\npoke"
	username, password := GetUserNameAndPasswordFromFile(content)
	assert.Equal(t, username, "lilo")
	assert.Equal(t, password, "poke")
}

package master

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRelativeFilePath(t *testing.T) {
	filepath := GetRelativeFilePath("assets/images")
	assert.IsType(t, "some/filepath", filepath)

}

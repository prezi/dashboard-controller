package slave

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetUp(t *testing.T) {
	port, slaveName, masterURL, OS := SetUp()

	assert.Equal(t, DEFAULT_LOCALHOST_PORT, port)
	assert.Equal(t, "SLAVE NAME UNSPECIFIED", slaveName)
	assert.Equal(t, "http://localhost:5000", masterURL)
	assert.IsType(t, "Some OS Name", OS)
}

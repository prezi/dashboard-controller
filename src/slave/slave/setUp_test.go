package slave

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSetUp(t *testing.T) {
	port, slaveName, masterURL, OS := SetUp()
	assert.Equal(t, "8080", port)
	assert.Equal(t, "SLAVE NAME UNSPECIFIED", slaveName)
	assert.Equal(t, "http://localhost:5000", masterURL)
	assert.IsType(t, "Some OS Name", OS)
}

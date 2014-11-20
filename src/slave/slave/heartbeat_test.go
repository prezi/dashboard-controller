package slave

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHeartbeat(t *testing.T) {
	err := Heartbeat(1, "slaveName", "someMasterAddress:5000")
	assert.NotNil(t, err)
}

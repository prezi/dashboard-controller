package slave

import (
	"net/http"
	"network"
	"time"
)

func Heartbeat(heartbeatInterval int, slaveName, ownPort, masterURL string) (err error) {
	err = nil
	masterURLForHeartbeat := masterURL + "/receive_heartbeat"
	beat := time.Tick(time.Duration(heartbeatInterval) * time.Second)

	client := &http.Client{}
	form := network.CreateFormWithInitialValues(map[string]string{"slaveName": slaveName, "slavePort": ownPort})

	for _ = range beat {
		_, err = client.PostForm(masterURLForHeartbeat, form)
		network.ErrorHandler(err, "Error communicating with master: %v\n")
	}
	return
}

package proxy

import (
	"net/http"
	"network"
	"time"
)

func Heartbeat(heartbeatInterval int, masterURL string) (err error) {
	err = nil
	masterURLForHeartbeat := masterURL + "/receive_proxy_heartbeat"
	beat := time.Tick(time.Duration(heartbeatInterval) * time.Second)

	client := &http.Client{}
	form := network.CreateFormWithInitialValues(map[string]string{})

	for _ = range beat {
		_, err = client.PostForm(masterURLForHeartbeat, form)
		network.ErrorHandler(err, "Error communicating with master: %v\n")
	}
	return
}

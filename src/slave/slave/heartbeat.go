package slave

import (
	"time"
	"net/http"
	"net/url"
	"network"
)

func Heartbeat(heartbeatInterval int, slaveName, masterURL string) (err error) {
	masterURLForHeartbeat := masterURL + "/receive_heartbeat"
	beat := time.Tick(time.Duration(heartbeatInterval) * time.Second)
	
	client := &http.Client{}
	form := url.Values{}
	form.Set("slaveName", slaveName)

    for now := range beat {
		form.Set("heartbeatTimestamp", now.String())
		_, err = client.PostForm(masterURLForHeartbeat, form)

		network.ErrorHandler(err, "Error communicating with master: %v\n")
    }
    return nil
}

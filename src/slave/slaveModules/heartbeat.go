package slaveModule

import (
	"time"
	"fmt"
	"net/http"
	"net/url"
	// "os"
)

func Heartbeat(heartbeatInterval int, slaveName string, masterIP string) {
	beat := time.Tick(time.Duration(heartbeatInterval) * time.Second)
    for now := range beat {
       	client := &http.Client{}
		form := url.Values{}
		form.Set("slaveName", slaveName)
		form.Set("heartbeatTimestamp", now.String())
		fmt.Println("hearbeat", now.String())

		masterIPAddressForHeartbeat := "http://localhost:5000/receive_heartbeat"
		_, err := client.PostForm(masterIPAddressForHeartbeat, form)

		if err != nil {
			fmt.Printf("Error communicating with master: %v\n", err)
			fmt.Println("Aborting program.")
			// os.Exit(1)
		}
    }
}

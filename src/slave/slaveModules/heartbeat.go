package slaveModule

import (
	"time"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	// "os"
)

func Heartbeat(heartbeatInterval int, slaveName string, masterIP string) (err error) {
	masterIPAddressForHeartbeat := getMasterReceiveHeartbeatAddress(masterIP)
	beat := time.Tick(time.Duration(heartbeatInterval) * time.Second)
	
	client := &http.Client{}
	form := url.Values{}
	form.Set("slaveName", slaveName)

    for now := range beat {
		form.Set("heartbeatTimestamp", now.String())
		_, err = client.PostForm(masterIPAddressForHeartbeat, form)

		if err != nil {
			fmt.Printf("Error communicating with master: %v\n", err)
			fmt.Println("Aborting program.")
			return err
			// os.Exit(1)
		}
    }
    return nil
}

func getMasterReceiveHeartbeatAddress(masterIP string) (masterAddress string) {
	masterIPAddressAndExtentionArray := []string{"http://", masterIP, "/receive_heartbeat"} 
	return strings.Join(masterIPAddressAndExtentionArray, "")
}

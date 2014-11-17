package slaveModule

import (
	"time"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func Heartbeat(slaveName string, masterIP string) {

	beat := time.Tick(1 * time.Second)
    for now := range beat {
       	client := &http.Client{}
		form := url.Values{}
		form.Set("slaveName", slaveName)
		// form.Set("slaveIPAddress", slaveIPAddress)
		form.Set("heartbeat timestamp", now.String())
		fmt.Println("hearbeat")

		masterIPAddressForHeartbeat := "http://localhost:5000"
		_, err := client.PostForm(masterIPAddressForHeartbeat, form)

		if err != nil {
			fmt.Printf("Error communicating with master: %v\n", err)
			fmt.Println("Aborting program.")
			os.Exit(1)
		}
    }
}

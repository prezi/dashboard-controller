package master

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

var test_mode = false

func ReceiveSlaveHeartbeat(request *http.Request, slaveMap map[string]Slave) {
	slaveName, slaveAddress := processRequest(request)

	if _, existsInMap := slaveMap[slaveName]; existsInMap {
		updateSlaveHeartbeat(slaveMap, slaveAddress, slaveName)
	} else {
		fmt.Printf("Slave added with name \"%v\", IP %v", slaveName, slaveAddress)
		slaveMap[slaveName] = Slave{URL: slaveAddress, heartbeat: time.Now()}
		sendSlaveListToWebserver(webServerAddress, slaveMap)
	}
}

func processRequest(request *http.Request) (slaveName, slaveAddress string) {
	slaveName = request.PostFormValue("slaveName")
	slavePort := request.PostFormValue("slavePort")

	slaveIP, _, _ := net.SplitHostPort(request.RemoteAddr)
	slaveAddress = "http://" + slaveIP + ":" + slavePort
	return
}

func updateSlaveHeartbeat(slaveMap map[string]Slave, slaveAddress, slaveName string) ( err error){
	slaveInstance := slaveMap[slaveName]
	if slaveInstance.URL != slaveAddress {
		fmt.Printf(`WARNING: Slave with name \"%v\"
			already exists with the IP address: %v. \n
			kill signal sent to slave with name \"%v\"
			with IP address: %v`,
			slaveName, slaveInstance.URL, slaveName, slaveAddress)
		err = sendKillSignalToSlave(slaveAddress)
	} else {
		slaveInstance.heartbeat = time.Now()
		slaveMap[slaveName] = slaveInstance
	}
	return
}

func sendKillSignalToSlave(slaveAddress string) (err error) {
	client := &http.Client{}
	form := url.Values{}
	form.Set("message", "die")
	_, err = client.PostForm(slaveAddress+"/receive_killsignal", form)
	return
}

func MonitorSlaves(timeInterval int, slaveMap map[string]Slave) {
	timer := time.Tick(time.Duration(timeInterval) * time.Second)
	for _ = range timer {
		removeDeadSlaves(timeInterval, slaveMap)
		if test_mode {
			break
		}
	}
}

func removeDeadSlaves(deadTime int, slaveMap map[string]Slave) {
	slavesToRemove := make([]string,0,len(slaveMap))
	remainingSlaves := make([]string,0,len(slaveMap))
	for slaveName, slave := range slaveMap {
		timeDifference := time.Now().Sub(slave.heartbeat)
		timeThreshold := time.Duration(deadTime)*time.Second
		if timeDifference > timeThreshold {
			slavesToRemove = append(slavesToRemove, slaveName)
		} else {
			remainingSlaves = append(remainingSlaves,slaveName)
		}
	}
	if len(slavesToRemove) > 0 {
		fmt.Printf("\nREMOVING DEAD SLAVES: %v\n", slavesToRemove)
		for _, deadSlaveName := range slavesToRemove {
			delete(slaveMap, deadSlaveName)
		}
		sendSlaveListToWebserver(webServerAddress, slaveMap)
		if len(remainingSlaves) >0 {
			fmt.Printf("Current slaves are: %v\n",remainingSlaves)
		} else {
			fmt.Println("No slaves available.\n")
		}
	}

}

func UpdateWebserverAddress(r *http.Request)(err error) {
	newWebServerAddress, err := getWebserverAddress(r)
	if webServerAddress != newWebServerAddress {
		fmt.Printf("Webserver address has changed from %v to %v\n", webServerAddress, newWebServerAddress)
		webServerAddress = newWebServerAddress
	}
	return
}


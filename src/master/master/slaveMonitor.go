package master

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

var test_mode = false

type Slave struct {
	URL          string
	heartbeat    time.Time
	displayedURL string // TODO: store currently displayed URL for each slave
}

func SetUp() (slaveMap map[string]Slave) {
	slaveMap = make(map[string]Slave)
	return
}

func ReceiveSlaveHeartbeat(request *http.Request, slaveMap map[string]Slave) {
	slaveName, slaveAddress := processSlaveHeartbeatRequest(request)

	if _, existsInMap := slaveMap[slaveName]; existsInMap {
		updateSlaveHeartbeat(slaveMap, slaveAddress, slaveName)
	} else {
		fmt.Printf("Slave added with name \"%v\", URL %v.\n\n", slaveName, slaveAddress)
		slaveMap[slaveName] = Slave{URL: slaveAddress, heartbeat: time.Now()}
		sendSlaveListToWebserver(webServerAddress, slaveMap)
	}
}

func processSlaveHeartbeatRequest(request *http.Request) (slaveName, slaveAddress string) {
	slaveName = request.PostFormValue("slaveName")
	slavePort := request.PostFormValue("slavePort")

	slaveIP, _, _ := net.SplitHostPort(request.RemoteAddr)
	slaveAddress = "http://" + slaveIP + ":" + slavePort
	return
}

func updateSlaveHeartbeat(slaveMap map[string]Slave, slaveAddress, slaveName string) (err error) {
	slaveInstance := slaveMap[slaveName]
	if slaveInstance.URL != slaveAddress {
		fmt.Println("WARNING: Received signal from slave with duplicate name.")
		fmt.Printf("Slave with name \"%v\" already exists.\n", slaveName)
		fmt.Printf("Sending kill signal to duplicate slave at URL %v.\n\n", slaveAddress)
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
	slavesToRemove := getDeadSlaves(deadTime, slaveMap)
	if len(slavesToRemove) > 0 {
		fmt.Printf("\nREMOVING DEAD SLAVES: %v\n", slavesToRemove)
		for _, deadSlaveName := range slavesToRemove {
			delete(slaveMap, deadSlaveName)
		}
		printSlaveNamesInMap(slaveMap)
		sendSlaveListToWebserver(webServerAddress, slaveMap)
	}
}

func getDeadSlaves(deadTime int, slaveMap map[string]Slave) (deadSlaves []string) {
	for slaveName, slave := range slaveMap {
		timeDifference := time.Now().Sub(slave.heartbeat)
		timeThreshold := time.Duration(deadTime) * time.Second

		if timeDifference > timeThreshold {
			deadSlaves = append(deadSlaves, slaveName)
		}
	}
	return
}

func printSlaveNamesInMap(slaveMap map[string]Slave) {
	fmt.Println("Current slaves are: ")
	if len(slaveMap) == 0 {
		fmt.Println("No slaves available.")
	} else {
		for slaveName, _ := range slaveMap {
			fmt.Println(slaveName)
		}
	}
}

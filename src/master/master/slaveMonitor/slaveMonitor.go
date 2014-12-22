package slaveMonitor

import (
	"fmt"
	"master/master"
	"net"
	"net/http"
	"network"
	"proxy/manageIPTables"
	"strings"
	"time"
)

// TODO: Check if proxyIP flag is given by user.
// Assume there is a proxy if user specifies proxyIP. If not, assume no proxy.
// Alternatively, have proxy send heartbeats to master. Then IS_USING_PROXY will be a var.

const IS_USING_PROXY = true

func ReceiveSlaveHeartbeat(request *http.Request, slaveMap map[string]master.Slave) (updatedSlaveMap map[string]master.Slave) {
	slaveName, slaveAddress := processSlaveHeartbeatRequest(request)

	if _, existsInMap := slaveMap[slaveName]; existsInMap {
		updateSlaveHeartbeat(slaveMap, slaveAddress, slaveName)
	} else {
		addNewSlaveToMap(slaveMap, slaveAddress, slaveName)
		if IS_USING_PROXY {
			manageIPTables.AddNewSlaveToIPTables(splitProtocolAndPortFromIP(slaveAddress))
		}
	}
	return slaveMap
}

func processSlaveHeartbeatRequest(request *http.Request) (slaveName, slaveAddress string) {
	slaveName = request.PostFormValue("slaveName")
	slavePort := request.PostFormValue("slavePort")
	slaveIP, _, _ := net.SplitHostPort(request.RemoteAddr)
	slaveAddress = "http://" + slaveIP + ":" + slavePort
	return
}

func updateSlaveHeartbeat(slaveMap map[string]master.Slave, slaveAddress, slaveName string) {
	slaveInstance := slaveMap[slaveName]
	if slaveInstance.URL != slaveAddress {
		killDuplicateSlave(slaveName, slaveAddress)
	} else {
		slaveInstance.Heartbeat = time.Now()
		slaveMap[slaveName] = slaveInstance
	}
}

func killDuplicateSlave(slaveName, slaveAddress string) {
	fmt.Println("WARNING: Received signal from slave with duplicate name.")
	fmt.Printf("Slave with name \"%v\" already exists.\n", slaveName)
	fmt.Printf("Sending kill signal to duplicate slave at URL %v.\n\n", slaveAddress)
	err := sendKillSignalToSlave(slaveAddress)
	network.ErrorHandler(err, "Error encountered killing slave: %v\n")
}

func sendKillSignalToSlave(slaveAddress string) (err error) {
	client := &http.Client{}
	form := network.CreateFormWithInitialValues(map[string]string{"message": "die"})
	_, err = client.PostForm(slaveAddress+"/receive_killsignal", form)
	return
}

func addNewSlaveToMap(slaveMap map[string]master.Slave, slaveAddress, slaveName string) {
	fmt.Printf("Slave added with name \"%v\", URL %v.\n\n", slaveName, slaveAddress)
	slaveMap[slaveName] = master.Slave{URL: slaveAddress, Heartbeat: time.Now(), PreviouslyDisplayedURL: "http://google.com", DisplayedURL: "http://google.com"}
	fmt.Println(slaveMap[slaveName])
}

func splitProtocolAndPortFromIP(address string) (ip string) {
	host, _, _ := net.SplitHostPort(address)
	ip = strings.TrimPrefix(host, "http://")
	return
}

func MonitorSlaves(timeInterval int, slaveMap map[string]master.Slave) {
	timer := time.Tick(time.Duration(timeInterval) * time.Second)
	for _ = range timer {
		removeDeadSlaves(timeInterval, slaveMap)
	}
}

func removeDeadSlaves(deadTime int, slaveMap map[string]master.Slave) {
	slavesToRemove := getDeadSlaves(deadTime, slaveMap)
	if len(slavesToRemove) > 0 {
		fmt.Printf("\nREMOVING DEAD SLAVES: %v\n", slavesToRemove)
		for _, deadSlaveName := range slavesToRemove {
			if IS_USING_PROXY {
				manageIPTables.RemoveDeadSlaveFromIPTables(splitProtocolAndPortFromIP(slaveMap[deadSlaveName].URL))
			}
			delete(slaveMap, deadSlaveName)
		}
		printSlaveNamesInMap(slaveMap)
	}
}

func getDeadSlaves(deadTime int, slaveMap map[string]master.Slave) (deadSlaves []string) {
	for slaveName, slave := range slaveMap {
		timeDifference := time.Now().Sub(slave.Heartbeat)
		timeThreshold := time.Duration(deadTime) * time.Second

		if timeDifference > timeThreshold {
			deadSlaves = append(deadSlaves, slaveName)
		}
	}
	return
}

func printSlaveNamesInMap(slaveMap map[string]master.Slave) {
	fmt.Println("Current slaves are: ")
	if len(slaveMap) == 0 {
		fmt.Println("No slaves available.\n")
	} else {
		for slaveName, _ := range slaveMap {
			fmt.Println(slaveName)
		}
	}
}

func ListSlaveNames(slaveMap map[string]master.Slave) (slaveNames []string) {
	slaveNames = make([]string, 0, len(slaveMap))
	for k := range slaveMap {
		slaveNames = append(slaveNames, k)
	}
	return slaveNames
}

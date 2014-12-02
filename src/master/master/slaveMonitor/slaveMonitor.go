package slaveMonitor

import (
	"fmt"
	"master/master"
	"master/master/webserverCommunication"
	"net"
	"net/http"
	"network"
	"time"
)

var test_mode = false

func ReceiveSlaveHeartbeat(request *http.Request, slaveMap map[string]master.Slave, webServerAddress string) (updatedSlaveMap map[string]master.Slave) {
	slaveName, slaveAddress := processSlaveHeartbeatRequest(request)

	if _, existsInMap := slaveMap[slaveName]; existsInMap {
		updateSlaveHeartbeat(slaveMap, slaveAddress, slaveName)
	} else {
		addNewSlaveToMap(slaveMap, slaveAddress, slaveName, webServerAddress)
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

func addNewSlaveToMap(slaveMap map[string]master.Slave, slaveAddress, slaveName, webServerAddress string) {
	fmt.Printf("Slave added with name \"%v\", URL %v.\n\n", slaveName, slaveAddress)
	slaveMap[slaveName] = master.Slave{URL: slaveAddress, Heartbeat: time.Now()}
	webserverCommunication.SendSlaveListToWebserver(webServerAddress, slaveMap)
}

func MonitorSlaves(timeInterval int, slaveMap map[string]master.Slave, webServerAddress string) {
	timer := time.Tick(time.Duration(timeInterval) * time.Second)
	for _ = range timer {
		removeDeadSlaves(timeInterval, slaveMap, webServerAddress)
		if test_mode {
			break
		}
	}
}

func removeDeadSlaves(deadTime int, slaveMap map[string]master.Slave, webServerAddress string) {
	slavesToRemove := getDeadSlaves(deadTime, slaveMap)
	if len(slavesToRemove) > 0 {
		fmt.Printf("\nREMOVING DEAD SLAVES: %v\n", slavesToRemove)
		for _, deadSlaveName := range slavesToRemove {
			delete(slaveMap, deadSlaveName)
		}
		printSlaveNamesInMap(slaveMap)
		webserverCommunication.SendSlaveListToWebserver(webServerAddress, slaveMap)
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

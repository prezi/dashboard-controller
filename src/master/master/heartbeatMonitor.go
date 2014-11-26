package master

import (
	"fmt"
	"net"
	"net/http"
	"time"
	"strconv"
)

func MonitorSlaveHeartbeats(_ http.ResponseWriter, request *http.Request, slaveMap map[string]Slave) {
	slaveName, slaveAddress := processRequest(request)

	if _, existsInMap := slaveMap[slaveName]; existsInMap {
		slaveMap = updateSlaveHeartbeat(slaveMap, slaveAddress, slaveName)
	} else {
		fmt.Printf("Slave added with name \"%v\", IP %v", slaveName, slaveAddress)
		slaveMap[slaveName] = Slave{URL: slaveAddress, heartbeat: time.Now()}
		sendSlaveListToWebserver(webserverAddress, slaveMap)
	}
}

func processRequest(request *http.Request) (slaveName, slaveAddress string) {
	slaveName = request.PostFormValue("slaveName")
	slavePort := request.PostFormValue("slavePort")

	slaveIP,_,_ := net.SplitHostPort(request.RemoteAddr)
	slaveAddress = "http://" + slaveIP + ":" + slavePort
	return
}

func updateSlaveHeartbeat(slaveMap map[string]Slave, slaveAddress, slaveName string) map[string]Slave {
	slaveInstance := slaveMap[slaveName]
	if slaveInstance.URL != slaveAddress {
		newSlaveName := getNewSlaveName(slaveMap,slaveName)
		slaveMap[newSlaveName] = Slave{URL: slaveAddress, heartbeat: time.Now()}
		fmt.Printf(`WARNING: Slave with name \"%v\"
			already exists with the IP address: %v. \n
			New slave with IP %v added with name %v\n`,
			slaveName, slaveInstance.URL, newSlaveName, slaveAddress)
	} else {
		slaveInstance.heartbeat = time.Now()
		slaveMap[slaveName] = slaveInstance
	}
	return slaveMap
}

func getNewSlaveName(slaveMap map[string]Slave,slaveName string) (newSlaveName string) {
	for number:=2; ;number++ {
		newSlaveName = slaveName + "_"+strconv.Itoa(number)
		fmt.Println(newSlaveName)
		_, ok := slaveMap[slaveName]
		if !ok {
			return
		}
//		if _, existInMap := slaveMap[slaveName];  !existInMap {
//			return
//		}
	}
}

func MonitorSlaves(timeInterval int, slaveMap map[string]Slave) {
	timer := time.Tick(time.Duration(timeInterval) * time.Second)
	for _ = range timer {
		removeDeadSlaves(timeInterval, slaveMap)
	}
}

func removeDeadSlaves(deadTime int, slaveMap map[string]Slave) {
	for slaveName, slave := range slaveMap {
		if time.Now().Sub(slave.heartbeat) > time.Duration(deadTime)*time.Second {
			fmt.Printf("\nREMOVING DEAD SLAVE: %v\n", slaveName)
			delete(slaveMap, slaveName)
			fmt.Println("Current slaves are: ")
			if len(slaveMap) == 0 {
				fmt.Println("No slaves available.")
			} else {
				for slaveName, _ := range slaveMap {
					fmt.Println(slaveName)
				}
			}
			fmt.Printf("\n\n")
			sendSlaveListToWebserver(webserverAddress, slaveMap)
		}
	}
}

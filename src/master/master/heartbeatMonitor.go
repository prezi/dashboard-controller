package master

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

func MonitorSlaveHeartbeats(_ http.ResponseWriter, request *http.Request, slaveMap map[string]Slave) {
	slaveName, slaveAddress := processRequest(request)

	if slaveInstance, existsInMap := slaveMap[slaveName]; existsInMap {
		slaveMap[slaveName] = updateSlaveHeartbeat(slaveInstance, slaveAddress, slaveName)
	} else {
		fmt.Printf("Slave added with name %v, IP %v", slaveName, slaveAddress)
		slaveMap[slaveName] = Slave{URL: slaveAddress, heartbeat: time.Now()}
		sendSlaveToWebserver(webserverAddress, slaveMap)
	}
}

func processRequest(request *http.Request) (slaveName, slaveAddress string) {
	slaveName = request.PostFormValue("slaveName")
	slavePort := request.PostFormValue("slavePort")

	slaveIP,_,_ := net.SplitHostPort(request.RemoteAddr)
	slaveAddress = "http://" + slaveIP + ":" + slavePort
	return
}

func updateSlaveHeartbeat(slaveInstance Slave, slaveAddress, slaveName string) Slave {
	if slaveInstance.URL != slaveAddress {
		slaveInstance.URL = slaveAddress
		fmt.Printf(`WARNING: Slave with name \"%v\" 
			already exists with the IP address: %v. \n 
			Updating %v's IP address to %v.\n`,
			slaveName, slaveInstance.URL, slaveName, slaveAddress)
	}
	slaveInstance.heartbeat = time.Now()
	return slaveInstance
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
			fmt.Printf("\ntime elapsed since last update: %v", time.Now().Sub(slave.heartbeat))
			fmt.Printf("\nREMOVING DEAD SLAVE: %v\n", slaveName)
			delete(slaveMap, slaveName)
			fmt.Println("Updated Slave Map: ")
			fmt.Println("Valid slave IDs are: ")
			for slaveName, _ := range slaveMap {
				fmt.Println(slaveName)
			}
			fmt.Printf("\n\n")
			sendSlaveToWebserver(webserverAddress, slaveMap)
		}
	}
}

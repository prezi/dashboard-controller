package master

import (
	"net/http"
	"fmt"
	"time"
	"net"
)

func MonitorSlaveHeartbeats(_ http.ResponseWriter, request *http.Request, slaveMap map[string]Slave) {
	slaveName := request.PostFormValue("slaveName")
	slavePort := request.PostFormValue("slavePort")
	slaveIP,_,_ := net.SplitHostPort(request.RemoteAddr)

	slaveAddress := slaveIP+":"+slavePort

	if returnedSlave, existsInMap := slaveMap[slaveName]; existsInMap {
		if returnedSlave.URL == slaveAddress {
			slaveToUpdate := slaveMap[slaveName]
			slaveToUpdate.heartbeat = time.Now()
			slaveMap[slaveName] = slaveToUpdate
		} else {
			slaveMap[slaveName] = Slave{URL: slaveAddress, heartbeat: time.Now()}
			fmt.Printf(`WARNING: Slave with name \"%v\" 
				already exists with the IP address: %v. \n 
				Updating %v's IP address to %v.\n`, 
				slaveName, returnedSlave.URL, slaveName, slaveAddress)
		}
	} else {
		for keySlaveName, valueSlave := range slaveMap {
			if valueSlave.URL == slaveAddress {
				delete(slaveMap, keySlaveName)
				fmt.Printf("WARNING: The following slave will be removed due IP conflict: %v",keySlaveName)
			}
		}
		fmt.Printf("Slave added with name %v, IP %v",slaveName,slaveAddress)
		slaveMap[slaveName] = Slave{URL: slaveAddress, heartbeat: time.Now()}
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
		if time.Now().Sub(slave.heartbeat) > time.Duration(deadTime) * time.Second {
			fmt.Printf("time elapsed since last update: %v",time.Now().Sub(slave.heartbeat))
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

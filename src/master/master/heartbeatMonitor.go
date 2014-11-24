package master

import (
	"net/http"
	"fmt"
	"time"
)

func MonitorSlaveHeartbeats(_ http.ResponseWriter, request *http.Request, slaveMap map[string]Slave) {
	slaveName := request.PostFormValue("slaveName")
	heartbeatTimestamp := time.Now()

	if returnedSlave, existsInMap := slaveMap[slaveName]; existsInMap {
		returnedSlave.heartbeat = heartbeatTimestamp
		slaveMap[slaveName] = returnedSlave
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

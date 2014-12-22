package manageIPTables

import (
	"fmt"
	"network"
	"os/exec"
	"proxy/proxy"
)

func AddNewSlaveToIPTables(slaveIP string) (err error) {
	if proxy.OS == "Linux" {
		err = exec.Command("sudo", "iptables", "-A", "INPUT", "-s", slaveIP, "-j", "ACCEPT", "-m", "tcp", "-p", "tcp", "--dport", proxy.PROXY_PORT).Run()
		if !network.ErrorHandler(err, "Error adding slave to iptables: %v\n") {
			fmt.Println("Slave added to proxy IP tables.")
		}
	}
	return
}

func RemoveDeadSlaveFromIPTables(slaveIP string) (err error) {
	if proxy.OS == "Linux" {
		err = exec.Command("sudo", "iptables", "-D", "INPUT", "-s", slaveIP, "-j", "ACCEPT", "-m", "tcp", "-p", "tcp", "--dport", proxy.PROXY_PORT).Run()
		if !network.ErrorHandler(err, "Error deleting slave from iptables: %v\n") {
			fmt.Println("Slave deleted from iptables.")
		}
	}
	return
}

package manageIPTables

import (
	"fmt"
	"net/http"
	"network"
	"os/exec"
	"proxy/proxy"
)

func UpdateIPTables(request *http.Request) {
	IPAddressToAdd := request.PostFormValue("IPAddressToAdd")
	IPAddressToDelete := request.PostFormValue("IPAddressToDelete")
	if IPAddressToAdd != "" {
		addNewSlaveToIPTables(IPAddressToAdd)
	}
	if IPAddressToDelete != "" {
		removeDeadSlaveFromIPTables(IPAddressToDelete)
	}
}

func addNewSlaveToIPTables(slaveIP string) (err error) {
	if proxy.OS == "Linux" {
		err = exec.Command("sudo", "iptables", "-A", "INPUT", "-s", slaveIP, "-j", "ACCEPT", "-m", "tcp", "-p", "tcp", "--dport", proxy.PROXY_PORT).Run()
		if !network.ErrorHandler(err, "Error adding slave to iptables: %v\n") {
			fmt.Println("Slave added to proxy IP tables.")
		}
	}
	return
}

func removeDeadSlaveFromIPTables(slaveIP string) (err error) {
	if proxy.OS == "Linux" {
		err = exec.Command("sudo", "iptables", "-D", "INPUT", "-s", slaveIP, "-j", "ACCEPT", "-m", "tcp", "-p", "tcp", "--dport", proxy.PROXY_PORT).Run()
		if !network.ErrorHandler(err, "Error deleting slave from iptables: %v\n") {
			fmt.Println("Slave deleted from iptables.")
		}
	}
	return
}

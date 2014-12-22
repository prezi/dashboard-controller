package proxy

import (
	"fmt"
	"master/master"
	"network"
	"os/exec"
)

// TODO: proxy independent from master - send requests to proxy's IP address; proxy will run relevant command-line executions

var (
	OS         = network.GetOS() // TODO: make all these switch cases, if we ever program for OS X platform
	PROXY_PORT = master.GetProxyPort()
)

func InitializeIPTables() {
	if OS == "Linux" {
		FlushIPTables()
		AcceptResponseFromDNSServer()
		AcceptRequestsOnMasterPort()
	}
}

func FlushIPTables() (error error) {
	error = exec.Command("sudo", "iptables", "-F").Run()
	if error != nil {
		fmt.Printf("Error flushing iptables: %v\n", error)
	}
	return
}

func AcceptResponseFromDNSServer() (error error) {
	error = exec.Command("sudo", "iptables", "-A", "INPUT", "-m", "conntrack", "--ctstate", "RELATED,ESTABLISHED", "-j", "ACCEPT").Run()
	if error != nil {
		fmt.Printf("Error setting rule for accepting responses from DNS server: %v\n", error)
	}
	return
}

func AcceptRequestsOnMasterPort() (error error) {
	error = exec.Command("sudo", "iptables", "-A", "INPUT", "-j", "ACCEPT", "-m", "tcp", "-p", "tcp", "--dport", "5000").Run()
	if error != nil {
		fmt.Printf("Error setting rule for accepting responses from DNS server: %v\n", error)
	}
	return
}

func AddNewSlaveToIPTables(slaveIP string) (error error) {
	if OS == "Linux" {
		error = exec.Command("sudo", "iptables", "-A", "INPUT", "-s", slaveIP, "-j", "ACCEPT", "-m", "tcp", "-p", "tcp", "--dport", PROXY_PORT).Run()
		if error != nil {
			fmt.Printf("Error adding slave to iptables: %v\n", error)
		}
		fmt.Println("Slave added to proxy IP tables.")
	}
	return
}

func RemoveDeadSlaveFromIPTables(slaveIP string) (error error) {
	if OS == "Linux" {
		error = exec.Command("sudo", "iptables", "-D", "INPUT", "-s", slaveIP, "-j", "ACCEPT", "-m", "tcp", "-p", "tcp", "--dport", PROXY_PORT).Run()
		if error != nil {
			fmt.Printf("Error deleting slave from iptables: %v\n", error)
		}
		fmt.Println("Slave deleted from iptables.")
	}
	return
}

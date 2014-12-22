package proxy

import (
	"flag"
	"fmt"
	"network"
	"os/exec"
)

// TODO: proxy independent from master - send requests to proxy's IP address; proxy will run relevant command-line executions

var (
	OS                       = network.GetOS() // TODO: make all these switch cases, if we ever program for OS X platform
	PROXY_PORT               = "8080"          // mitmproxy runs on port 8080
	PROXY_CONFIGURATION_FILE = network.GetRelativeFilePath("proxyConfig.py")
	DEFAULT_MASTER_IP        = "localhost"
	DEFAULT_MASTER_PORT      = "5000"
)

func SetUp() (masterURL string) {
	startProxy()
	masterIP, masterPort := configFlags()
	initializeIPTables(masterIP)
	return getMasterURL(masterIP, masterPort)
}

func startProxy() {
	fmt.Println("Starting mitmproxy with command: mitmproxy -s ", PROXY_CONFIGURATION_FILE)
	err := exec.Command("mitmproxy", "-s", PROXY_CONFIGURATION_FILE).Run()
	network.ErrorHandler(err, "Error starting mitmproxy: %v\n")
}

func configFlags() (masterIP, masterPort string) {
	flag.StringVar(&masterIP, "masterIP", DEFAULT_MASTER_IP, "master IP")
	flag.StringVar(&masterPort, "masterPort", DEFAULT_MASTER_PORT, "master port")
	flag.Parse()
	return
}

func initializeIPTables(masterIP string) {
	if OS == "Linux" {
		flushIPTables()
		acceptResponseFromDNSServer()
		acceptRequestsFromMaster(masterIP)
	}
}

func flushIPTables() (err error) {
	err = exec.Command("sudo", "iptables", "-F").Run()
	network.ErrorHandler(err, "Error flushing iptables: %v\n")
	return
}

func acceptResponseFromDNSServer() (err error) {
	err = exec.Command("sudo", "iptables", "-A", "INPUT", "-m", "conntrack", "--ctstate", "RELATED,ESTABLISHED", "-j", "ACCEPT").Run()
	network.ErrorHandler(err, "Error setting rule for accepting responses from DNS server: %v\n")
	return
}

func acceptRequestsFromMaster(masterIP string) (err error) {
	err = exec.Command("sudo", "iptables", "-A", "INPUT", "-s", masterIP, "-j", "ACCEPT", "-m", "tcp", "-p", "tcp", "--dport", PROXY_PORT).Run()
	network.ErrorHandler(err, "Error setting rule for accepting requests from master: %v\n")
	return
}

func getMasterURL(masterIP, masterPort string) (masterURL string) {
	return "http://" + masterIP + ":" + masterPort
}

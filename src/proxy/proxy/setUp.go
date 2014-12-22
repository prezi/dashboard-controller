package proxy

import (
	"flag"
	"fmt"
	"net"
	"network"
	"os/exec"
	"strings"
)

// TODO: proxy independent from master - send requests to proxy's IP address; proxy will run relevant command-line executions
var (
	OS                       = network.GetOS() // TODO: make all these switch cases, if we ever program for OS X platform
	PROXY_PORT               = "8080"          // mitmproxy runs on port 8080
	PROXY_CONFIGURATION_FILE = network.GetRelativeFilePath("proxyConfig.py")
	DEFAULT_MASTER_URL       = "http://localhost:5000"
)

func Start() (masterURL string) {
	fmt.Println("Starting mitmproxy with command: mitmproxy -s ", PROXY_CONFIGURATION_FILE)
	err := exec.Command("mitmproxy", "-s", PROXY_CONFIGURATION_FILE).Run()
	network.ErrorHandler(err, "Error starting mitmproxy: %v\n")
	masterURL = getMasterURL()
	initializeIPTables(masterURL)
	return
}

func initializeIPTables(masterURL string) {
	masterIP := splitProtocolAndPortFromIP(masterURL)
	if OS == "Linux" {
		flushIPTables()
		acceptResponseFromDNSServer()
		acceptRequestsFromMaster(masterIP)
	}
}

func getMasterURL() (masterURL string) {
	flag.StringVar(&masterURL, "masterURL", DEFAULT_MASTER_URL, "master URL")
	flag.Parse()
	return
}

func splitProtocolAndPortFromIP(URL string) (IP string) {
	host, _, _ := net.SplitHostPort(URL)
	IP = strings.TrimPrefix(host, "http://")
	return
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

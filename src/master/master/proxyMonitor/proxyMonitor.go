package proxyMonitor

import (
	"fmt"
	"master/master"
	"net"
	"net/http"
	"network"
	"regexp"
	"time"
)

var (
	IS_USING_PROXY   = false
	PROXY_DEAD_TIME  = 0
	PROXY_PORT       string
	PROXY_IP_ADDRESS string
	PROXY_URL        string
)

func ReceiveProxyHeartbeat(request *http.Request, slaveMap map[string]master.Slave) {
	if !IS_USING_PROXY {
		PROXY_IP_ADDRESS = getOriginIPAddressFromRequest(request)
		PROXY_PORT = request.PostFormValue("ProxyHTTPServerPort")
		PROXY_URL = "http://" + PROXY_IP_ADDRESS + ":" + PROXY_PORT
		fmt.Printf("Proxy detected at %v.\n", PROXY_IP_ADDRESS)
		sendCurrentSlaveIPAddressesToProxy(slaveMap)
		IS_USING_PROXY = true
	}
	PROXY_DEAD_TIME = 0
}

func getOriginIPAddressFromRequest(request *http.Request) (proxyIP string) {
	proxyIP, _, _ = net.SplitHostPort(request.RemoteAddr)
	return
}

func sendCurrentSlaveIPAddressesToProxy(slaveMap map[string]master.Slave) {
	IPAddresses := getSlaveIPAddresses(slaveMap)
	sendSlaveIPAddressesToProxy(IPAddresses)
}

func getSlaveIPAddresses(slaveMap map[string]master.Slave) (IPAddresses []string) {
	for key := range slaveMap {
		slave := slaveMap[key]
		IPAddresses = append(IPAddresses, getIPInString(slave.URL))
	}
	return IPAddresses
}

func getIPInString(input string) string {
	numBlock := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	regexPattern := numBlock + "\\." + numBlock + "\\." + numBlock + "(\\." + numBlock + ")?"

	regEx := regexp.MustCompile(regexPattern)
	return regEx.FindString(input)
}

func sendSlaveIPAddressesToProxy(IPAddresses []string) {
	for i := range IPAddresses {
		RequestProxyToAddNewSlaveToIPTables(PROXY_URL, IPAddresses[i])
	}
}

func MonitorProxy(timeInterval int) {
	timer := time.Tick(time.Duration(timeInterval) * time.Second)
	for _ = range timer {
		if IS_USING_PROXY {
			PROXY_DEAD_TIME += 1
			checkProxyHealth(timeInterval)
		}
	}
}

func checkProxyHealth(timeInterval int) {
	if PROXY_DEAD_TIME > timeInterval {
		fmt.Println("Proxy has been disconnected.")
		IS_USING_PROXY = false
	}
}

func RequestProxyToAddNewSlaveToIPTables(proxyURL, IPAddressToAdd string) (err error) {
	proxyURLToUpdateIPTables := proxyURL + "/update_iptables"

	client := &http.Client{}
	form := network.CreateFormWithInitialValues(map[string]string{"IPAddressToAdd": IPAddressToAdd})

	_, err = client.PostForm(proxyURLToUpdateIPTables, form)
	network.ErrorHandler(err, "Error communicating with proxy: %v\n")
	return
}

func RequestProxyToRemoveDeadSlaveFromIPTables(proxyURL, IPAddressToDelete string) (err error) {
	proxyURLToUpdateIPTables := proxyURL + "/update_iptables"

	client := &http.Client{}
	form := network.CreateFormWithInitialValues(map[string]string{"IPAddressToDelete": IPAddressToDelete})

	_, err = client.PostForm(proxyURLToUpdateIPTables, form)
	network.ErrorHandler(err, "Error communicating with proxy: %v\n")
	return
}

package main

import(
	"flag"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	DEFAULT_LOCALHOST_PORT = 4000
)

var port int
var OS string
var err error

func main() {
	setUp()
	http.HandleFunc("/", handleRequest)

	fmt.Printf("Listening on port: %v\n", port)
	fmt.Printf("You can send HTTP POST requests with a 'url' parameter to open it in a browser.\n")
	fmt.Printf("e.g.: curl localhost:%v -X POST -d \"url=http://www.google.com\"\n", port)

	// start HTTP server with given address and handler
	// handler=nil will default handler to DefaultServeMux
	err = http.ListenAndServe(":" + strconv.Itoa(port), nil)
	if err != nil {
		fmt.Printf("Error starting HTTP server: %v\n", err)
		fmt.Println("Abort process.")
	}
}

func setUp() {
	setOS()
	if (OS == "Unknown") {
		fmt.Printf("ERROR: Failed to detect operating system.\n")
		fmt.Println("Abort process.")
	} else {
		fmt.Printf("Operating system detected: %v\n", OS)
	}

	flag.IntVar(&port, "port", DEFAULT_LOCALHOST_PORT, "the port to listen on for commands")
	// can pass flag argument: $ ./slave -port=8080
	// if flag not specified, will set port=DEFAULT_LOCALHOST_PORT
	flag.Parse()
}

func setOS() {
	// func (c *Cmd) Output() ([]byte, error)
	operatingSystemBytes, err := exec.Command("uname", "-a").Output() // display operating system name...why do we need the -a?
	operatingSystemName := string(operatingSystemBytes)

	var kernel string
	// fmt.Println("cmd", operatingSystemName)

	if err != nil {
		fmt.Printf("Error encountered while reading kernal: %v\n", err)
		kernel = "Unknown"
	} else {
		kernel = strings.Split(operatingSystemName, " ")[0]
	}
	fmt.Println("Kernal detected: ", kernel)

	switch kernel {
	case "Linux":
		OS = "Linux"
	case "Darwin":
		OS = "OS X"
	default:
		OS = "Unknown"
	}
}

func killBrowser() {
	switch OS {
	case "Linux":
		fmt.Printf("Executing command: kill -SIGTERM $(pidof chromium)")
		err = exec.Command("kill -SIGTERM $(pidof chromium)").Run() // TODO: this needs testing
	case "OS X":
		fmt.Printf("Executing command: killall 'Google Chrome'\n")
		err = exec.Command("killall", "Google Chrome").Run()
	}

	if err != nil {
		fmt.Printf("Error killing current browser: %v\n", err)
	} else {
	// sleep the code so that the browser can finish closing 
	time.Sleep(1 * time.Second)		
	}
}

func openBrowser(url string){
	switch OS {
	case "Linux":
		fmt.Printf("Executing command: chromium --kiosk %v\n", url)
		err = exec.Command("nohup", "chromium", "--kiosk", url).Run()		
	case "OS X":
		fmt.Printf("Executing command: open -a 'Google Chrome' --args --kiosk %v\n", url)
		err = exec.Command("open", "-a", "Google Chrome", "--args", "--kiosk", url).Run()
	}

	if err != nil {
		fmt.Printf("Error opening URL: %v\n", err)
	}	
}

func handleRequest(writer http.ResponseWriter, request *http.Request) {
	url := request.PostFormValue("url")
	killBrowser()
	openBrowser(url)
}

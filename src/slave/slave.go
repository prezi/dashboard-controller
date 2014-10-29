package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"bytes"
)

const DEFAULT_LOCALHOST_PORT = 8080

var port int
var OS string
var err error

func main() {
	setUp()
	http.HandleFunc("/", handleRequest)

	fmt.Printf("Listening on port: %v\n", port)
	fmt.Println("You can send HTTP POST requests with a 'url' parameter to open it in a browser.")
	fmt.Printf("e.g.: curl localhost:%v -X POST -d \"url=http://www.google.com\"\n", port)

	// start HTTP server with given address and handler
	// handler=nil will default handler to DefaultServeMux
	err = http.ListenAndServe(":" + strconv.Itoa(port), nil)
	if err != nil {
		fmt.Printf("Error starting HTTP server: %v\n", err)
		fmt.Println("ERROR: Failed to start HTTP server.")
		fmt.Println("Aborting program.")
		os.Exit(1)
	}
}

func setUp() {
	setOS()
	if (OS == "Unknown") {
		fmt.Println("ERROR: Failed to detect operating system.")
		fmt.Println("Aborting program.")
		os.Exit(1)
	} else {
		fmt.Printf("Operating system detected: %v\n", OS)
	}
	// can pass flag argument: $ ./slave -port=8080
	// if flag not specified, will set port=DEFAULT_LOCALHOST_PORT
	flag.IntVar(&port, "port", DEFAULT_LOCALHOST_PORT, "the port to listen on for commands")
	flag.Parse()

	// :0.0 indicates the first screen attached to the first display in localhost
	err = os.Setenv("DISPLAY",":0.0")
	if err != nil {
		fmt.Printf("Error setting DISPLAY environment variable: %v\n", err)
	}
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

func blockProgramWhileBrowserCloses() {
	// block the code so that the browser can finish closing 
	// use a while-loop to check if there is a browser process running, 
	// exit the loop once all browser processes have closed

	var existingProcess []byte
	emptyByteArray := make([]byte, 0)
	
	for {
		switch OS {
		case "Linux":
			existingProcess, err = exec.Command("pgrep", "chromium").CombinedOutput()		
		case "OS X":
			existingProcess, err = exec.Command("pgrep", "Google Chrome").CombinedOutput()
		}

		if bytes.Equal(existingProcess, emptyByteArray) { 
				break
		}
	}
}

func killBrowser() {
	switch OS {
	case "Linux":
		fmt.Println("Executing command: killall chromium")
		err = exec.Command("killall", "chromium").Run() 
	case "OS X":
		fmt.Println("Executing command: killall 'Google Chrome'")
		err = exec.Command("killall", "Google Chrome").Run()
	}

	if err != nil {
		fmt.Printf("Error killing current browser: %v\n", err)
	} else {
		blockProgramWhileBrowserCloses()	
	}
}

func openBrowser(url string){
	switch OS {
	case "Linux":
		fmt.Printf("Executing command: chromium --kiosk %v &\n", url)
		err = exec.Command("chromium", "--kiosk", url, "&").Run()		
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

package main

import(
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	DEFAULT_LOCALHOST_PORT = 4000
	DEFAULT_LOG_FILE = "/log/commandListener.log"

	LINUX_DEFAULT_BROWSER_CMD = "chromium"
	LINUX_DEFAULT_BROWSER_ARGS = "--kiosk"

	OSX_DEFAULT_BROWSER_CMD = "open"
	OSX_DEFAULT_BROWSER_ARGS = "-a 'Google Chrome' --args --kiosk"
)

var port int
var browser_cmd string
var browser_args string

func main() {
	setUp()
	http.HandleFunc("/", handleRequest)
	fmt.Printf("listening on port: %v\n", port)
	fmt.Printf("you can send HTTP POST requests with an 'url' parameter to open it in a browser\n")
	fmt.Printf("e.g.: curl localhost:4000 -X POST -d \"url=http://www.google.com\"\n")
	err := http.ListenAndServe(":" + strconv.Itoa(port), nil)
	if err != nil {
		fmt.Printf("error starting HTTP server: %v\n", err)
	}
}

func setUp() {
	OS := getOs()
	fmt.Printf("detected operating system: %v\n", OS)
	switch OS {
	case "Linux":
		browser_cmd = LINUX_DEFAULT_BROWSER_CMD
		browser_args = LINUX_DEFAULT_BROWSER_ARGS
	case "OSX":
		browser_cmd = OSX_DEFAULT_BROWSER_CMD
		browser_args = OSX_DEFAULT_BROWSER_ARGS
	default:
		print("error: unknown operating system")
	}
	flag.IntVar(&port, "port", DEFAULT_LOCALHOST_PORT, "the port to listen on for commands")
	flag.Parse()
	err := os.Setenv("DISPLAY",":0.0")
	if err != nil {
		fmt.Printf("error setting env variable: %v\n", err)
	}
}

func getOs() string {
	cmd := exec.Command( "uname", "-a" )
	var kernel string
	output, err := cmd.Output()
	if( err != nil ) {
		fmt.Printf("getting kernel: %v\n", err)
		kernel = "unknown"
	} else {
		kernel = strings.Split( string(output), " " )[0]
	}
	var OS string
	switch kernel {
	case "Linux":
		OS = "Linux"
	case "Darwin":
		OS = "OSX"
	default:
		OS = "unknown"
	}
	return OS
}

func handleRequest(writer http.ResponseWriter, request *http.Request) {
	url := request.PostFormValue("url")
	fmt.Printf("executing: %v %v %v\n", browser_cmd, browser_args, url)
	err := exec.Command(browser_cmd, browser_args, url).Run()
//	err := exec.Command(browser_cmd, url).Run()
	if err != nil {
		fmt.Printf("error opening URL: %v\n", err)
	}
}

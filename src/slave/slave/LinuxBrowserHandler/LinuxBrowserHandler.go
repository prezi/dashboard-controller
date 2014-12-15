package LinuxBrowserHandler

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

func BrowserHandler(writer http.ResponseWriter, request *http.Request, browserProcess *exec.Cmd) *exec.Cmd {
	url := request.PostFormValue("url")
	killBrowser()
	browserProcess, _ = openBrowser(url)
	fmt.Fprintf(writer, "SUCCESS. \"%v\" has been posted.\n", url)

	return browserProcess
}

func killBrowser() (err error) {
	fmt.Println("Executing command: killall chromium")
	err = exec.Command("killall", "-TERM", "chromium").Run()
	if err != nil {
		fmt.Printf("Error killing current browser: %v\n", err)
	} else {
		blockWhileBrowserCloses()
	}
	return
}

func blockWhileBrowserCloses() (err error) {
	var existingProcess []byte
	for {
		time.Sleep(75 * time.Millisecond)
		existingProcess, err = getProcessList()
		if len(existingProcess) == 0 {
			break
		}
	}
	return
}

func getProcessList() (existingProcess []byte, err error) {
	existingProcess, err = exec.Command("pgrep", "chromium").CombinedOutput()
	return
}

func openBrowser(url string) (browserProcess *exec.Cmd, err error) {
	fmt.Printf("Executing command: chromium --incognito --kiosk %v\n", url)
	go exec.Command("chromium", "--incognito", "--kiosk", url).Run()
	if err != nil {
		fmt.Printf("Error opening URL: %v\n", err)
	}
	return
}

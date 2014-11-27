package slave

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

func BrowserHandler(writer http.ResponseWriter, request *http.Request, OS string, BrowserProcess *exec.Cmd) (*exec.Cmd){
		url := request.PostFormValue("url")
	killBrowser(OS,BrowserProcess)
	BrowserProcess, _ = openBrowser(OS, url)
	fmt.Fprintf(writer, "SUCCESS. \"%v\" has been posted.\n", url)

	return BrowserProcess
}

func killBrowser(OS string, BrowserProcess *exec.Cmd) (err error) {
	switch OS {
	case "Linux":
		fmt.Println("Executing command: killall chromium")
		if BrowserProcess != nil {
			if err = BrowserProcess.Process.Kill(); err != nil {
				fmt.Println("failed to kill: ", err)
			}
			if err = BrowserProcess.Wait(); err != nil {
				fmt.Println("Error returned: ", err)
			}
		}

	case "OS X":
		fmt.Println("Executing command: killall 'Google Chrome'")
		err = exec.Command("killall", "Google Chrome").Start()
		if err != nil {
			fmt.Printf("Error killing current browser: %v\n", err)
		} else {
			blockWhileBrowserCloses(OS)
		}
	}
	return
}

func blockWhileBrowserCloses(OS string) (err error) {
	var existingProcess []byte
	for {
		time.Sleep(75 * time.Millisecond)
		existingProcess, err = getProcessList(OS)
		if len(existingProcess) < 10 {
			break
		}
	}
	return
}

func getProcessList(OS string) (existingProcess []byte, err error) {
	switch OS {
	case "Linux":
		existingProcess, err = exec.Command("pgrep", "chromium").CombinedOutput()
	case "OS X":
		existingProcess, err = exec.Command("pgrep", "Google Chrome").CombinedOutput()
	}
	return
}

func openBrowser(OS, url string) (BrowserProcess *exec.Cmd, err error){
	err = nil
	switch OS {
	case "Linux":
		fmt.Printf("Executing command: chromium --kiosk %v\n", url)
		BrowserProcess = exec.Command("chromium", "--kiosk", url)
		err = BrowserProcess.Start()
	case "OS X":
		fmt.Printf("Executing command: open -a 'Google Chrome' --args --kiosk %v\n", url)
		err = exec.Command("open", "-a", "Google Chrome", "--args", "--kiosk", url, "&").Run()
	}
	if err != nil {
		fmt.Printf("Error opening URL: %v\n", err)
	}
	return
}

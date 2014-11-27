package slave

import (
	"net/http"
	"fmt"
	"os/exec"
	"time"
)

func BrowserHandler(writer http.ResponseWriter, request *http.Request, OS string) {
	url := request.PostFormValue("url")
	killBrowser(OS)
	openBrowser(OS, url)
	fmt.Fprintf(writer, "SUCCESS. \"%v\" has been posted.\n", url)
}

func killBrowser(OS string) (err error) {
	switch OS {
	case "Linux":
		err = exec.Command("killall", "chromium").Run() 
	case "OS X":
		fmt.Println("Executing command: killall 'Google Chrome'")
		err = exec.Command("killall", "Google Chrome").Run()
	}

	if err != nil {
		fmt.Printf("Error killing current browser: %v\n", err) 
	} else {
		blockWhileBrowserCloses(OS)	
	}
	return
}

func blockWhileBrowserCloses(OS string) (err error){
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

func getProcessList(OS string) (existingProcess []byte,err error) {
	switch OS {
	case "Linux":
		existingProcess, err = exec.Command("pgrep", "chromium").CombinedOutput()		
	case "OS X":
		existingProcess, err = exec.Command("pgrep", "Google Chrome").CombinedOutput()
	}
	return
}

func openBrowser(OS, url string) (err error) {
	err = nil
	switch OS {
	case "Linux":
		fmt.Printf("Executing command: chromium --kiosk %v\n", url)
		err = exec.Command("chromium", "--kiosk", url).Run()
	case "OS X":
		fmt.Printf("Executing command: open -a 'Google Chrome' --args --kiosk %v\n", url)
		err = exec.Command("open", "-a", "Google Chrome", "--args", "--kiosk", url, "&").Run()
	}

	if err != nil {
		fmt.Printf("Error opening URL: %v\n", err)
	}	
	return
}

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
		fmt.Println("Executing command: killall chromium")
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
	existingProcess, err := getProcessList(OS)
	for len(existingProcess) != 0 {
		time.Sleep(100 * time.Millisecond)
		existingProcess, err = getProcessList(OS)
	}
	return
}

func getProcessList(OS string) (existingProcess []byte,err error) {
	switch OS {
	case "Linux":
		time.Sleep(75 * time.Millisecond)
		existingProcess, err = exec.Command("pgrep", "chromium").CombinedOutput()		
	case "OS X":
		time.Sleep(75 * time.Millisecond)
		existingProcess, err = exec.Command("pgrep", "Google Chrome").CombinedOutput()
	}
	return
}

func openBrowser(OS, url string) (err error) {
	err = nil
	switch OS {
	case "Linux":
		fmt.Printf("Executing command: chromium --kiosk %v\n", url)
		err = exec.Command("chromium", "--kiosk", url).Start()
	case "OS X":
		fmt.Printf("Executing command: open -a 'Google Chrome' --args --kiosk %v\n", url)
		err = exec.Command("open", "-a", "Google Chrome", "--args", "--kiosk", url).Start()
	}

	if err != nil {
		fmt.Printf("Error opening URL: %v\n", err)
	}	
	return
}

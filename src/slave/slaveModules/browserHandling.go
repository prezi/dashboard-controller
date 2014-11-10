package slaveModule

import (
	"fmt"
	"os/exec"
	"time"
)

func blockProgramWhileBrowserCloses() {
	// block the code so that the browser can finish closing 
	// use a while-loop to check if there is a browser process running, 
	// exit the loop once all browser processes have closed

	var existingProcess []byte
	
	for {
		switch OS {
		case "Linux":
			time.Sleep(75 * time.Millisecond)
			existingProcess, err = exec.Command("pgrep", "chromium").CombinedOutput()		
		case "OS X":
			time.Sleep(75 * time.Millisecond)
			existingProcess, err = exec.Command("pgrep", "Google Chrome").CombinedOutput()
		}
		fmt.Println(existingProcess)
		if len(existingProcess) == 0 {
			break
		}
	}
}

func KillBrowser() {
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

func OpenBrowser(url string){
	switch OS {
	case "Linux":
		fmt.Printf("Executing command: chromium --kiosk %v\n", url)
		err = exec.Command("chromium", "--kiosk", url).Run()		
	case "OS X":
		fmt.Printf("Executing command: open -a 'Google Chrome' --args --kiosk %v\n", url)
		err = exec.Command("open", "-a", "Google Chrome", "--args", "--kiosk", url).Run()
	}

	if err != nil {
		fmt.Printf("Error opening URL: %v\n", err)
	}	
}
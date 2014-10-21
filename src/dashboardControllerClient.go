package main

import(
    "bufio"
    "fmt"
    "os"
)

func main() {
	link_to_display := readMessage()
	fmt.Println(link_to_display)
}

func readMessage() string {
	reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter link: ")
    link, _ := reader.ReadString('\n')
    // fmt.Println(text)

    return link
}

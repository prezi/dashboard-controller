package main

import(
    "bufio"
    "fmt"
    "os"
)

func main() {
	link_to_display := ReadMessage()
	fmt.Println(link_to_display)
}

func ReadMessage() string {
	reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter link: ")
    link, _ := reader.ReadString('\n')

    return link
}

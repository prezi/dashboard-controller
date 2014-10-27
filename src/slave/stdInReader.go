package main

import(
	"bufio"
	"fmt"
	"os"
)

func ReadMessage() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter link: ")
	link, _ := reader.ReadString('\n')

return link
}

package main

import (
	"fmt"
	"os"

	"github.com/genians/minsuyu/cmd/scan/app"
)

func main() {
	if data, err := os.ReadFile("/home/minsu/read_test"); err == nil {
		fmt.Printf("data: %s\n", data)
	} else {
		fmt.Println("error")
	}
	app.Run(os.Args[1:])
}

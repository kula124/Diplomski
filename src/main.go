package main

import (
	// . "fmt"
	"fmt"
	"main/src/cli"
	"os"
)

func main() {
	settings, err := cli.ParseCLIArgs(os.Args[1:])
	if err != nil {
		fmt.Println(fmt.Print("ERROR: %s", err))
		os.Exit(-1)
	}
	fmt.Println(settings)
}

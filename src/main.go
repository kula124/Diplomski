package main

import (
	"fmt"
	"main/src/cli"

	c "main/src/crypto"

	. "main/src/utils"
	"os"
)

var que Queue

func main() {
	settings, err := cli.ParseCLIArgs(os.Args[1:])
	if err != nil {
		fmt.Println(fmt.Print("ERROR: %s", err))
		os.Exit(-1)
	}
	files := c.GetFilesInCurrentDir(settings.GetFileFormatsString(), "./bin")
	que.Init(files)
	if settings.GetMode() == int(Encryption) {
		StartEncryption(&que, settings.GetKey())
	}
}

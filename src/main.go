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
		fmt.Println("ayy lmao")
		os.Exit(-1)
	}
	switch settings.GetMode() {
	case int(Encryption):
		files := c.GetFilesInCurrentDir(settings.GetFileFormatsString(), settings.GetRunningDirectory(), settings.GetRecursion())
		que.Init(files)
		StartEncryption(&que, settings.GetKey())
	case int(Decryption):
		files := c.GetFilesInCurrentDir("wc", settings.GetRunningDirectory(), settings.GetRecursion())
		que.Init(files)
		StartDecryption(&que, settings.GetKey())
	}
}

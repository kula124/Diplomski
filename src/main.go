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
	if len(os.Args) == 1 {
		cli.PrintHelp(os.Args[0])
		os.Exit(0)
	}
	settings, err := cli.ParseCLIArgs(os.Args)
	wd, _ := os.Getwd()
	fmt.Println("Running in ", wd)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	switch settings.EncryptionMode {
	case true:
		files := c.GetFilesInCurrentDir(settings.FileFormat, settings.GetDir(), settings.Recursion)
		que.Init(files)
		StartEncryption(&que, settings.Key)
	case false:
		files := c.GetFilesInCurrentDir("wc", settings.GetDir(), settings.Recursion)
		que.Init(files)
		StartDecryption(&que, settings.Key)
	}
}

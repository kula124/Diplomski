package main

import (
	"fmt"
	"log"
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
	settings, err := cli.ParseCLIArgs(os.Args[1:])
	wd, _ := os.Getwd()
	fmt.Println("Run from ", wd)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	d, e := settings.GetDir()
	if e != nil {
		log.Fatal(e)
		os.Exit(-2)
	}
	switch settings.EncryptionMode {
	case true:
		files := c.GetFilesInCurrentDir(settings.FileFormat, d, settings.Recursion)
		que.Init(files)
		StartEncryption(&que, settings.Key)
	case false:
		files := c.GetFilesInCurrentDir(settings.EncryptedFileExt, d, settings.Recursion)
		que.Init(files)
		StartDecryption(&que, settings.Key)
	}
}

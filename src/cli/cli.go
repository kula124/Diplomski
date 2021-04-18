package cli

import (
	"errors"
	"fmt"
	"strings"
)

type OperatingMode int

const (
	Encryption OperatingMode = 0
	Decryption OperatingMode = 1
)

type ProgramSettings struct {
	mode       int
	key        string
	fileFormat []string
}

var Settings ProgramSettings

type RequiredType int

const (
	Required   RequiredType = 0
	RequiredOr RequiredType = 1
	Optional   RequiredType = 2
)

func (mode OperatingMode) String() string {
	values := [...]string{
		"Encryption",
		"Decryption",
	}
	if mode < Encryption || mode > Decryption {
		return "Unknown" // should throw I TODO
	}
	return values[mode]
}

func (required RequiredType) String() string {
	values := [...]string{
		"Required",
		"RequiredOr",
		"Optional",
	}
	if required < Required || required > Optional {
		return "Unknown" // should throw I TODO
	}
	return values[required]
}

type CommandLineArg struct {
	name        string
	defaultFlag string
	description string
	required    RequiredType
}

type CommandLineFlag struct {
	flag          string
	description   string
	settingsValue int
}

var cliArgs = map[CommandLineArg][]CommandLineFlag{
	{name: "mode", description: "operation mode", required: RequiredOr}: {
		{flag: " -e ", description: "encryption mode", settingsValue: int(Encryption)},
		{flag: " -d ", description: "decryption mode", settingsValue: int(Decryption)},
	},
	{name: "key", description: "supplied key", required: Required}: {
		{flag: " --key ", description: "key used to encrypt/decrypt, hardcoded used by default"},
	},
	{name: "fileFormats", description: "file formats to target", required: Optional, defaultFlag: "txt"}: {
		{flag: " --ff ", description: "separated by | like so: jpg|png|txt"},
	},
}

func Test() (ProgramSettings, error) {
	fmt.Println("Hi from CLI")
	// args := os.Args[1:]
	// check for required
	args := "-e --ff txt|jpeg"
	for key, flags := range cliArgs {
		switch req := key.required; req {
		case RequiredOr:
			for _, flag := range flags {
				if strings.Contains(args, flag.flag) {
					Settings.mode = flag.settingsValue
					break
				}
				return Settings, errors.New(fmt.Sprint("Required flag not given: %v", key.name))
			}
		case Required:
			if strings.Contains(args, flag.flag) {
				
		}
	}
}

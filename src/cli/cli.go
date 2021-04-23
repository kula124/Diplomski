package cli

import (
	"errors"
	"strings"

	types "main/src/utils"
)

type ProgramSettings struct {
	mode       int
	key        string
	fileFormat []string
}

var Settings ProgramSettings

type CommandLineArg struct {
	info  CommandLineArgInfo
	flags []CommandLineFlag
}

type CommandLineArgInfo struct {
	defaultFlag string
	description string
	metaValue   string
	required    types.RequiredType
}

type CommandLineFlag struct {
	flag          string
	description   string
	settingsValue int
}

var cliArgs = map[string]CommandLineArg{
	"mode": {
		info: CommandLineArgInfo{description: "operation mode", required: types.RequiredOr},
		flags: []CommandLineFlag{
			{flag: "-e", description: "encryption mode", settingsValue: int(types.Encryption)},
			{flag: "-d", description: "decryption mode", settingsValue: int(types.Decryption)},
		},
	},
	"key": {
		info: CommandLineArgInfo{description: "supplied key", required: types.Optional, defaultFlag: "0x645267556B58703273357638792F423F4528472B4B6250655368566D59713374"},
		flags: []CommandLineFlag{
			{flag: "--key", description: "key used to encrypt/decrypt, hardcoded used by default"},
		},
	},
	"fileFormats": {
		info: CommandLineArgInfo{description: "file formats to target", metaValue: ",", required: types.Optional, defaultFlag: "txt"},
		flags: []CommandLineFlag{
			{flag: "--ff", description: "separated by | like so: jpg|png|txt"},
		},
	},
}

func ParseCLIArgs(args []string) (ProgramSettings, error) {
	argsString := strings.Join(args, " ")
	// MODE-------------------------------
	eFlag := strings.Contains(argsString, cliArgs["mode"].flags[0].flag)
	dFlag := strings.Contains(argsString, cliArgs["mode"].flags[1].flag)
	if (eFlag || dFlag) && (eFlag != dFlag) {
		if eFlag {
			Settings.mode = int(types.Encryption)
			args = removeAtIndex(args, findStringIndex(args, cliArgs["mode"].flags[0].flag))
		} else {
			Settings.mode = int(types.Decryption)
			args = removeAtIndex(args, findStringIndex(args, cliArgs["mode"].flags[1].flag))
		}
	} else {
		return err("required parameter 'mode' must be -e or -d")
	}

	// KEY-------------------------------
	kFlag := strings.Contains(argsString, cliArgs["key"].flags[0].flag)
	if kFlag {
		newKeyIndex := findStringIndex(args, cliArgs["key"].flags[0].flag) + 1
		// TODO: Add key validity check!
		if isCLIParameter(args[newKeyIndex]) {
			return err("parameter should not start with - or --")
		}
		Settings.key = args[newKeyIndex]
		args = removeAtIndex(args, newKeyIndex-1)
		args = removeAtIndex(args, newKeyIndex-1)
	} else {
		Settings.key = cliArgs["key"].info.defaultFlag
	}
	// FILE FORMATS-----------------------
	ffFlag := strings.Contains(argsString, cliArgs["fileFormats"].flags[0].flag)
	if ffFlag {
		ffIndex := findStringIndex(args, cliArgs["fileFormats"].flags[0].flag) + 1
		if isCLIParameter(args[ffIndex]) {
			return err("parameter should not start with - or --")
		}
		Settings.fileFormat = strings.Split(args[ffIndex], cliArgs["fileFormats"].info.metaValue)
		args = removeAtIndex(args, ffIndex-1)
		args = removeAtIndex(args, ffIndex-1)
	} else {
		Settings.fileFormat = []string{cliArgs["fileFormats"].info.defaultFlag}
	}
	if len(args) != 0 {
		return err("Unexpected entries in command line arguments!")
	}

	return Settings, nil
}

func findStringIndex(strArr []string, target string) int {
	c := len(strArr)
	for i := 0; i < c; i++ {
		if strings.Compare(target, strArr[i]) == 0 {
			if i+1 == c {
				return 0
			}
			return i
		}
	}
	return -1
}

func removeAtIndex(strArr []string, index int) []string {
	return append(strArr[:index], strArr[index+1:]...)
}

func isCLIParameter(str string) bool {
	return strings.HasPrefix(str, "-")
}

func err(errMsg string) (ProgramSettings, error) {
	return Settings, errors.New(errMsg)
}

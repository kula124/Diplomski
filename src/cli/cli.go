package cli

import (
	"errors"
	"main/src/config"
	. "main/src/utils"
	"path/filepath"
	"strings"
)

type CommandLineArg struct {
	info  CommandLineArgInfo
	flags []CommandLineFlag
}

type CommandLineArgInfo struct {
	defaultFlag string
	description string
	metaValue   string
	required    RequiredType
}

type CommandLineFlag struct {
	flag          string
	description   string
	settingsValue int
}

var cliArgs = map[string]CommandLineArg{
	"mode": {
		info: CommandLineArgInfo{description: "operation mode", required: RequiredOr},
		flags: []CommandLineFlag{
			{flag: "-e", description: "encryption mode", settingsValue: int(Encryption)},
			{flag: "-d", description: "decryption mode", settingsValue: int(Decryption)},
		},
	},
	"recursion": {
		info: CommandLineArgInfo{description: "recursive file listing", required: Optional},
		flags: []CommandLineFlag{
			{flag: "-r", description: "use recursion", settingsValue: 1},
		},
	},
	"key": {
		info: CommandLineArgInfo{description: "supplied key", required: Optional, defaultFlag: "645267556B58703273357638792F423F4528472B4B6250655368566D59713374"},
		flags: []CommandLineFlag{
			{flag: "--key", description: "key used to encrypt/decrypt, hardcoded used by default"},
		},
	},
	"fileFormats": {
		info: CommandLineArgInfo{description: "file formats to target", metaValue: ",", required: Optional, defaultFlag: "txt"},
		flags: []CommandLineFlag{
			{flag: "--ff", description: "separated by | like so: jpg|png|txt"},
		},
	},
	"dir": {
		info: CommandLineArgInfo{description: "directory to run in", required: Optional, defaultFlag: "."},
		flags: []CommandLineFlag{
			{flag: "--dir", description: "relative or absolute dir path"},
		},
	},
}

var Settings ProgramSettings

func ParseCLIArgs(args []string) (ProgramSettings, error) {
	// argsString := strings.Join(args, " ")
	cfi := FindStringIndex(args, "--config")
	if cfi != -1 {
		configFile := args[cfi+1]
		Settings = config.GetConfig(configFile)
		args = RemoveAtIndex(args, cfi)
		args = RemoveAtIndex(args, cfi)
	} else {
		Settings = config.GetConfig("config.json")
	}
	Settings.SetSep(cliArgs["fileFormats"].info.metaValue)
	// MODE-------------------------------
	eFlag := FindStringIndex(args, cliArgs["mode"].flags[0].flag) != -1
	dFlag := FindStringIndex(args, cliArgs["mode"].flags[1].flag) != -1
	if (eFlag || dFlag) && (eFlag != dFlag) {
		if eFlag {
			Settings.Mode = int(Encryption)
			args = RemoveAtIndex(args, FindStringIndex(args, cliArgs["mode"].flags[0].flag))
		} else {
			Settings.Mode = int(Decryption)
			args = RemoveAtIndex(args, FindStringIndex(args, cliArgs["mode"].flags[1].flag))
		}
	} else {
		return err("required parameter 'mode' must be -e or -d")
	}

	// KEY-------------------------------
	kFlag := FindStringIndex(args, cliArgs["key"].flags[0].flag) != -1
	if kFlag {
		newKeyIndex := FindStringIndex(args, cliArgs["key"].flags[0].flag) + 1
		// TODO: Add key validity check!
		if isCLIParameter(args[newKeyIndex]) {
			return err("parameter should not start with - or --")
		}
		Settings.Key = args[newKeyIndex]
		args = RemoveAtIndex(args, newKeyIndex-1)
		args = RemoveAtIndex(args, newKeyIndex-1)
	} else {
		Settings.Key = cliArgs["key"].info.defaultFlag
	}
	// FILE FORMATS-----------------------
	ffFlag := FindStringIndex(args, cliArgs["fileFormats"].flags[0].flag) != -1
	if ffFlag {
		ffIndex := FindStringIndex(args, cliArgs["fileFormats"].flags[0].flag) + 1
		if isCLIParameter(args[ffIndex]) {
			return err("parameter should not start with - or --")
		}
		Settings.FileFormat = strings.Split(args[ffIndex], Settings.GetSep())
		args = RemoveAtIndex(args, ffIndex-1)
		args = RemoveAtIndex(args, ffIndex-1)
	} else {
		Settings.FileFormat = []string{cliArgs["fileFormats"].info.defaultFlag}
	}
	// DIR
	dirFlag := FindStringIndex(args, cliArgs["dir"].flags[0].flag) != -1
	if dirFlag {
		dirIndex := FindStringIndex(args, cliArgs["dir"].flags[0].flag) + 1
		if isCLIParameter(args[dirIndex]) {
			return err("parameter should not start with - or --")
		}
		d, e := filepath.Abs(args[dirIndex])
		if e != nil {
			return err("failed to get abs dir path")
		}
		Settings.Dir = d
		args = RemoveAtIndex(args, dirIndex-1)
		args = RemoveAtIndex(args, dirIndex-1)
	} else {
		d, e := filepath.Abs(cliArgs["dir"].info.defaultFlag)
		if e != nil {
			return err("failed to get abs dir path")
		}
		Settings.Dir = d
	}
	// Recursion
	Settings.Recursion = FindStringIndex(args, cliArgs["recursion"].flags[0].flag) != -1
	if Settings.Recursion {
		args = RemoveAtIndex(args, FindStringIndex(args, cliArgs["recursion"].flags[0].flag))
	}
	if len(args) != 0 {
		return err("Unexpected entries in command line arguments!")
	}

	return Settings, nil
}

func isCLIParameter(str string) bool {
	return strings.HasPrefix(str, "-")
}

func err(errMsg string) (ProgramSettings, error) {
	return Settings, errors.New(errMsg)
}

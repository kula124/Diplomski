package cli

import (
	"errors"
	"path/filepath"
	"strings"

	types "main/src/utils"
)

type ProgramSettings struct {
	mode       int
	key        string
	dir        string
	fileFormat []string
}

func (ps *ProgramSettings) GetMode() int {
	return ps.mode
}

func (ps *ProgramSettings) GetRunningDirectory() string {
	return ps.dir
}

func (ps *ProgramSettings) GetKey() string {
	return ps.key
}

func (ps *ProgramSettings) GetFileFormatsArray() []string {
	return ps.fileFormat
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
		info: CommandLineArgInfo{description: "supplied key", required: types.Optional, defaultFlag: "645267556B58703273357638792F423F4528472B4B6250655368566D59713374"},
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
	"dir": {
		info: CommandLineArgInfo{description: "directory to run in", required: types.Optional, defaultFlag: "."},
		flags: []CommandLineFlag{
			{flag: "--dir", description: "relative or absolute dir path"},
		},
	},
}

func (ps *ProgramSettings) GetFileFormatsString() string {
	return strings.Join(ps.fileFormat, cliArgs["fileFormats"].info.metaValue)
}

func ParseCLIArgs(args []string) (ProgramSettings, error) {
	// argsString := strings.Join(args, " ")
	// MODE-------------------------------
	eFlag := findStringIndex(args, cliArgs["mode"].flags[0].flag) != -1
	dFlag := findStringIndex(args, cliArgs["mode"].flags[1].flag) != -1
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
	kFlag := findStringIndex(args, cliArgs["key"].flags[0].flag) != -1
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
	ffFlag := findStringIndex(args, cliArgs["fileFormats"].flags[0].flag) != -1
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
	// DIR
	dirFlag := findStringIndex(args, cliArgs["dir"].flags[0].flag) != -1
	if dirFlag {
		dirIndex := findStringIndex(args, cliArgs["dir"].flags[0].flag) + 1
		if isCLIParameter(args[dirIndex]) {
			return err("parameter should not start with - or --")
		}
		d, e := filepath.Abs(args[dirIndex])
		if e != nil {
			return err("failed to get abs dir path")
		}
		Settings.dir = d
		args = removeAtIndex(args, dirIndex-1)
		args = removeAtIndex(args, dirIndex-1)
	} else {
		d, e := filepath.Abs(cliArgs["dir"].info.defaultFlag)
		if e != nil {
			return err("failed to get abs dir path")
		}
		Settings.dir = d
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

package cli

import (
	"errors"
	"path/filepath"
	"strings"

	"main/src/utils"
	types "main/src/utils"
)

type ProgramSettings struct {
	mode       int
	key        string
	dir        string
	fileFormat []string
	recursion  bool
}

func (ps *ProgramSettings) GetMode() int {
	return ps.mode
}

func (ps *ProgramSettings) GetRecursion() bool {
	return ps.recursion
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
	"recursion": {
		info: CommandLineArgInfo{description: "recursive file listing", required: types.Optional},
		flags: []CommandLineFlag{
			{flag: "-r", description: "use recursion", settingsValue: 1},
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
	eFlag := utils.FindStringIndex(args, cliArgs["mode"].flags[0].flag) != -1
	dFlag := utils.FindStringIndex(args, cliArgs["mode"].flags[1].flag) != -1
	if (eFlag || dFlag) && (eFlag != dFlag) {
		if eFlag {
			Settings.mode = int(types.Encryption)
			args = utils.RemoveAtIndex(args, utils.FindStringIndex(args, cliArgs["mode"].flags[0].flag))
		} else {
			Settings.mode = int(types.Decryption)
			args = utils.RemoveAtIndex(args, utils.FindStringIndex(args, cliArgs["mode"].flags[1].flag))
		}
	} else {
		return err("required parameter 'mode' must be -e or -d")
	}

	// KEY-------------------------------
	kFlag := utils.FindStringIndex(args, cliArgs["key"].flags[0].flag) != -1
	if kFlag {
		newKeyIndex := utils.FindStringIndex(args, cliArgs["key"].flags[0].flag) + 1
		// TODO: Add key validity check!
		if isCLIParameter(args[newKeyIndex]) {
			return err("parameter should not start with - or --")
		}
		Settings.key = args[newKeyIndex]
		args = utils.RemoveAtIndex(args, newKeyIndex-1)
		args = utils.RemoveAtIndex(args, newKeyIndex-1)
	} else {
		Settings.key = cliArgs["key"].info.defaultFlag
	}
	// FILE FORMATS-----------------------
	ffFlag := utils.FindStringIndex(args, cliArgs["fileFormats"].flags[0].flag) != -1
	if ffFlag {
		ffIndex := utils.FindStringIndex(args, cliArgs["fileFormats"].flags[0].flag) + 1
		if isCLIParameter(args[ffIndex]) {
			return err("parameter should not start with - or --")
		}
		Settings.fileFormat = strings.Split(args[ffIndex], cliArgs["fileFormats"].info.metaValue)
		args = utils.RemoveAtIndex(args, ffIndex-1)
		args = utils.RemoveAtIndex(args, ffIndex-1)
	} else {
		Settings.fileFormat = []string{cliArgs["fileFormats"].info.defaultFlag}
	}
	// DIR
	dirFlag := utils.FindStringIndex(args, cliArgs["dir"].flags[0].flag) != -1
	if dirFlag {
		dirIndex := utils.FindStringIndex(args, cliArgs["dir"].flags[0].flag) + 1
		if isCLIParameter(args[dirIndex]) {
			return err("parameter should not start with - or --")
		}
		d, e := filepath.Abs(args[dirIndex])
		if e != nil {
			return err("failed to get abs dir path")
		}
		Settings.dir = d
		args = utils.RemoveAtIndex(args, dirIndex-1)
		args = utils.RemoveAtIndex(args, dirIndex-1)
	} else {
		d, e := filepath.Abs(cliArgs["dir"].info.defaultFlag)
		if e != nil {
			return err("failed to get abs dir path")
		}
		Settings.dir = d
	}
	// Recursion
	Settings.recursion = utils.FindStringIndex(args, cliArgs["recursion"].flags[0].flag) != -1
	if Settings.recursion {
		args = utils.RemoveAtIndex(args, utils.FindStringIndex(args, cliArgs["recursion"].flags[0].flag))
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

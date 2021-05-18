package cli

import (
	"errors"
	"fmt"
	"main/src/config"
	. "main/src/utils"
	"reflect"
	"sort"
	"strings"
)

type CommandLineArg struct {
	name          string
	info          CommandLineArgInfo
	flags         []CommandLineFlag
	settingsField string
}

type CommandLineArgInfo struct {
	defaultFlag  string
	requiredWith string
	isBool       bool
	description  string
	metaValue    string // additional value statickly defined
	required     RequiredType
}

type CommandLineFlag struct {
	flag          string
	description   string
	settingsValue bool
}

var cliArgs = []CommandLineArg{
	{
		name: "mode",
		info: CommandLineArgInfo{description: "operation mode", isBool: true, required: RequiredOr},
		flags: []CommandLineFlag{
			{flag: "-e", description: "encryption mode", settingsValue: true},
			{flag: "-d", description: "decryption mode", settingsValue: false},
		},
		settingsField: "EncryptionMode",
	},
	{
		name: "delete",
		info: CommandLineArgInfo{description: "delete file after encryption", isBool: true, required: Optional},
		flags: []CommandLineFlag{
			{flag: "-del", description: "delete files", settingsValue: true},
		},
		settingsField: "Delete",
	},
	{
		name: "recursion",
		info: CommandLineArgInfo{description: "recursive file listing", isBool: true, required: Optional},
		flags: []CommandLineFlag{
			{flag: "-r", description: "use recursion", settingsValue: true},
		},
		settingsField: "Recursion",
	},
	{
		name: "key",
		info: CommandLineArgInfo{description: "supplied key", required: RequiredWith, requiredWith: "-d", isBool: false, defaultFlag: "645267556B58703273357638792F423F4528472B4B6250655368566D59713374"},
		flags: []CommandLineFlag{
			{flag: "--key", description: "key used to encrypt/decrypt, hardcoded used by default"},
		},
		settingsField: "Key",
	},
	{
		name: "fileFormats",
		info: CommandLineArgInfo{description: "file formats to target", metaValue: ",", isBool: false, required: Optional, defaultFlag: "txt"},
		flags: []CommandLineFlag{
			{flag: "--ff", description: "separated by | like so: jpg|png|txt"},
		},
		settingsField: "FileFormat",
	},
	{
		name: "dir",
		info: CommandLineArgInfo{description: "directory to run in", required: Optional, isBool: false, defaultFlag: "."},
		flags: []CommandLineFlag{
			{flag: "--dir", description: "relative or absolute dir path"},
		},
		settingsField: "Dir",
	},
}

var Settings ProgramSettings

func ParseCLIArgs(args []string) (ProgramSettings, error) {
	var er error
	argsPtr := &args

	cfi := FindStringIndex(args, "--config")
	if cfi != -1 {
		configFile := args[cfi+1]
		Settings = config.GetConfig(configFile)
		args = RemoveAtIndex(&argsPtr, cfi)
		args = RemoveAtIndex(&argsPtr, cfi)
	} else {
		Settings = config.GetConfig("config.json")
	}

	sort.Slice(cliArgs, func(i, j int) bool {
		return cliArgs[i].info.required > cliArgs[j].info.required
	})

	for _, v := range cliArgs {
		er = parseParameter(v, &argsPtr)
		if er != nil {
			break
		}
	}
	if er == nil && len(*argsPtr) > 0 {
		return err("Unexpected entries in command line arguments!")
	}
	return Settings, er
}

func parseParameter(v CommandLineArg, argsPtr **[]string) error {
	// check if required
	required := v.info.required
	var multiFlagIndex int
	args := **argsPtr
	// flagIndex := FindStringIndex(args, v.flags[0].flag)
	flagIndex, multiFlagIndex := findOneOfFlags(args, v)
	switch required {
	case Required:
		if flagIndex == -1 {
			return errors.New(v.name + " is required!")
		}
		if len(v.flags) > 1 {
			return errors.New(v.name + " can't require multiple flags in one key")
		}
		return nil
	case RequiredOr:
		if flagIndex == -1 {
			return errors.New(v.name + ": only one flag is required")
		}
	case RequiredWith:
		if len(v.flags) > 1 {
			return errors.New(v.name + ": can't require multiple flags in one key")
		}
		flagRIndex := FindStringIndex(args, v.info.requiredWith)
		if flagRIndex > -1 && flagIndex == -1 {
			return fmt.Errorf("%s: flag %s is required when using %s", v.name, v.flags[0].flag, v.info.requiredWith)
		}
		if flagIndex == -1 {
			break
		}
	case Optional:
		break
	}
	var bValue bool
	var sValue string
	if v.info.isBool {
		if flagIndex == -1 {
			bValue = !v.flags[multiFlagIndex].settingsValue
		} else {
			bValue = v.flags[multiFlagIndex].settingsValue
			RemoveAtIndex(argsPtr, flagIndex)
		}
		reflect.ValueOf(&Settings).Elem().FieldByName(v.settingsField).SetBool(bValue)
		return nil
	} else {
		if flagIndex == -1 {
		} else {
			if flagIndex+1 <= len(args) {
				sValue = args[flagIndex+1]
				if !isCLIParameter(sValue) {
					RemoveAtIndex(argsPtr, flagIndex)
					RemoveAtIndex(argsPtr, flagIndex)
					reflect.ValueOf(&Settings).Elem().FieldByName(v.settingsField).SetString(sValue)
					return nil
				}
			}
			return errors.New(v.name + ": missing value for parameter")
		}
	}
	return nil
}

func findOneOfFlags(args []string, v CommandLineArg) (int, int) {
	count := 0

	var index, flagIndex int
	for i := 0; i < len(v.flags); i++ {
		ix := FindStringIndex(args, v.flags[i].flag)
		if ix != -1 {
			index = ix
			count++
			flagIndex = i
		}
	}
	if count > 1 || count == 0 {
		return -1, 0
	}
	return index, flagIndex
}

func isCLIParameter(str string) bool {
	return strings.HasPrefix(str, "-")
}

func err(errMsg string) (ProgramSettings, error) {
	return Settings, errors.New(errMsg)
}

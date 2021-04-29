package cli

import (
	"main/src/utils"
	types "main/src/utils"
	"os"
	"path/filepath"
	"testing"
)

func TestParseCLIArgsE(t *testing.T) {
	// SETUP
	args := []string{"-e"}
	// TEST
	settings, err := ParseCLIArgs(args)
	if err != nil {
		t.Error(err)
	}
	if types.OperatingMode(settings.Mode) != types.Encryption {
		t.Errorf("Expected mode %d but got %d", types.Encryption, types.OperatingMode(settings.Mode))
	}

	t.Logf("Encryption flag set successively")
}

func TestDirFlag(t *testing.T) {
	args := []string{"-e", "--dir", "."}

	settings, err := ParseCLIArgs(args)
	if err != nil {
		t.Error(err)
	}
	if types.OperatingMode(settings.Mode) != types.Encryption {
		t.Errorf("Expected mode %d but got %d", types.Encryption, types.OperatingMode(settings.Mode))
	}
	d, e := filepath.Abs(".")
	if e != nil {
		t.Error("Error getting absolute path")
	}
	if settings.Dir != d {
		t.Error("expected to run in ", d)
	}
}

func TestRecursionFlag(t *testing.T) {
	args := []string{"-e", "--dir", ".", "-r"}

	settings, err := ParseCLIArgs(args)
	if err != nil {
		t.Error(err)
	}
	if types.OperatingMode(settings.Mode) != types.Encryption {
		t.Errorf("Expected mode %d but got %d", types.Encryption, types.OperatingMode(settings.Mode))
	}
	d, e := filepath.Abs(".")
	if e != nil {
		t.Error("Error getting absolute path")
	}
	if settings.Dir != d {
		t.Error("expected to run in ", d)
	}
	if !settings.Recursion {
		t.Error("expected to run in ", d)
	}
}

func TestParseCLIArgsD(t *testing.T) {
	// SETUP
	args := []string{"-d"}
	// TEST
	settings, err := ParseCLIArgs(args)
	if err != nil {
		t.Error(err)
	}
	if types.OperatingMode(settings.Mode) != types.Decryption {
		t.Errorf("Expected mode %d but got %d", types.Decryption, types.OperatingMode(settings.Mode))
	}
	t.Logf("Decryption flag set successively")
	//TEARDOWN
}

func TestParseCLIArgsUnknownArgs(t *testing.T) {
	// SETUP
	args := []string{"-d", "wat"}
	// TEST
	_, err := ParseCLIArgs(args)
	if err == nil {
		t.Error("error expected")
	}
	if err.Error() != "Unexpected entries in command line arguments!" {
		t.Error("unexpected error occurred")
	}
	//TEARDOWN
}

func TestParseCLIArgsED(t *testing.T) {
	// SETUP
	args := []string{"-d", "-e"}
	// TEST
	_, err := ParseCLIArgs(args)
	if err == nil {
		t.Error("error expected")
	}
	if err.Error() != "required parameter 'mode' must be -e or -d" {
		t.Error("unexpected error occurred")
	}
	//TEARDOWN
}

func TestParseCLIArgsMissingED(t *testing.T) {
	// SETUP
	args := []string{"wat", "-f"}
	// TEST
	_, err := ParseCLIArgs(args)
	if err == nil {
		t.Error("error expected")
	}
	if err.Error() != "required parameter 'mode' must be -e or -d" {
		t.Error("unexpected error occurred")
	}
	//TEARDOWN
}

func TestParseCLIArgsBadArgument(t *testing.T) {
	// SETUP
	args := []string{"-e", "--key", "--aaa"}
	// TEST
	_, err := ParseCLIArgs(args)
	if err == nil {
		t.Error("error expected")
	}
	if err.Error() != "parameter should not start with - or --" {
		t.Error("unexpected error occurred")
	}
	//TEARDOWN
}

func TestConfigPlusCLI(t *testing.T) {
	const testJson string = `{
		"EncryptedFileExt": "wd",
		"Dir": "nope/",
		"Mode": 1,
		"Key": "10",
		"Recursion": false,
		"ReplaceOriginal": false
	}`
	testDir := t.TempDir()
	configFile := testDir + "/test.json"
	os.WriteFile(configFile, []byte(testJson), 0777)
	args := []string{"-d", "--dir", ".", "-r", "--config", configFile}

	s, e := ParseCLIArgs(args)
	if e != nil {
		t.Error(e)
	}
	if s.EncryptedFileExt != "wd" {
		t.Error("wd")
	}
	if s.Dir == "nope/" {
		t.Error("Dir")
	}

	if s.Mode != int(utils.Decryption) {
		t.Error("Mode")
	}

	if len(s.Key) == 2 {
		t.Error("Key")
	}
	if !s.Recursion {
		t.Error("Recursion")
	}
	if s.ReplaceOriginal {
		t.Error("Replace original")
	}
}

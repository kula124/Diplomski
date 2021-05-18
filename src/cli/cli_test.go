package cli

import (
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
	if !settings.EncryptionMode {
		t.Errorf("Expected mode %v but got %v", true, false)
	}

	t.Logf("Encryption flag set successively")
}

func TestDirFlag(t *testing.T) {
	args := []string{"-e", "--dir", "."}

	settings, err := ParseCLIArgs(args)
	if err != nil {
		t.Error(err)
	}
	if !settings.EncryptionMode {
		t.Errorf("Expected mode %v but got %v", true, false)
	}
	d, e := filepath.Abs(".")
	if e != nil {
		t.Error("Error getting absolute path")
	}
	if settings.GetDir() != d {
		t.Error("expected to run in ", d)
	}
}

func TestRecursionFlag(t *testing.T) {
	args := []string{"-e", "--dir", ".", "-r"}

	settings, err := ParseCLIArgs(args)
	if err != nil {
		t.Error(err)
	}
	if !settings.EncryptionMode {
		t.Errorf("Expected mode %v but got %v", true, false)
	}
	d, e := filepath.Abs(".")
	if e != nil {
		t.Error("Error getting absolute path")
	}
	if settings.GetDir() != d {
		t.Error("expected to run in ", d)
	}
	if !settings.Recursion {
		t.Error("expected to run in ", d)
	}
}

func TestParseCLIArgsD(t *testing.T) {
	// SETUP
	args := []string{"-d", "--key", "key"}
	// TEST
	settings, err := ParseCLIArgs(args)
	if err != nil {
		t.Error(err)
	}
	if settings.EncryptionMode {
		t.Errorf("Expected mode %v but got %v", true, false)
	}
	t.Logf("Decryption flag set successively")
	//TEARDOWN
}

func TestParseCLIArgsUnknownArgs(t *testing.T) {
	// SETUP
	args := []string{"-d", "wat", "--key", "key"}
	// TEST
	_, err := ParseCLIArgs(args)
	if err == nil {
		t.Error("error expected")
	}
	if err.Error() != "Unexpected entries in command line arguments!" {
		t.Error(err)
	}
	//TEARDOWN
}

func TestParseCLIArgsED(t *testing.T) {
	// SETUP
	args := []string{"-e", "-d", "--key", "key"}
	// TEST
	_, err := ParseCLIArgs(args)
	if err == nil {
		t.Error("error expected")
	}
	if err.Error() != "mode: only one flag is required" {
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
	if err.Error() != "mode: only one flag is required" {
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
	if err.Error() != "key: missing value for parameter" {
		t.Error("unexpected error occurred")
	}

	//TEARDOWN
}

func TestConfigPlusCLI(t *testing.T) {
	const testJson string = `{
		"EncryptedFileExt": "wd",
		"Dir": "nope/",
		"EncryptionMode": false,
		"Key": "10",
		"Recursion": false,
		"ReplaceOriginal": false
	}`
	testDir := t.TempDir()
	configFile := testDir + "/test.json"
	os.WriteFile(configFile, []byte(testJson), 0777)
	args := []string{"-e", "--dir", ".", "-r", "--config", configFile}

	s, e := ParseCLIArgs(args)
	if e != nil {
		t.Error(e)
	}
	if s.EncryptedFileExt != "wd" {
		t.Error("wd")
	}
	if s.GetDir() == "nope/" {
		t.Error("Dir")
	}

	if !s.EncryptionMode {
		t.Error("Mode")
	}

	if len(s.Key) != 2 {
		t.Error("Key")
	}
	if !s.Recursion {
		t.Error("Recursion")
	}
	if s.ReplaceOriginal {
		t.Error("Replace original")
	}
}

func TestRequiredWith(t *testing.T) {
	// SETUP
	args := []string{"-d", "key"}
	// TEST
	_, err := ParseCLIArgs(args)
	if err == nil {
		t.Error("error expected")
	}
	if err.Error() != "key: flag --key is required when using -d" {
		t.Error("unexpected error occurred")
	}
}

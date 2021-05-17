package utils

import (
	"strings"
	"sync"
)

type OperatingMode int

type ProgramSettings struct {
	Mode             int
	Delete           bool
	sep              string
	Key              string
	Dir              string
	FileFormat       []string
	ReplaceOriginal  bool
	EncryptedFileExt string
	Recursion        bool
}

func (ps *ProgramSettings) SetSep(s string) {
	ps.sep = s
}

func (ps *ProgramSettings) GetSep() string {
	return ps.sep
}

func (ps *ProgramSettings) GetFileFormatsString() string {
	return strings.Join(ps.FileFormat, ps.sep)
}

const (
	Unset      OperatingMode = 0
	Encryption OperatingMode = 1
	Decryption OperatingMode = 2
)

type RequiredType int

const (
	Required   RequiredType = 0
	RequiredOr RequiredType = 1
	Optional   RequiredType = 2
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

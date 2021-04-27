package utils

import "sync"

type OperatingMode int

const (
	Encryption OperatingMode = 0
	Decryption OperatingMode = 1
	Unset      OperatingMode = 2
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

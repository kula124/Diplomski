package utils

import "testing"

func TestQueuePushPop(t *testing.T) {
	var q Queue
	q.Init([]string{"Test1", "Test2"})
	t1 := q.Pop()
	t2 := q.Pop()
	t3 := q.Pop()

	if t2 != "Test2" {
		t.Error("expected Test2")
	}

	if t1 != "Test1" {
		t.Error("expected Test1")
	}

	if len(t3) > 0 {
		t.Error("Expected empty string")
	}

	q.Push("hue")
	if q.queue[0] != "hue" {
		t.Error("Expected hue")
	}

	q.Pop()
	if len(q.queue) != 0 {
		t.Error("Expected empty array")
	}
}

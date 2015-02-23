package main

import "testing"
import "math"

type testMsg string

func (m testMsg) Read(p []byte) (n int, err error) {
	err = nil
	n = int(math.Min(float64(len(m)), float64(len(p))))

	for i := 0; i < n; i++ {
		p[i] = m[i]
	}

	return
}

func TestReadMessage(t *testing.T) {
	expected := "What's up?"
	m := testMsg(expected + "\r\n")

	msg, err := ReadMessage(m)

	if msg != expected || err != nil {
		t.Error("Expected message to be \"" + expected + "\".")
	}
}

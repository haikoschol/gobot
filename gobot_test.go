package main

import "testing"
import "math"

type testMsg struct {
	content   string
	position  int
	chunksize int
}

func (m testMsg) Read(p []byte) (n int, err error) {
	err = nil
	srcLen := float64(len(m.content) - m.position)
	dstLen := float64(len(p))

	if m.chunksize > 0 {
		n = int(math.Min(float64(m.chunksize), dstLen))
	} else {
		n = int(math.Min(srcLen, dstLen))
	}

	for i := m.position; i < n; i++ {
		p[i] = m.content[i]
	}

	m.position += n
	return
}

func makeMsg(s string, cs int) *testMsg {
	return &testMsg{
		content:   s + "\r\n",
		position:  0,
		chunksize: cs,
	}
}

func TestReadMessage(t *testing.T) {
	expected := "What's up?"
	m := makeMsg(expected, 0)

	msg, err := ReadMessage(m)

	if msg != expected || err != nil {
		t.Error("Expected message to be \"" + expected + "\".")
	}
}

func TestReadMessage_empty_message(t *testing.T) {
	msg, err := ReadMessage(makeMsg("", 0))
	if msg != "" || err != nil {
		t.Error("Expected empty message. Instead got:", msg)
	}
}

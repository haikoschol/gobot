package main

import "testing"
import "math"

type testMsg struct {
	content   string
	position  int
	chunksize int
}

func (m *testMsg) Read(p []byte) (n int, err error) {
	err = nil
	srcLen := float64(len(m.content) - m.position)
	dstLen := float64(len(p))

	if m.chunksize > 0 {
		tmp := math.Min(float64(m.chunksize), dstLen)
		n = int(math.Min(tmp, srcLen))
	} else {
		n = int(math.Min(srcLen, dstLen))
	}

	for i := 0; i < n; i++ {
		p[i] = m.content[m.position+i]
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

func assert(actual string, expected string, err error, t *testing.T) {
	if actual != expected || err != nil {
		t.Errorf("Expected message to be %q. Instead got %q. error: %v",
			expected, actual, err)
	}
}

func TestReadMessage(t *testing.T) {
	expected := "What's up?"
	m := makeMsg(expected, 0)

	msg, err := ReadMessage(m)

	assert(msg, expected, err, t)
}

func TestReadMessage_empty_message(t *testing.T) {
	msg, err := ReadMessage(makeMsg("", 0))
	assert(msg, "", err, t)
}

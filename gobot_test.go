package main

import "testing"
import (
	"fmt"
	"math"
)

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

func makeMsgOfLen(char byte, length int, cs int) *testMsg {
	str := makeStrOfLen(char, length)
	return makeMsg(str, cs)
}

func makeStrOfLen(char byte, length int) string {
	buf := make([]byte, length)

	for i := 0; i < length; i++ {
		buf[i] = char
	}

	return string(buf)
}

func assert(what string, actual string, expected string, err error,
	t *testing.T) {

	if actual != expected || err != nil {
		t.Errorf("Expected %s to be %q. Instead got %q. error: %v", what,
			expected, actual, err)
	}
}

func Test_readMessage(t *testing.T) {
	expected := "What's up?"
	m := makeMsg(expected, 0)

	msg, err := readMessage(m)

	assert("message", msg, expected, err, t)
}

func Test_readMessage_empty_message(t *testing.T) {
	msg, err := readMessage(makeMsg("", 0))
	assert("message", msg, "", err, t)
}

func Test_readMessage_partial_message(t *testing.T) {
	expected := "Hello"
	m := makeMsg(expected, 3)

	msg, err := readMessage(m)

	assert("message", msg, expected, err, t)
}

func Test_readMessage_longest_valid_message(t *testing.T) {
	expected := makeStrOfLen('a', maxMsgLen)
	m := makeMsg(expected, 0)

	msg, err := readMessage(m)

	assert("message", msg, expected, err, t)
}

func Test_readMessage_message_too_long(t *testing.T) {
	m := makeMsgOfLen('a', maxMsgLen+1, 0)

	_, err := readMessage(m)

	if err == nil {
		t.Error("ReadMessage() accepts messages longer than", maxMsgLen,
			"bytes.")
	}
}

func Test_parseMessage_nonempty_prefix_command_and_parameters(t *testing.T) {
	raw := ":Angel!wings@irc.org PRIVMSG Wiz :Are you receiving this message"

	msg, err := parseMessage(raw)

	if msg.prefix == "" || msg.command == "" || len(msg.parameters) == 0 ||
		err != nil {
		t.Error("Expected non-empty prefix, command and parameters. prefix:",
			msg.prefix, "command: ", msg.command, "parameters: ",
			msg.parameters)
	}
}

func Test_parseMessage_no_prefix(t *testing.T) {
	raw := "PRIVMSG Wiz :Are you receiving this message"

	msg, err := parseMessage(raw)

	if msg.prefix != "" || msg.command == "" || len(msg.parameters) == 0 ||
		err != nil {
		t.Error("Expected empty prefix and non-empty command and parameters. "+
			"prefix:", msg.prefix, "command: ", msg.command, "parameters: ",
			msg.parameters)
	}
}

func Test_parseMessage_prefix_marker_gets_removed(t *testing.T) {
	expected_prefix := "Angel!wings@irc.org"

	raw := fmt.Sprintf(":%s PRIVMSG Wiz Are you receiving this message",
		expected_prefix)

	msg, err := parseMessage(raw)

	assert("prefix", msg.prefix, expected_prefix, err, t)
}

func Test_dispatch_no_handlers(t *testing.T) {
	handlers := &CommandHandlers{}
	msg := &message{"", "foo", make([]string, 0)}

	handler, err := dispatch(msg, handlers)

	if handler != nil || err != nil {
		t.Error("Expected nil handler and no error. Istead got:", handler, err)
	}
}

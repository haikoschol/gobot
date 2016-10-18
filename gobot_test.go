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

func assertError(err error, t *testing.T) {
	if err == nil {
		t.Error("Expected error, got nil")
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

	assertError(err, t)
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

	assert("parsed message", msg.prefix, "", err, t)
	assert("parsed message", msg.command, "PRIVMSG", err, t)

	if len(msg.parameters) < 2 || err != nil {
		t.Error("Expected empty prefix and non-empty command and parameters. "+
			"prefix:", msg.prefix, "command:", msg.command, "parameters:",
			msg.parameters)
	}
}

func Test_parseMessage_prefix_marker_gets_removed(t *testing.T) {
	expected_prefix := "Angel!wings@irc.org"

	raw := fmt.Sprintf(":%s PRIVMSG Wiz :Are you receiving this message",
		expected_prefix)

	msg, err := parseMessage(raw)

	assert("prefix", msg.prefix, expected_prefix, err, t)
}

func Test_parseMessage_privmsg_recipient_and_body_in_parameters(t *testing.T) {
	expected_command := "PRIVMSG"
	expected_recipient := "Wiz"
	expected_body := "Are you receiving this message"

	raw := fmt.Sprintf("%s %s :%s", expected_command, expected_recipient,
		expected_body)

	msg, err := parseMessage(raw)

	assert("parsed message", msg.command, expected_command, err, t)
	assert("parsed message", msg.parameters[0], expected_recipient, err, t)
	assert("parsed message", msg.parameters[1], expected_body, err, t)
}

func runInvalidPrivmsgTest(msg string, t *testing.T) {
	_, err := parseMessage(msg)
	assertError(err, t)
}

func Test_parseMessage_invalid_privmsg_missing_body(t *testing.T) {
	runInvalidPrivmsgTest(":Angel!wings@irc.org PRIVMSG Wiz", t)
}

func Test_parseMessage_invalid_privmsg_missing_parameters(t *testing.T) {
	runInvalidPrivmsgTest(":Angel!wings@irc.org PRIVMSG", t)
}

func Test_parseMessage_invalid_privmsg_without_prefix_missing_parameters(t *testing.T) {
	runInvalidPrivmsgTest("PRIVMSG", t)
}

func Test_parseMessage_invalid_privmsg_prefix_marker_without_prefix(t *testing.T) {
	runInvalidPrivmsgTest(": PRIVMSG Wiz :Whassup?", t)
}

func Test_parseMessage_invalid_privmsg_empty_body(t *testing.T) {
	runInvalidPrivmsgTest(":Angel!wings@irc.org PRIVMSG Wiz :", t)
}

func Test_dispatch_no_handlers(t *testing.T) {
	handlers := &CommandHandlers{}
	msg := &message{"", "foo", make([]string, 0)}

	handler := dispatch(msg, handlers)

	if handler != nil {
		t.Error("Expected nil handler and no error. Instead got:", handler)
	}
}

func Test_dispatch_finds_handler(t *testing.T) {
	expectedResponse := "bar"
	expectedHandler := func(_ []string) string { return expectedResponse }
	handlers := &CommandHandlers{"foo": expectedHandler}
	msg := &message{"", "foo", make([]string, 0)}

	handler := dispatch(msg, handlers)

	if handler == nil || handler([]string{}) != expectedResponse {
		t.Error("dispatch() did not return the expected command handler.")
	}
}

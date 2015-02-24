package main

import "io"
import "errors"

const MaxMsgLen = 510

type Message struct {
	prefix     string
	command    string
	parameters []string
}

func isCompleteMessage(buf []byte) bool {
	l := len(buf)

	if l < 2 {
		return false
	}

	return buf[l-2] == '\r' && buf[l-1] == '\n'
}

func ReadMessage(reader io.Reader) (string, error) {
	buf := make([]byte, MaxMsgLen+2)
	pos := 0

	for len(buf[pos:]) != 0 {
		n, err := reader.Read(buf[pos:])
		pos += n

		if pos >= 2 && isCompleteMessage(buf[:pos]) {
			return string(buf[:pos-2]), err
		}
	}

	return "", errors.New("message too long")
}

func ParseMessage(_ string) (*Message, error) {
	return &Message{
		prefix:     "",
		command:    "",
		parameters: make([]string, 0),
	}, nil
}

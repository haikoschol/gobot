package main

import (
	"errors"
	"io"
	"strings"
)

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

func ParseMessage(raw string) (*Message, error) {
	parts := strings.Split(raw, " ")
	var prefix, command string
	var parameters []string

	if strings.HasPrefix(parts[0], ":") {
		prefix = parts[0][1:]
		command = parts[1]
		parameters = parts[2:]
	} else {
		prefix = ""
		command = parts[0]
		parameters = parts[1:]
	}

	return &Message{
		prefix:     prefix,
		command:    command,
		parameters: parameters,
	}, nil
}

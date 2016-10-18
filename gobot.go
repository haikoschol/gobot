package main

import (
	"errors"
	"io"
	"strings"
	"fmt"
)

const maxMsgLen = 510

type message struct {
	prefix     string
	command    string
	parameters []string
}

type CommandHandler func(parameters []string) string
type CommandHandlers map[string]CommandHandler

func isCompleteMessage(buf []byte) bool {
	l := len(buf)

	if l < 2 {
		return false
	}

	return buf[l-2] == '\r' && buf[l-1] == '\n'
}

func readMessage(reader io.Reader) (string, error) {
	buf := make([]byte, maxMsgLen+2)
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

func parseMessage(raw string) (*message, error) {
	var prefix, command, rest string
	var parameters []string

	if strings.HasPrefix(raw, ":") {
		parts := strings.SplitN(raw, " ", 3)
		if len(parts) != 3 {
			return nil, errors.New("message invalid: " + raw)
		}
		prefix, command, rest = parts[0][1:], parts[1], parts[2]
		if prefix == "" {
			return nil, errors.New("message invalid: " + raw)
		}
	} else {
		parts := strings.SplitN(raw, " ", 2)
		if len(parts) != 2 {
			return nil, errors.New("message invalid: " + raw)
		}
		command, rest = parts[0], parts[1]
	}

	if command == "PRIVMSG" {
		parameters = strings.SplitN(rest, " ", 2)
		if len(parameters) != 2 {
			return nil, errors.New("message invalid: " + raw)
		}
		parameters[1] = parameters[1][1:]
		if parameters[1] == "" {
			return nil, errors.New("message invalid: " + raw)
		}
	} else {
		parameters = strings.Split(rest, " ")
	}

	return &message{
		prefix:     prefix,
		command:    command,
		parameters: parameters,
	}, nil
}

func dispatch(msg *message, handlers *CommandHandlers) CommandHandler {
	return (*handlers)[msg.command]
}

func main() {
	fmt.Print("It's alive!\n")
}

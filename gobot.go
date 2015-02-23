package main

import "io"

func ReadMessage(reader io.Reader) (string, error) {
	msg := make([]byte, 12)
	n, err := reader.Read(msg)
	return string(msg[0 : n-2]), err
}

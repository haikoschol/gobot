package main

import "io"

func ReadMessage(reader io.Reader) (string, error) {
	buf := make([]byte, 12)
	pos := 0

	for {
		n, err := reader.Read(buf[pos:])
		pos += n

		if pos >= 2 && buf[pos-2] == '\r' && buf[pos-1] == '\n' {
			return string(buf[:pos-2]), err
		}
	}
}

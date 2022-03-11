package connection

import (
	"bytes"
	"errors"
	"io"
	"mjpeg_multiplexer/src/mjpeg"
	"net"
	"strconv"
)

type InputHTTP struct {
	connection net.Conn
}

func NewInputHTTP(url string) (source InputHTTP) {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		panic("Socket error")
	}

	return InputHTTP{conn}
}

func (source InputHTTP) Open() {
	_, err := source.connection.Write([]byte("GET /?action=stream HTTP/1.1\r\nHost:%s\r\n\r\n"))
	if err != nil {
		panic("cannot send GET request to " + source.connection.LocalAddr().String())
	}
}

func (source InputHTTP) ReceiveFrame() (mjpeg.Frame, error) {
	// todo optimize
	// Read header
	var buffer = make([]byte, 1)
	for {
		var bufferTmp = make([]byte, 1)
		var _, err = source.connection.Read(bufferTmp[:])
		if err != nil {
			println("Can't read from connection")
			return mjpeg.Frame{}, err
		}

		buffer = append(buffer, bufferTmp...)

		// Check if jpeg start has been reached
		if len(buffer) > len(mjpeg.JPEG_PREFIX) && bytes.Compare(buffer[len(buffer)-len(mjpeg.JPEG_PREFIX):], mjpeg.JPEG_PREFIX) == 0 {
			break
		}
	}

	// find index of content-length number
	var index = 0
	var wordIndex = 0
	var word = "Content-Length: "

	for {
		if index >= len(buffer) {
			break
		}

		if wordIndex == len(word) {
			break
		}

		if buffer[index] == word[wordIndex] {
			wordIndex++
		} else {
			wordIndex = 0
		}

		index++
	}

	if wordIndex != len(word) {
		return mjpeg.Frame{}, errors.New("empty frame")
	}

	// parse content size number
	var length = 0
	for {
		var el = buffer[index+length]
		if el < '0' || el > '9' {
			break
		}
		length++
	}
	var contentLength, _ = strconv.Atoi(string(buffer[index : index+length]))

	// read rest of the frame
	var bufferBody = make([]byte, contentLength-len(mjpeg.JPEG_PREFIX)) //jpge prefix has already been read
	var n, err = io.ReadFull(source.connection, bufferBody)

	if err != nil {
		println("Can't read from connection")
		return mjpeg.Frame{}, err
	}

	if n != contentLength-len(mjpeg.JPEG_PREFIX) {
		println("Cannot read all bytes")
		return mjpeg.Frame{}, errors.New("can't the expected amount of bytes")
	}

	var body = append(mjpeg.JPEG_PREFIX, bufferBody...)

	return mjpeg.Frame{Body: body}, nil
}

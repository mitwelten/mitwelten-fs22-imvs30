package connection

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mjpeg_multiplexer/src/mjpeg"
	"net"
	"strconv"
	"time"
)

type InputHTTP struct {
	url        string
	connection net.Conn
}

// NewInputHTTP todo TEST: Test this function by creating an input and checking if it runs
func NewInputHTTP(url string) *InputHTTP {
	return &InputHTTP{url: url}
}

func (source *InputHTTP) Start() {
	var err error
	source.connection, err = net.DialTimeout("tcp", source.url, 3*time.Second)
	if err != nil {
		panic("Socket error") // todo: error handling...
	}
	_, err = source.connection.Write([]byte("GET /?action=stream HTTP/1.1\r\nHost:%s\r\n\r\n"))
	if err != nil {
		panic("cannot send GET request to " + source.connection.LocalAddr().String())
	}
}

func (source *InputHTTP) ReceiveFrame() (mjpeg.MjpegFrame, error) {
	// todo optimize
	// Read header
	var buffer = make([]byte, 1)
	for {
		var bufferTmp = make([]byte, 1)
		var _, err = source.connection.Read(bufferTmp[:])
		if err != nil {
			return mjpeg.MjpegFrame{}, fmt.Errorf("error while reading header: %w", err)
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
		return mjpeg.MjpegFrame{}, errors.New("empty frame")
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
		return mjpeg.MjpegFrame{}, fmt.Errorf("error while reading frame: %w", err)
	}

	if n != contentLength-len(mjpeg.JPEG_PREFIX) {
		println("Cannot read all bytes")
		return mjpeg.MjpegFrame{}, fmt.Errorf("could not read expected amount of bytes: %w", err)
	}

	var body = append(mjpeg.JPEG_PREFIX, bufferBody...)

	return mjpeg.MjpegFrame{Body: body}, nil
}

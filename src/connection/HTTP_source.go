package connection

import (
	"bytes"
	"errors"
	"io"
	"mjpeg_multiplexer/src/mjpeg"
	"net"
	"strconv"
)

type HTTPSource struct {
	connection net.Conn
}

func NewHTTPSource(url string) (source HTTPSource) {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		panic("Socket error")
	}

	return HTTPSource{conn}
}

func (source HTTPSource) Open() {
	source.connection.Write([]byte("GET /?action=stream HTTP/1.1\r\nHost:%s\r\n\r\n"))
}

func (source HTTPSource) ReceiveFrame() (mjpeg.Frame, error) {
	//todo optimize

	// Read header
	var buffer = make([]byte, 1)
	for {
		var buffer_tmp = make([]byte, 1)
		var _, err = source.connection.Read(buffer_tmp[:])
		if err != nil {
			println("Can't read from connection")
			panic(err)
		}

		buffer = append(buffer, buffer_tmp...)

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
	var content_length, _ = strconv.Atoi(string(buffer[index : index+length]))

	// read rest of the frame
	var buffer_body = make([]byte, content_length-len(mjpeg.JPEG_PREFIX)) //jpge prefix has already been read
	var n, err = io.ReadFull(source.connection, buffer_body)

	if err != nil {
		println("Can't read from connection")
		panic(err)
	}
	if n != content_length-len(mjpeg.JPEG_PREFIX) {
		println(n)
		println(content_length)
		panic("Cannot read all bytes")
	}

	var body = append(mjpeg.JPEG_PREFIX, buffer_body...)

	return mjpeg.Frame{Header: buffer[:len(buffer)-len(mjpeg.JPEG_PREFIX)], Body: body}, nil
}

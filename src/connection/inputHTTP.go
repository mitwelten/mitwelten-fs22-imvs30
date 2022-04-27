package connection

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"mjpeg_multiplexer/src/customErrors"
	"mjpeg_multiplexer/src/mjpeg"
	"net"
	"strconv"
	"strings"
	"time"
)

type InputHTTP struct {
	url                string
	connection         net.Conn
	bufferedConnection *bufio.Reader
}

// NewInputHTTP todo TEST: Test this function by creating an input and checking if it runs
func NewInputHTTP(url string) *InputHTTP {
	return &InputHTTP{url: url}
}

func (source *InputHTTP) Info() string {
	return source.url
}
func (source *InputHTTP) open() error {
	var err error
	source.connection, err = net.DialTimeout("tcp", source.url, 3*time.Second)
	source.bufferedConnection = bufio.NewReader(source.connection)

	if err != nil {
		return &customErrors.ErrHttpOpenInputSocketDial{IP: source.url}
	}

	return nil
}

func (source *InputHTTP) sendHeader() error {
	var err error

	_, err = source.connection.Write([]byte("GET /?action=stream HTTP/1.1\r\nHost:%s\r\n\r\n"))
	if err != nil {
		return &customErrors.ErrHttpWriteHeader{IP: source.connection.LocalAddr().String()}
	}

	return nil
}
func (source *InputHTTP) Start() error {
	var err error

	err = source.open()
	if err != nil {
		return err
	}

	err = source.sendHeader()
	if err != nil {
		return err
	}

	return nil
}

func (source *InputHTTP) ReceiveFrameFast() (mjpeg.MjpegFrame, error) {
	header, err := source.bufferedConnection.ReadString(mjpeg.JPEG_PREFIX[0])
	if err != nil {
		// could not read from header
		return mjpeg.MjpegFrame{}, err
	}

	field := "Content-Length: "
	startIndex := strings.LastIndex(header, field)
	if startIndex == -1 {
		//invalid header: no content length
		return mjpeg.MjpegFrame{}, err
	}

	// count n digits after field
	var length = 0
	for {
		var el = header[startIndex+len(field)+length]
		if el < '0' || el > '9' {
			break
		}
		length++
	}
	//Content-Length: 301
	contentLengthStart := startIndex + len(field)
	contentLengthEnd := contentLengthStart + length
	contentLength, err := strconv.Atoi(header[contentLengthStart:contentLengthEnd])
	if err != nil {
		// cant parse content length
		return mjpeg.MjpegFrame{}, err
	}

	body := make([]byte, contentLength-1) // first byte of jpge prefix has already been read

	n, err := io.ReadFull(source.bufferedConnection, body)
	if n != contentLength-1 {
		log.Println("error: cannot read all bytes")
		return mjpeg.MjpegFrame{}, &customErrors.ErrHttpReadEntireFrame{}
	}

	return mjpeg.MjpegFrame{Body: append(mjpeg.JPEG_PREFIX[0:1], body...)}, nil
}

func (source *InputHTTP) ReceiveFrame() (mjpeg.MjpegFrame, error) {
	// todo optimize
	// Read header
	var buffer = make([]byte, 1)
	for {
		var bufferTmp = make([]byte, 1)
		var _, err = source.connection.Read(bufferTmp[:])
		if err != nil {
			return mjpeg.MjpegFrame{}, &customErrors.ErrHttpReadHeader{}
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
		log.Println("error: empty frame received")
		return mjpeg.MjpegFrame{}, &customErrors.ErrHttpEmptyFrame{}
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
		log.Println("error: could not read frame")
		return mjpeg.MjpegFrame{}, &customErrors.ErrHttpReadFrame{}
	}

	if n != contentLength-len(mjpeg.JPEG_PREFIX) {
		log.Println("error: cannot read all bytes")
		return mjpeg.MjpegFrame{}, &customErrors.ErrHttpReadEntireFrame{}
	}

	var body = append(mjpeg.JPEG_PREFIX, bufferBody...)

	return mjpeg.MjpegFrame{Body: body}, nil
}

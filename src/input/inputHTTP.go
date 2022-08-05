package input

import (
	"bufio"
	"io"
	"log"
	"mjpeg_multiplexer/src/customErrors"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/imageUtils"
	"mjpeg_multiplexer/src/mjpeg"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	header = "GET /?action=stream HTTP/1.1\r\n" +
		"Host:%s\r\n"
	delim          = "\r\n"
	field          = "Content-Length: "
	authentication = "Authorization: Basic "
	NewMJPEGFrameNewMJPEGFrame
)

var JPEG_PREFIX = []byte("\xff\xd8")

type InputHTTP struct {
	data               InputData
	configIndex        int
	url                string
	connection         net.Conn
	bufferedConnection *bufio.Reader
}

// NewInputHTTP ctor
func NewInputHTTP(configIndex int, url string) *InputHTTP {
	return &InputHTTP{configIndex: configIndex, url: url, data: InputData{InputStorage: mjpeg.NewFrameStorage()}}
}

func (source *InputHTTP) GetInputData() *InputData {
	return &source.data
}
func (source *InputHTTP) GetInfo() string {
	return source.url
}

// open tries to open a TCP connection to the sources URL
func (source *InputHTTP) open() error {
	var err error
	source.connection, err = net.DialTimeout("tcp", source.url, 3*time.Second)
	source.bufferedConnection = bufio.NewReader(source.connection)

	if err != nil {
		return &customErrors.ErrHttpOpenInputSocketDial{IP: source.url}
	}

	return nil
}

// sendHeader tries to send open the connection by sending a header to the contents of a mjpeg-streamer stream.
func (source *InputHTTP) sendHeader() error {
	var err error

	_, err = source.connection.Write([]byte(header))
	if err != nil {
		return &customErrors.ErrHttpWriteHeader{IP: source.connection.LocalAddr().String()}
	}

	// Also send the authentication if available
	if global.Config.Debug {
		log.Printf("InputConfigs is %+v\n", global.Config.InputConfigs[source.configIndex])
	}
	if global.Config.UseAuth && global.Config.InputConfigs[source.configIndex].Authentication != "" {
		if global.Config.Debug {
			log.Printf("DEBUG: Sending authenticaion to input source %v\n", source.url)
		}
		_, err = source.connection.Write([]byte(authentication + global.Config.InputConfigs[source.configIndex].Authentication + delim))
		if err != nil {
			return &customErrors.ErrHttpWriteHeader{IP: source.connection.LocalAddr().String()}
		}
	}

	_, err = source.connection.Write([]byte(delim))
	if err != nil {
		return &customErrors.ErrHttpWriteHeader{IP: source.connection.LocalAddr().String()}
	}

	// Get the first frame to test if we have permission to access the source
	frame, err := source.ReceiveFrame()

	if err != nil {
		return &customErrors.ErrHttpOpenInputAuthentication{Text: err.Error()}
	}

	source.data.InputStorage.Store(&frame)
	imageUtils.Decode(source.data.InputStorage)

	return nil
}

// Init opens connection + sends header
func (source *InputHTTP) Init() error {
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

// ReceiveFrame reads a mjpeg-stream and parses the next received frame as mjpeg.MjpegFrame
func (source *InputHTTP) ReceiveFrame() (mjpeg.MjpegFrame, error) {
	header, err := source.bufferedConnection.ReadString(JPEG_PREFIX[0])
	if err != nil {
		return mjpeg.MjpegFrame{}, err
	}

	startIndex := strings.LastIndex(header, field)
	if startIndex == -1 {
		return mjpeg.MjpegFrame{}, &customErrors.ErrInvalidFrame{Text: "invalid header: no content length"}
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
	contentLengthStart := startIndex + len(field)
	contentLengthEnd := contentLengthStart + length
	contentLength, err := strconv.Atoi(header[contentLengthStart:contentLengthEnd])
	if err != nil {
		// cant parse content length
		return mjpeg.MjpegFrame{}, &customErrors.ErrInvalidFrame{Text: "cant parse content length"}
	}

	body := make([]byte, contentLength-1) // first byte of jpeg prefix has already been read

	n, err := io.ReadFull(source.bufferedConnection, body)
	if n != contentLength-1 {
		return mjpeg.MjpegFrame{}, &customErrors.ErrInvalidFrame{Text: "cannot read all bytes"}
	}

	return mjpeg.MjpegFrame{Body: append(JPEG_PREFIX[0:1], body...)}, nil
}

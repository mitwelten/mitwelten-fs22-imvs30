package connection

import (
	_ "embed"
	"mjpeg_multiplexer/src/mjpeg"
	"net"
)

// This is only used for integration testing

//go:embed red.jpg
var redJPG []byte

//go:embed blue.jpg
var blueJPG []byte

func RedFrame() mjpeg.Frame {
	return mjpeg.Frame{Body: redJPG}
}

func BlueFrame() mjpeg.Frame {
	return mjpeg.Frame{Body: blueJPG}
}

func SimpleServer(port string, frame mjpeg.Frame) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	conn, err := listener.Accept()
	if err != nil {
		return err
	}

	var server = Server{nil, conn}
	err = server.sendHeader()
	if err != nil {
		return err
	}
	err = server.sendFrame(frame)
	if err != nil {
		return err
	}

	return nil
}

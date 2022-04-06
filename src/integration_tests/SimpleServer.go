package integration_tests

import (
	_ "embed"
	"mjpeg_multiplexer/src/connection"
	"mjpeg_multiplexer/src/mjpeg"
	"net"
)

// This is only used for integration testing
//go:embed red.jpg
var redJPG []byte

//go:embed blue.jpg
var blueJPG []byte

func RedFrame() mjpeg.MjpegFrame {
	return mjpeg.MjpegFrame{Body: redJPG}
}

func BlueFrame() mjpeg.MjpegFrame {
	return mjpeg.MjpegFrame{Body: blueJPG}
}

func SimpleServer(port string, frame mjpeg.MjpegFrame) error {
	go func() {
		println("Listening...")
		listener, err := net.Listen("tcp", ":"+port)

		conn, err := listener.Accept()
		if err != nil {
			panic(err.Error())
		}

		var server = connection.ClientConnection{Connection: conn}
		err = server.SendHeader()
		if err != nil {
			panic(err.Error())
		}
		err = server.SendFrame(frame)
		if err != nil {
			panic(err.Error())
		}

	}()

	return nil
}

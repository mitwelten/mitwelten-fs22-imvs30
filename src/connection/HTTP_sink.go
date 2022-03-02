package connection

import (
	"mjpeg_multiplexer/src/mjpeg"
	"net"
	"strconv"
)

var Connections = make([]net.Conn, 0)

type HTTPSink struct {
}

func sendHeader(connection net.Conn) {
	var header = "HTTP/1.0 200 OK\r\n" +
		"Access-Control-Allow-Origin: *\r\n" +
		"Content-Type: multipart/x-mixed-replace;boundary=boundarydonotcross\r\n" +
		"\r\n" +
		"--boundarydonotcross\r\n"
	_, err := connection.Write([]byte(header))
	if err != nil {
		panic("Can't send header")
	}
}
func sendFrame(connection net.Conn, frame mjpeg.Frame) {
	var header = "Content-Type: image/jpeg\r\n" +
		"Content-Length: " + strconv.Itoa(len(frame.Body)) + "\r\n"
	//"X-Timestamp: TODO \r\n"

	println(header)
	_, err := connection.Write([]byte(header))
	if err != nil {
		panic("Can't send header")
	}
	_, err = connection.Write(frame.Body)
	if err != nil {
		panic("Can't send frame")
	}
	_, err = connection.Write([]byte("\r\n--boundarydonotcross\r\n"))
	if err != nil {
		panic("Can't send separator")
	}
}

func NewHTTPSink(port string) HTTPSink {
	//todo this is trash
	listener, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		panic("Socket error")
	}

	go func() {
		for {
			conn, err := listener.Accept()
			println(conn.LocalAddr(), " connected!")
			if err != nil {
				panic("Invalid connection or w/e")
			}
			//TODO locks etc
			Connections = append(Connections, conn)
			sendHeader(conn)
		}
	}()

	return HTTPSink{}
}

func (sink HTTPSink) ProcessFrame(frame mjpeg.Frame) {
	println(len(Connections))
	for _, connection := range Connections {
		sendFrame(connection, frame)
	}
}

package connection

import (
	"errors"
	"mjpeg_multiplexer/src/mjpeg"
	"net"
	"strconv"
	"sync"
)

type OutputHTTP struct {
}

var servers = make([]Server, 0)
var serversMutex sync.Mutex

type Server struct {
	channel    chan mjpeg.Frame
	connection net.Conn
}

var HEADER = "HTTP/1.1 200 OK\r\n" +
	"Server: mjpeg-multiplexer\r\n" +
	"Connection: close\r\n" +
	"Max-Age: 0\r\n" +
	"Expires: 0\r\n" +
	"Cache-Control: no-cache, private\r\n" +
	"Pragma: no-cache\r\n" +
	"Content-Type: multipart/x-mixed-replace; boundary=--boundarydonotcross\r\n" +
	"\r\n" +
	"--boundarydonotcross\r\n"

var DELIM = "\r\n----boundarydonotcross\r\n"

func remove(server Server) {
	//remove this server from the list of servers
	serversMutex.Lock()
	for i, s := range servers {
		if server == s {
			servers = append(servers[:i], servers[i+1:]...)
			break
		}
	}
	serversMutex.Unlock()
}

func serve(server Server) {
	defer func(connection net.Conn) {
		err := connection.Close()
		if err != nil {
			println("can't close connection to " + connection.LocalAddr().String() + ", potential leak!")
		}
	}(server.connection)

	defer close(server.channel)

	var err = server.sendHeader()
	if err != nil {
		println("error when sending header to " + server.connection.LocalAddr().String() + ", closing connection")
		println(err.Error())
		remove(server)
		return
	}

	for {
		var frame = <-server.channel
		var err = server.sendFrame(frame)
		if err != nil {
			//todo Counter that closes after X errors
			println("error when sending frame to " + server.connection.LocalAddr().String() + ", closing connection")
			println(err.Error())
			remove(server)
			return
		}
	}
}
func (server Server) sendHeader() error {
	var header = HEADER
	_, err := server.connection.Write([]byte(header))
	if err != nil {
		return err
	}
	return nil
}
func (server Server) sendFrame(frame mjpeg.Frame) error {
	//Format must be not be change, else it will not work on some browsers!
	var header = "Content-Type: image/jpeg\r\n" +
		"Content-Length: " + strconv.Itoa(len(frame.Body)) + "\r\n" +
		"\r\n"

	var data = []byte(header)
	data = append(data, frame.Body...)
	data = append(data, []byte(DELIM)...)

	_, err := server.connection.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func NewOutputHTTP(port string) (Output, error) {
	//todo this is trash
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return OutputHTTP{}, errors.New("can't open socket")
	}

	go func() {
		for {
			conn, err := listener.Accept()
			println(conn.LocalAddr(), " connected!")
			if err != nil {
				println("Invalid connection")
				continue
			}

			var server = Server{make(chan mjpeg.Frame), conn}

			serversMutex.Lock()
			servers = append(servers, server)
			serversMutex.Unlock()

			go serve(server)
		}
	}()

	return OutputHTTP{}, nil
}

func (sink OutputHTTP) SendFrame(frame mjpeg.Frame) error {
	for _, server := range servers {
		server.channel <- frame
	}
	return nil
}

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

var delim = "--boundarydonotcross"

func serve(server Server) {
	defer server.connection.Close()
	defer close(server.channel)

	var err = server.sendHeader()
	if err != nil {
		println(err.Error())
		return
	}

	for {
		var frame = <-server.channel
		var err = server.sendFrame(frame)
		if err != nil {
			print(err.Error())
			serversMutex.Lock()
			//remove this server from the list of servers
			for i, s := range servers {
				if server == s {
					servers = append(servers[:i], servers[i+1:]...)
					break
				}
			}
			serversMutex.Unlock()
			return
		}
	}
}
func (server Server) sendHeader() error {
	var header = "HTTP/1.0 200 OK\r\n" +
		"Access-Control-Allow-Origin: *\r\n" +
		"Content-Type: multipart/x-mixed-replace;boundary=" + delim + "\r\n" +
		"\r\n" +
		"--" + delim + "\r\n"
	_, err := server.connection.Write([]byte(header))
	if err != nil {
		return errors.New("Can't send header")
	}
	return nil
}
func (server Server) sendFrame(frame mjpeg.Frame) error {
	var header = "Content-Type: image/jpeg\r\n" +
		"Content-Length: " + strconv.Itoa(len(frame.Body)) + "\r\n"
	//"X-Timestamp: TODO \r\n"

	_, err := server.connection.Write([]byte(header))
	if err != nil {
		return errors.New("can't send header")
	}

	_, err = server.connection.Write(frame.Body)
	if err != nil {
		return errors.New("can't send frame body")
	}
	_, err = server.connection.Write([]byte("\r\n--" + delim + "\r\n"))
	if err != nil {
		return errors.New("Can't send delim")
	}

	return nil
}

func NewHTTPSink(port string) OutputHTTP {
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

	return OutputHTTP{}
}

func (sink OutputHTTP) ProcessFrame(frame mjpeg.Frame) {
	for _, server := range servers {
		server.channel <- frame
	}
}

package connection

import (
	"errors"
	"mjpeg_multiplexer/src/communication"
	"mjpeg_multiplexer/src/mjpeg"
	"net"
	"strconv"
	"sync"
)

type OutputHTTP struct {
}

var clients = make([]ClientConnection, 0)
var clientsMutex sync.Mutex

type ClientConnection struct {
	channel    chan mjpeg.MjpegFrame
	connection net.Conn
	isClosed   bool
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

func remove(client ClientConnection) {
	//remove this SimpleServer from the list of clients
	for i, s := range clients {
		if client == s {
			//remove this client from the client list
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

func serve(client ClientConnection) {
	defer func(client_ ClientConnection) {
		//safely remove client from client list and close its channel
		println("Trying to remove...")
		remove(client_)
		<-client_.channel
		close(client_.channel)

		println("Successfully removed!!")

		err := client_.connection.Close()
		if err != nil {
			println("can't close connection to " + client_.connection.LocalAddr().String() + ", potential leak!")
		}
	}(client)

	var err = client.sendHeader()
	if err != nil {
		println("error when sending header to " + client.connection.LocalAddr().String() + ", closing connection")
		println(err.Error())
		return
	}

	for {
		var frame = <-client.channel
		var err = client.sendFrame(frame)
		if err != nil {
			//todo Counter that closes after X errors
			println("error when sending frame to " + client.connection.LocalAddr().String() + ", closing connection")
			println(err.Error())
			return
		}
	}
}
func (client ClientConnection) sendHeader() error {
	var header = HEADER
	_, err := client.connection.Write([]byte(header))
	if err != nil {
		return err
	}
	return nil
}
func (client ClientConnection) sendFrame(frame mjpeg.MjpegFrame) error {
	//Format must be not be change, else it will not work on some browsers!
	var header = "Content-Type: image/jpeg\r\n" +
		"Content-Length: " + strconv.Itoa(len(frame.Body)) + "\r\n" +
		"\r\n"

	var data = []byte(header)
	data = append(data, frame.Body...)
	data = append(data, []byte(DELIM)...)

	_, err := client.connection.Write(data)
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
			println(conn.RemoteAddr().String(), " connected!")

			if err != nil {
				println("Invalid connection")
				continue
			}

			client := ClientConnection{make(chan mjpeg.MjpegFrame), conn, false}

			clients = append(clients, client)

			go serve(client)
		}
	}()

	return OutputHTTP{}, nil
}

func (output OutputHTTP) SendFrame(frame mjpeg.MjpegFrame) error {
	//TODO CRITICAL: FIX DATARACE
	for _, client := range clients {
		client.channel <- frame
	}

	return nil
}
func (output OutputHTTP) Run(storage *communication.FrameStorage) {
	go func(storage_ *communication.FrameStorage) {
		for {
			frame := storage_.Get()
			err := output.SendFrame(frame)
			if err != nil {
				println("Error while trying to send frame to output")
				println(err.Error())
				continue
			}
		}
	}(storage)
}

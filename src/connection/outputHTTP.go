package connection

import (
	"errors"
	"log"
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/mjpeg"
	"net"
	"strconv"
	"sync"
)

type OutputHTTP struct {
	lastFrame mjpeg.MjpegFrame
}

var clients = make([]ClientConnection, 0)
var clientsMutex sync.RWMutex

type ClientConnection struct {
	channel    chan mjpeg.MjpegFrame
	Connection net.Conn
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

var DELIM = "\r\n--boundarydonotcross\r\n"

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

func (output *OutputHTTP) serve(client ClientConnection) {
	// On disconnect, close connection and cleanup
	defer func(client_ ClientConnection) {
		//safely remove client from client list and close its channel
		clientsMutex.Lock()
		remove(client_)
		close(client_.channel)
		clientsMutex.Unlock()

		err := client_.Connection.Close()
		if err != nil {
			log.Println("can't close Connection to " + client_.Connection.LocalAddr().String() + ", potential leak!")
		}
	}(client)

	// Send the stream header
	var err = client.SendHeader()
	if err != nil {
		log.Println("error when sending header to " + client.Connection.LocalAddr().String() + ", closing Connection")
		log.Println(err.Error())
		return
	}

	// Send the cached frame to the client
	_ = client.SendFrame(output.lastFrame)

	// Know issue in chromium: Chromium's stream always lags one frame behind, to show the first frame immediately it is sent twice here
	// reference: chromium bug tracker issue #527446
	// status: open
	// link: https://bugs.chromium.org/p/chromium/issues/detail?id=527446
	_ = client.SendFrame(output.lastFrame)

	// Send all receive frames
	for {
		var frame = <-client.channel
		var err = client.SendFrame(frame)
		if err != nil {
			//todo Counter that closes after X errors
			log.Println("error when sending frame to " + client.Connection.LocalAddr().String() + ", closing Connection")
			log.Println(err.Error())
			return
		}
	}
}
func (client ClientConnection) SendHeader() error {
	var header = HEADER
	_, err := client.Connection.Write([]byte(header))
	if err != nil {
		return err
	}
	return nil
}
func (client ClientConnection) SendFrame(frame mjpeg.MjpegFrame) error {
	//Format must be not be changed, else it will not work on some browsers!
	var header = "Content-Type: image/jpg\r\n" +
		"Content-Length: " + strconv.Itoa(len(frame.Body)) + "\r\n" +
		"\r\n"

	data := []byte(header)
	data = append(data, frame.Body...)
	data = append(data, []byte(DELIM)...)

	_, err := client.Connection.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func NewOutputHTTP(port string) (Output, error) {
	//todo this is trash
	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		return &OutputHTTP{}, errors.New("can't open socket")
	}

	output := &OutputHTTP{}
	go func() {
		for {
			conn, err := listener.Accept()
			log.Println(conn.RemoteAddr().String(), " connected!")

			if err != nil {
				log.Println("Invalid Connection")
				continue
			}

			client := ClientConnection{make(chan mjpeg.MjpegFrame), conn, false}

			clientsMutex.Lock()
			clients = append(clients, client)
			clientsMutex.Unlock()

			go output.serve(client)
		}
	}()

	return output, nil
}

func (output *OutputHTTP) SendFrame(frame mjpeg.MjpegFrame) error {
	defer clientsMutex.RUnlock()
	clientsMutex.RLock()

	for _, client := range clients {
		select {
		case client.channel <- frame:
		default:
		}
	}

	return nil
}

func (output *OutputHTTP) Run(aggregator aggregator.Aggregator) {

	lock := sync.Mutex{}
	lock.Lock()
	condition := sync.NewCond(&lock)

	aggregator.GetAggregatorData().OutputCondition = condition

	go func(storage_ *mjpeg.FrameStorage) {
		for {
			condition.Wait()
			frame := storage_.GetLatest()

			/*
				if reflect.DeepEqual(frame, previousFrame) {
					continue
				}
			*/

			output.lastFrame = frame

			err := output.SendFrame(frame)
			if err != nil {
				log.Printf("Error while trying to send frame to output: %s\n", err.Error())
				continue
			}
		}
	}(aggregator.GetAggregatorData().OutputStorage)
}

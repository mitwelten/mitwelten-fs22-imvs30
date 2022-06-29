package connection

import (
	"log"
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/mjpeg"
	"net"
	"strconv"
	"sync"
)

type OutputHTTP struct {
	lastFrame    *mjpeg.MjpegFrame
	clients      []ClientConnection
	clientsMutex *sync.RWMutex
	aggregator   aggregator.Aggregator
	condition    *sync.Cond
	errorCounter int
}

type ClientConnection struct {
	channel    chan *mjpeg.MjpegFrame
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

func NewOutputHTTP(port string, aggregator aggregator.Aggregator) Output {
	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatalf("can't open socket %s: %s", port, err.Error())
	}

	output := OutputHTTP{}
	output.aggregator = aggregator

	lock := sync.Mutex{}
	lock.Lock()
	output.condition = sync.NewCond(&lock)
	output.aggregator.GetAggregatorData().OutputCondition = output.condition

	output.clients = make([]ClientConnection, 0)
	output.clientsMutex = &sync.RWMutex{}
	output.errorCounter = 0

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Invalid connection: %s", err.Error())
				continue
			}

			log.Printf("%s connected\n", conn.RemoteAddr().String())

			client := ClientConnection{make(chan *mjpeg.MjpegFrame), conn, false}

			output.clientsMutex.Lock()
			output.clients = append(output.clients, client)
			if len(output.clients) == 1 && output.aggregator != nil {
				output.aggregator.GetAggregatorData().Enabled = true
				log.Printf("Client connected, starting aggregator\n")
			}
			output.clientsMutex.Unlock()

			go output.serve(client)
		}
	}()

	return &output
}

func (output *OutputHTTP) SendFrame(frame *mjpeg.MjpegFrame) {
	defer output.clientsMutex.RUnlock()
	output.clientsMutex.RLock()

	for _, client := range output.clients {
		select {
		case client.channel <- frame:
		default:
		}
	}
}

func (output *OutputHTTP) remove(client ClientConnection) {
	//remove this SimpleServer from the list of clients
	for i, s := range output.clients {
		if client == s {
			//remove this client from the client list
			output.clients = append(output.clients[:i], output.clients[i+1:]...)
			break
		}
	}
}

func (output *OutputHTTP) serve(client ClientConnection) {
	// On disconnect, close connection and cleanup
	defer func(client_ ClientConnection) {
		//safely remove client from client list and close its channel
		output.clientsMutex.Lock()
		output.remove(client_)
		close(client_.channel)
		if len(output.clients) == 0 && output.aggregator != nil {
			output.aggregator.GetAggregatorData().Enabled = false
			log.Printf("No more clients, stopping aggregator\n")
		}
		output.clientsMutex.Unlock()

		err := client_.Connection.Close()
		if err != nil {
			log.Printf("Can't close Connection to %s\n", client_.Connection.LocalAddr().String())
		}
	}(client)

	// Send the stream header
	var err = client.SendHeader()
	if err != nil {
		log.Printf("Closing connection, error when sending header to %s: %s\n", client.Connection.LocalAddr().String(), err.Error())
		return
	}

	if output.lastFrame != nil {
		// Send the cached frame to the client
		_ = client.SendFrame(output.lastFrame)

		// Know issue in chromium: Chromium's stream always lags one frame behind, to show the first frame immediately it is sent twice here
		// reference: chromium bug tracker issue #527446
		// status: open
		// link: https://bugs.chromium.org/p/chromium/issues/detail?id=527446
		_ = client.SendFrame(output.lastFrame)
	}

	// Send all receive frames
	for {
		var frame = <-client.channel
		var err = client.SendFrame(frame)
		if err != nil {
			output.errorCounter++
			if output.errorCounter == 5 {
				log.Printf("Closing connection, error while sending frames to %s: %s\n", client.Connection.LocalAddr().String(), err.Error())
				return
			}
			continue
		}
		output.errorCounter = 0
	}
}
func (client *ClientConnection) SendHeader() error {
	var header = HEADER
	_, err := client.Connection.Write([]byte(header))
	if err != nil {
		return err
	}
	return nil
}
func (client *ClientConnection) SendFrame(frame *mjpeg.MjpegFrame) error {
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

func (output *OutputHTTP) Run() {
	go func(storage_ *mjpeg.FrameStorage) {
		for {
			output.condition.Wait()

			frame := storage_.GetLatestPtr()
			output.lastFrame = frame
			output.SendFrame(frame)
		}
	}(output.aggregator.GetAggregatorData().OutputStorage)
}

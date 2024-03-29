package output

import (
	"fmt"
	"log"
	"mjpeg_multiplexer/src/aggregator"
	"mjpeg_multiplexer/src/global"
	"mjpeg_multiplexer/src/mjpeg"
	"net"
	"sync"
)

type OutputHTTP struct {
	//port is the tcp port of the output
	port string
	//lastFrame is the last frame sent which may be resent in certain conditions
	lastFrame *mjpeg.MjpegFrame
	//clients is the list of connect tcp connections
	clients []ClientConnection
	//clientsMutex is the mutex which must be used when editing the clients list
	clientsMutex *sync.RWMutex
	//aggregator provides the frames
	aggregator aggregator.Aggregator
	//condition is used to wait for new frames
	condition *sync.Cond
}

type ClientConnection struct {
	channel    chan *mjpeg.MjpegFrame
	Connection net.Conn
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

var CONTENT = "Content-Type: image/jpeg\r\n" +
	"Content-Length: %d\r\n" +
	"\r\n"

var DELIM = "\r\n--boundarydonotcross\r\n"

// NewOutputHTTP ctor
func NewOutputHTTP(port string) Output {
	output := OutputHTTP{}
	output.port = port

	output.clients = make([]ClientConnection, 0)
	output.clientsMutex = &sync.RWMutex{}

	return &output
}

// SendFrame sends a single frame to all connected outputs
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

// remove removes a client from the list of connected clients
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

// serve handles the connection to a client by sending the header and all received frames
// to be started in a new go routine
func (output *OutputHTTP) serve(client ClientConnection) {
	// On disconnect, close connection and cleanup
	defer func(client_ ClientConnection) {
		//safely remove client from client list and close its channel
		output.clientsMutex.Lock()
		output.remove(client_)
		close(client_.channel)
		if len(output.clients) == 0 && output.aggregator != nil {
			global.Config.AggregatorMutex.Lock()
			global.Config.AggregatorEnabled = false
			global.Config.AggregatorMutex.Unlock()
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
	errorCounter := 0
	for {
		var frame = <-client.channel
		var err = client.SendFrame(frame)
		if err != nil {
			errorCounter++
			if errorCounter == 5 {
				log.Printf("Closing connection, error while sending frames to %s: %s\n", client.Connection.LocalAddr().String(), err.Error())
				return
			}
			continue
		}
		errorCounter = 0
	}
}

// SendHeader sends the MJPEG-header to the client
func (client *ClientConnection) SendHeader() error {
	var header = HEADER
	_, err := client.Connection.Write([]byte(header))
	if err != nil {
		return err
	}
	return nil
}

// SendFrame sends a single MJPEG-Frame to the client
func (client *ClientConnection) SendFrame(frame *mjpeg.MjpegFrame) error {
	//Format must be not be changed, else it will not work on some browsers!
	content := fmt.Sprintf(CONTENT, len(frame.Body))

	data := []byte(content)
	data = append(data, frame.Body...)
	data = append(data, []byte(DELIM)...)

	_, err := client.Connection.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// connectionLoop listens to the specified tcp port and accepts new clients, before sending them to the serve method
func (output *OutputHTTP) connectionLoop() {
	listener, err := net.Listen("tcp", ":"+output.port)

	if err != nil {
		log.Fatalf("can't open socket %s: %s", output.port, err.Error())
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			//can't open connection => ignore
			log.Printf("Invalid connection: %s", err.Error())
			continue
		}

		log.Printf("%s connected\n", conn.RemoteAddr().String())

		//handle new connection
		client := ClientConnection{make(chan *mjpeg.MjpegFrame), conn}
		output.clientsMutex.Lock()
		output.clients = append(output.clients, client)
		if len(output.clients) == 1 && output.aggregator != nil {
			global.Config.AggregatorMutex.Lock()
			global.Config.AggregatorEnabled = true
			global.Config.AggregatorMutex.Unlock()
			log.Printf("Client connected, starting aggregator\n")
		}
		output.clientsMutex.Unlock()

		go output.serve(client)
	}
}

func (output *OutputHTTP) distributeFramesLoop() {
	for {
		output.condition.Wait()
		frame := output.aggregator.GetAggregatorData().AggregatorStorage.GetFrame()
		output.lastFrame = frame
		output.SendFrame(frame)
	}
}

// StartOutput starts the output by waiting for new connections and distributing new frames to them
func (output *OutputHTTP) StartOutput(agg *aggregator.Aggregator) {
	output.aggregator = *agg
	output.condition = mjpeg.CreateUpdateCondition(output.aggregator.GetAggregatorData().AggregatorStorage)

	go output.connectionLoop()
	go output.distributeFramesLoop()
}

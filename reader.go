package main

import (
	"bufio"
	"bytes"
	"net"
	"os"
	"sync"
	"time"
)

var delim = []byte("--boundarydonotcross\r\n")
var frame_start = []byte("\xff\xd8\xff\xe0")

func read_til_boundry(r *bufio.Reader) (data []byte) {
	for {
		buffer, err := r.ReadString('\n')

		if err != nil {
			panic("ERRROR")
		}

		data = append(data, []byte(buffer)...)

		if bytes.HasSuffix(data, delim) {
			return data[:len(data)-len(delim)]
		}
	}
}

func parse_frame(data []byte) (header []byte, body []byte) {
	for i := 0; i < len(data); i++ {
		if bytes.Compare(data[i:i+4], frame_start) == 0 {
			return data[:i], data[i:]
		}
	}
	panic("invalid frame")
}

func start_socket(port string, channel chan []byte) {
	conn, err := net.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		panic("Socket error")
	}
	conn.Write([]byte("GET /?action=stream HTTP/1.1\r\nHost:%s\r\n\r\n"))

	reader := bufio.NewReader(conn)

	//header
	_ = read_til_boundry(reader)

	for {
		var data = read_til_boundry(reader)
		var _, frame = parse_frame(data)
		//os.WriteFile("./out_" +port+ ".jpg", frame, 0644)
		channel <- frame
	}

}

func sink(channels []chan []byte) {
	for {
		for _, channel := range channels {
			frame := <-channel
			os.WriteFile("./out_.jpg", frame, 0644)
		time.Sleep(250 * time.Millisecond)
		}
	}
}

func main() {
	var args = os.Args[1:]

	var channels []chan []byte
	//for i := range args {
	//  channels = append(channels, make(chan []byte))

	var wg sync.WaitGroup
	for _, port := range args {
		var channel = make(chan []byte)
		//channels[i] = channel
		channels = append(channels, channel)
		wg.Add(1)
		go start_socket(port, channel)
	}

	go sink(channels)

	wg.Wait()
}

//https://go101.org/article/channel-use-cases.html
package main

import (
	"bufio"
	"bytes"
	"net"
	"os"
	"sync"
	"time"
)

type Frame struct {
	header []byte
	body   []byte
}

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

func parse_frame(data []byte) (frame Frame) {
	for i := 0; i < len(data); i++ {
		if bytes.Compare(data[i:i+4], frame_start) == 0 {
			return Frame{data[:i], data[i:]}
		}
	}
	panic("invalid frame")
}

func start_socket(port string) (reader *bufio.Reader) {
	conn, err := net.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		panic("Socket error")
	}
	conn.Write([]byte("GET /?action=stream HTTP/1.1\r\nHost:%s\r\n\r\n"))

	return bufio.NewReader(conn)

}

func source(reader *bufio.Reader) chan Frame {
	var source = make(chan Frame)
	go func() {
		//header
		_ = read_til_boundry(reader)
		for {
			var data = read_til_boundry(reader)
			var frame = parse_frame(data)
			//os.WriteFile("./out_" +port+ ".jpg", frame, 0644)
			source <- frame
		}
	}()

	return source
}

func sink(sources []chan Frame) {
	go func() {
		for {
			for _, source := range sources {
				var frame = <-source
				os.WriteFile("./out_.jpg", frame.body, 0644)
				time.Sleep(250 * time.Millisecond)
			}
		}
	}()
}

func main() {
	var args = os.Args[1:]
	var wg sync.WaitGroup

	var channels []chan Frame

	for _, port := range args {
		wg.Add(1)
		var reader = start_socket(port)
		var channel = source(reader)
		channels = append(channels, channel)
	}

	wg.Add(1)
	sink(channels)

	wg.Wait()
}

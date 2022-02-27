//https://go101.org/article/channel-use-cases.html
package main

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
)

type Frame struct {
	header []byte
	body   []byte
}

var delim = []byte("--boundarydonotcross\r\n")

// as per https://docs.fileformat.com/image/jpeg/
var frame_start = []byte("\xff\xd8")

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
		if bytes.Compare(data[i:i+2], frame_start) == 0 {
			return Frame{data[:i], data[i:]}
		}
	}

	println(string(data[77]))
	println(string(data[78]))
	println(bytes.Compare(data[77:77+2], frame_start) == 0)
	println(string(frame_start[0]))
	println(string(frame_start[1]))

	os.WriteFile("./parse_frame_panic_data", data, 0644)
	panic("invalid frame")
}

func receive_frame(connection net.Conn) (frame Frame) {
	//todo optimize

	// Read header
	var buffer = make([]byte, 1)
	for {

		var buffer_tmp = make([]byte, 1)
		var _, err = connection.Read(buffer_tmp[:])
		if err != nil {
			println("Can't read from connection")
			panic(err)
		}

		buffer = append(buffer, buffer_tmp...)

		// Check if jpeg start has been reached
		if len(buffer) > 2 && bytes.Compare(buffer[len(buffer)-2:], frame_start) == 0 {
			break
		}
	}

	// find index of content-length number
	var index = 0
	var wordIndex = 0
	var word = "Content-Length: "

	for {
		if wordIndex == len(word) {
			break
		}

		if buffer[index] == word[wordIndex] {
			wordIndex++
		} else {
			wordIndex = 0
		}

		index++
	}

	if wordIndex != len(word) {
		panic("Cannot find field Content-Length:")
	}

	// parse content size number
	var length = 0
	for {
		var el = buffer[index+length]
		if el < '0' || el > '9' {
			break
		}
		length++
	}
	var content_length, _ = strconv.Atoi(string(buffer[index : index+length]))

	// read rest
	var buffer_body = make([]byte, content_length-2)
	var n, err = io.ReadFull(connection, buffer_body)

	if err != nil {
		println("Can't read from connection")
		panic(err)
	}
	if n != content_length-2 {
		println(n)
		println(content_length)
		panic("Cannot read all bytes")
	}

	var body = append(frame_start, buffer_body...)
	//println(content_length)
	return Frame{buffer[:len(buffer)-2], body}
}

func start_socket(port string) (connection net.Conn) {
	conn, err := net.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		panic("Socket error")
	}
	conn.Write([]byte("GET /?action=stream HTTP/1.1\r\nHost:%s\r\n\r\n"))

	return conn

}

func source(connection net.Conn) chan Frame {
	var source = make(chan Frame)
	go func() {
		for {
			var frame = receive_frame(connection)
			source <- frame
			//select {
			// case source  <- frame:
			//default:

			//}
		}

		//header
		//_ = read_til_boundry(reader)
		//for {
		//	var data = read_til_boundry(reader)
		//	var frame = parse_frame(data)
		//	//os.WriteFile("./out_" +port+ ".jpg", frame, 0644)
		//	source <- frame

	}()

	return source
}

func sink(sources []chan Frame) {
	agg := make(chan Frame)
	for _, ch := range sources {
		go func(c chan Frame) {
			for msg := range c {
				agg <- msg
			}
		}(ch)
	}

	go func(agg chan Frame) {
		for {
			frame := <-agg
			os.WriteFile("./out_.jpg", frame.body, 0644)

		}
	}(agg)
}


func run(args []string){
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

func main() {
	var args = os.Args[1:]
  run(args)
}

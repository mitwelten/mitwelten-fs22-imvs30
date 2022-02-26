package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
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

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		println("Socket error")
	}
	conn.Write([]byte("GET /?action=stream HTTP/1.1\r\nHost:%s\r\n\r\n"))

	//reply := make([]byte, 2000)
	//_, err = conn.Read(reply)
	//fmt.Println(string(reply))

	reader := bufio.NewReader(conn)

	//header
	_ = read_til_boundry(reader)

	//first frame
	var data = read_til_boundry(reader)
	var _, frame = parse_frame(data)
	fmt.Println(string(frame))
}

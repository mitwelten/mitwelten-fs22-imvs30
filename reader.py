#!/bin/bash
from typing import Optional
import socket
import urllib.request

sep : list = list(map(lambda x:  bytes(x, 'utf-8'), "--boundarydonotcross\r\n"))
sep_len : int = len(sep)


def next_boundry(socket: socket) -> list:
    bytes = []

    #for i in range (0,335):
    while True:
        data = s.recv(1)
        bytes.append(data)

        if len(bytes) >= len(sep) and bytes[-sep_len:] == sep:
            break


    # return all but seperator
    #print(bytes)
    return bytes[0:-sep_len]


def remove_header(data: list) -> list:
    for i in range(len(data)):
        if data[i:i+4] == [b'\xff', b'\xd8', b'\xff', b'\xe0']:
            return data[i:]
    print("error")
    exit(1)



if __name__ == "__main__":
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    #s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    #s._fileobject.default_bufsize = 0
    s.connect(("127.0.0.1", 8080))
    request = "GET /?action=stream HTTP/1.1\r\nHost:%s\r\n\r\n"
    s.send(request.encode())

    # HTTP Header
    next_boundry(s)

    # first frame
    frame = next_boundry(s)
    frame = remove_header(frame)
    text = b''.join(frame)
    f = open('out1.jpg', 'wb')
    f.write(text)

    # next frame
    frame = next_boundry(s)
    frame = remove_header(frame)
    text = b''.join(frame)
    f = open('out2.jpg', 'wb')
    f.write(text)



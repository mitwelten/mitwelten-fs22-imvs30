#!/bin/bash
from typing import Optional
import socket
import urllib.request
from dataclasses import dataclass

sep: list = list(map(lambda x: bytes(x, "utf-8"), "--boundarydonotcross\r\n"))
sep_len: int = len(sep)
frame_start: list = [b"\xff", b"\xd8", b"\xff", b"\xe0"]

@dataclass
class Frame:
    header: list
    body: list

    def body_to_bytes(self) -> bytes:
        return b"".join(self.body)

    def write_to_file(self, name: str) -> None:
        f = open(name, "wb")
        f.write(self.body_to_bytes())
        f.close()


def read_til_boundry(socket: socket.socket) -> list:
    bytes = []

    # for i in range (0,335):
    while True:
        data = socket.recv(1)
        bytes.append(data)

        if len(bytes) >= len(sep) and bytes[-sep_len:] == sep:
            break

    # return all but seperator
    # print(bytes)
    return bytes[0:-sep_len]


def parse_frame(data: list) -> Frame:
    for i in range(len(data)):
        if data[i : i + 4] == frame_start:
            return Frame(data[:i], data[i:])

    raise Exception("Cannot parse frame")



def next_frame(socket: socket.socket) -> Frame:
    data = read_til_boundry(socket)
    return parse_frame(data)

if __name__ == "__main__":
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    # s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    # s._fileobject.default_bufsize = 0
    s.connect(("127.0.0.1", 8080))
    request = "GET /?action=stream HTTP/1.1\r\nHost:%s\r\n\r\n"
    s.send(request.encode())

    # HTTP Header
    read_til_boundry(s)

    # first frame
    frame = next_frame(s)
    frame.write_to_file("out1.jpg")

    # second frame
    frame = next_frame(s)
    frame.write_to_file("out2.jpg")


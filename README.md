# Mitwelten FS22 IMVS30
mjpg-multiplexer is a command line tool that takes one or more MJPEG streams as input and returns an MJPEG stream as output, which combines the input streams.

## Building & Installation
tbd: directly from github release..
- Clone repo:
    ```
    $ git clone git@github.com:mitwelten/mitwelten-fs22-imvs30.git
    ```
## Usage

The script can be parameterised and has two main purposes:
Redirect output to a file or provide output as a stream. 

- Redirect output to a file:
    ```
    $ go run ./src/main.go -input "192.168.137.216:8080 192.168.137.59:8080" -output "file" -output_filename "out.jpg" -method "grid"
    ```

- Output stream:
    ```
    $ go run ./src/main.go -input "192.168.137.102:8080 192.168.137.187:8080" -output "stream" -output_port "8088" -method "grid"
    ```

## Badge
[![Go](https://github.com/mitwelten/mitwelten-fs22-imvs30/actions/workflows/go.yml/badge.svg)](https://github.com/mitwelten/mitwelten-fs22-imvs30/actions/workflows/go.yml)

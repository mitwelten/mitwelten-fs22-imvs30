# Mitwelten FS22 IMVS30
mjpg-multiplexer is a command line tool that takes one or more MJPEG streams as input and returns an MJPEG stream as output, which combines the input streams.


## Prerequisites

- libjpeg:
  libjpeg is a free library with functions for handling the JPEG image data format.
    ```
    $ sudo apt-get install -y libjpeg-dev
    ```

## Building & Installation
tbd: directly from github release..
- Clone repo:
    ```
    $ git clone git@github.com:mitwelten/mitwelten-fs22-imvs30.git
    ```
## Usage

This script can be parameterised and has different modes. First argument determines the mode.

    
    Usage:
    multiplexer motion input <URLs> output <PORT> [options]
    multiplexer grid --grid_dimension <GRID_ROWS> <GRID_COLUMNS> input <URLs> output <PORT> [options]
    multiplexer carousel --duration <CAROUSEL_DURATION> input <URLs> output <PORT> [options]
    multiplexer panel --cycle --duration <CYCLE_DURATION> input <URLs> output <PORT> [options]
    multiplexer -h | --help
    multiplexer --version
    
    Options:
    -h --help               Shows this screen.
    --input_framerate=n     input framerate in fps [default: -1]
    --output_framerate=n    output framerate in fps[default: -1]
    --output_max_width=n    output width in pixel [default: -1]
    --output_max_height=n   output height in pixel [default: -1]  
    --use_auth              Use Authentication [default: false]
    --log_time              Log Time verbose [default: false]
    --verbose               Shows details. [default: false]
    --version               Shows version.`

## Examples

- Output stream (grid):
    ```
    $ go run ./src/main.go grid --dimension 2 1 input localhost:8080,localhost:8080 output 8088 --log_time --use_auth --input_framerate 10 --output_framerate 10 
    ```

- Output stream (motion):
    ```
    $ go run ./src/main.go motion input localhost:8080,localhost:8080 output 8088 --log_time --use_auth --input_framerate 10 --output_framerate 10  
    ```

- Output stream (carousel):
    ```
    $ go run ./src/main.go carousel --duration 5 input localhost:8080,localhost:8080 output 8088 --log_time --use_auth --input_framerate 10 --output_framerate 10  
    ```

- Output stream (panel):
    ```
    $ go run ./src/main.go panel input localhost:8080,localhost:8080 output 8088 --log_time --use_auth --input_framerate 10 --output_framerate 10  
    ```

- Output stream (panel-cycling):
    ```
    $ go run ./src/main.go panel --cycle input localhost:8080,localhost:8080 output 8088 --log_time --use_auth --input_framerate 10 --output_framerate 10  
    ```

- Output stream (panel-cycling with custom duration [in seconds]):
    ```
    $ go run ./src/main.go panel --cycle --duration 10 input localhost:8080,localhost:8080 output 8088 --log_time --use_auth --input_framerate 10 --output_framerate 10  
    ```

## Build Automation
We use [Go release Action](https://github.com/wangyoucao577/go-release-action) to automatically publish Go binaries to Github Release Assets through Github Action. 

### Build Badge
[![Go](https://github.com/mitwelten/mitwelten-fs22-imvs30/actions/workflows/go.yml/badge.svg)](https://github.com/mitwelten/mitwelten-fs22-imvs30/actions/workflows/go.yml)

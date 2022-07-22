[![Run Tests](https://github.com/mitwelten/mitwelten-fs22-imvs30/actions/workflows/tests.yml/badge.svg)](https://github.com/mitwelten/mitwelten-fs22-imvs30/actions/workflows/tests.yml)

# Mitwelten FS22 IMVS30
mjpg-multiplexer is a command line tool that takes one or more MJPEG streams as input and returns an MJPEG stream as output, which combines the input streams.


## Prerequisites

- libjpeg-turbo:
  [libjpeg-turbo](https://libjpeg-turbo.org/) is a free library with functions for efficient handling of the JPEG image data format.

    
    Installation on Ubuntu:
    ```
    $ sudo apt-get install -y libjpeg-turbo8
    ```

## Building & Installation

- Prerequisites: go 1.18+

- Clone repo
- Change into folder
- Build with the go toolchain using `go build`

    ```
    $ git clone git@github.com:mitwelten/mitwelten-fs22-imvs30.git
    $ cd mitwelten-fs22-imvs30
    $ go build -o "multiplexer" ./src/main.go
    ```
  
## Usage
This script can be parameterised and has different modes. First argument determines the mode.

```
 $ ./multiplexer [grid | panel | carousel] input [URL] output [URL] [options...]
                 <--------- mode --------> <- input -> <- output ->
```
    
```
    $ ./multiplexer --help
    
    Usage: multiplexer [grid | panel | carousel] input [URL] output [URL] [options...]
                       <--------- mode --------> <- input -> <- output ->
    Mode:
      grid: static grid of images with X rows and Y columns
      panel: dynamic panel of.... Can be used with motion (see --motion)
      carousel: dynamic carousel view.... Can be used with motion (see --motion)
    Input:  comma separated list of urls including port
    Output: output url including port
    
    Examples: 
      ./multiplexer grid input localhost:8080,localhost:8081 output 8088
      ./multiplexer panel input :8080,:8081,:8082 output 8088 --cycle --width 800 
      ./multiplexer carousel input 192.168.0.1:8080 192.168.0.2:8081 output 8088 --motion
    
    Options:
      --grid_dimension=ROWS,COLUMNS    Number of cells in the grid mode
      --motion                         Enables motion detection to focus the most active frame on selected mode
      --cycle                          Enables cycling of the panel layout, see also [--duration] 
      --duration=n                     Duration in seconds before changing the layout (panel and carousel only) [default: 15]
      --width=n                        Output width in pixel [default: -1]
      --height=n                       Output height in pixel [default: -1]
      --ignore_aspect_ratio            Stretches the frames instead of adding a letterbox on resize
      --framerate=n                    Output framerate in fps[default: -1]
      --quality=n                      Output jpeg quality in percentage [default: -1]
      --use_auth                       Use Authentication
      --show_border                    Enables a border in the grid and panel layout between the images
      --show_label                     Show label for input streams
      --labels=n                       Comma separated list of names to show instead of the camera input url
      --label_font_size=n              Input label font size in pixel [default: 32]
      --log_fps                        Logs the current FPS 
      -v --version                     Shows version.
      -h --help                        Shows this screen.`
```

## Examples

TODO update

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

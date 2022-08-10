[![Run Tests](https://github.com/mitwelten/mitwelten-fs22-imvs30/actions/workflows/tests.yml/badge.svg)](https://github.com/mitwelten/mitwelten-fs22-imvs30/actions/workflows/tests.yml)

# Mitwelten FS22 IMVS30

mjpg-multiplexer is a command line tool that takes one or more MJPEG streams as input and returns an
MJPEG stream as output, which combines the input streams.

## Prerequisites

- libjpeg-turbo:
  [libjpeg-turbo](https://libjpeg-turbo.org/) is a free library with functions for efficient
  handling of the JPEG image data format.

    Installation on Ubuntu:
    ```
    $ sudo apt-get install -y libjpeg-turbo8
    ```

## Building & Installation

- Prerequisites: go 1.19+

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
 $ ./mjpeg_multiplexer [grid | panel | carousel] input [URL] output [URL] [options...]
                 <--------- mode --------> <- input -> <- output ->
```

```
    $ ./mjpeg_multiplexer --help
    
  Usage: multiplexer [grid | panel | carousel] input [URL] output [URL] [options...]
                   <--------- mode --------> <- input -> <- output ->

The multiplexer combines multiple multiple input streams to an output stream using a mode.

Mode:
  grid: static grid of images with X rows and Y columns
  panel: dynamic panel of.... Can be used with activity detection (see --activity)
  carousel: dynamic carousel view.... Can be used with activity detection (see --activity)
Input:  comma separated list of urls including port
Output: output url including port

Examples: 
  $ ./mjpeg_multiplexer grid input localhost:8080,localhost:8081 output 8088
  $ ./mjpeg_multiplexer panel input :8080,:8081,:8082 output 8088 --panel_cycle --width 800 
  $ ./mjpeg_multiplexer carousel input 192.168.0.1:8080 192.168.0.2:8081 output 8088 --activity

Options:
  --grid_dimension [list]          Comma separated list of the number of cells in the grid mode, eg. '--grid_dimension "3,2"'
  --activity                       Enables activity detection to focus the most active frame on selected mode
  --panel_cycle                    Enables cycling of the panel layout, see also [--duration] 
  --duration [number]              Duration in seconds before changing the layout (panel and carousel only) [default: 15]
  --width [number]                 Total output width in pixel
  --height [number of number]      Total output height in pixel
  => if only the height or width is specified, the other will be adjusted with regards to the ascpect ratio
  --ignore_aspect_ratio            Stretches the frames instead of adding a letterbox on resize
  --framerate [number]             Limit the output framerate per second
  --quality [number]               Output jpeg quality in percentage [default: 80]
  --use_auth                       Use Authentication
  --show_border                    Enables a border in the grid and panel layout between the images
  --show_label                     Show label for input streams
  --labels [list]                  Comma separated list of alternative label text, eg. '--labels "label 1, label 2"'
  --label_font_size [number]       Input label font size in pixel [default: 32]
  --log_fps                        Logs the current FPS 
  -v --version                     Shows version.
  -h --help                        Shows this screen

Authentication to connect to mjpeg_streamer streams secured by credentials can be enabled using the [--use_auth] flag. Add the credentials to the 'authentication.json' file. See 'authentication_example.json' as an example.`
```

## Examples

- Grid 
    ```
    $ ./mjpeg_multiplexer input :8080,:8081,:8082,:8083 output 8088 --grid_dimension 2,2 --width 1280
    ```

- Panel (with `--activity` and custom labels)
    ```
    $ ./mjpeg_multiplexer panel input :8080,:8081,:8082,:8083,:8084 output 8088 --quality 100 --activity --log_fps --show_label --labels "Wild-Cam 1,Wild-Cam 2,Sea-Cam 1,Sea-Cam 2"
    ```

- Carousel (with passthrough for high fps)
    ```
    $ ./mjpeg_multiplexer carousel input :8080,:8081,:8082,:8083,:8084 output 8088"
    ```

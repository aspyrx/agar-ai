package main

// #cgo LDFLAGS: -lX11
// #include <X11/Xlib.h>
//
// XImage *ScreenImage(Display *dpy, Window window, int x, int y, unsigned int width, unsigned int height) {
//     return XGetImage(dpy, window, x, y, width, height, AllPlanes, ZPixmap);
// }
//
// void FreeImage(XImage *image) {
//     XDestroyImage(image);
// }
import "C"

import (
	"fmt"
	"math"
	"net"
)

var dpy *C.Display
var root_window C.Window
var image *C.XImage
var width int
var height int
var capWidth int
var capHeight int
const border int = 50

func move(x, y int) {
	C.XWarpPointer(dpy, C.None, root_window, 0, 0, 0, 0, C.int(x), C.int(y))
	C.XFlush(dpy)
}

func handle(conn net.Conn) {
	for {
		// Reset global variables.
		image = C.ScreenImage(dpy, root_window, C.int(border), C.int(border), C.uint(capWidth), C.uint(capHeight))
		checked = make([][]bool, capWidth)
		checkedPx := make([]bool, capWidth * capHeight)
		for i := range checked {
		    checked[i], checkedPx = checkedPx[:capHeight], checkedPx[capHeight:]
		}
		cluster = make([]cell, 1)
		cluster[0].dist = 1<<63 - 1
		clusterI = 0
		process(conn)
		C.FreeImage(image)
        ourSize := 0
        for _, value := range cluster {
            if value.dist == 0 {
                ourSize = value.size
                break
            }
        }
		
		var fieldX, fieldY float64
		for _, value := range cluster {
		    if value.dist != 0 {
		        charge := value.size - ourSize
		        angle := math.Atan2(float64(capHeight / 2 - value.pY), float64(capWidth / 2 - value.pX))
		        mag := float64(charge) / float64(value.dist)
		        fieldX += mag * math.Cos(angle)
		        fieldY += mag * math.Sin(angle)
		    }
		}
		
		angle := math.Atan2(float64(fieldY), float64(fieldX))
		resultantX := int(float64(capHeight) / 2.0 * math.Cos(angle))
		resultantY :=  int(float64(capHeight) / 2.0 * math.Sin(angle))
		move(width / 2 + resultantX, height / 2 + resultantY)
		fmt.Fprint(conn, angle, fieldX, fieldY, "\r\n")
	}
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		handle(conn)
	}
}

func init() {
	zero := C.char(0)
	dpy = C.XOpenDisplay(&zero)
	root_window = C.XRootWindow(dpy, 0)
	width = int(C.XDisplayWidth(dpy, 0))
	height = int(C.XDisplayHeight(dpy, 0))
	capWidth = width - border * 2
	capHeight = height - border * 2
	fmt.Println(width, height, capWidth, capHeight)
	C.XSelectInput(dpy, root_window, C.KeyReleaseMask)
}

package main

// #cgo LDFLAGS: -lX11
// #include <X11/Xlib.h>
//
// int get_shift (int mask) {
//     int shift = 0;
//     while (mask) {
//         if (mask & 1) break;
//             shift++;
//             mask >>=1;
//       }
//       return shift;
//     }
//
// unsigned long RedValue(XImage *image, int x, int y) {
//     return (XGetPixel(image, x, y) & image->red_mask) >> get_shift(image->red_mask);
// }
import "C"

import (
	"net"
)

var checked [][]bool

type cell struct {
	size, dist, pX, pY int
}

var cluster []cell

// Keep track of the number of clusters.
var clusterI int

func floodFill(x, y, depth int) {
	// If target-color is equal to replacement-color, return.
	if checked[x][y] {
		return
	}

	// If the color of node is not equal to target-color, return.
	r := uint64(C.RedValue(image, C.int(x), C.int(y)))
	if r < 64 {
		// Value is black. Return.
		return
	}

	// Set the pixel as "checked".
	checked[x][y] = true

	// Increase number of pixels in block.
	cluster[clusterI].size++

	// Determine square of distance to center.
	dist := (capWidth / 2 - x) * (capWidth / 2 - x) + (capHeight / 2 - y) * (capHeight / 2 - y)
	if cluster[clusterI].dist > dist {
		cluster[clusterI].dist = dist
		cluster[clusterI].pX = x
		cluster[clusterI].pY = y
	}

	if x + 1 < capWidth {
		floodFill(x+1, y, depth+1)
	}
	if x - 1 > -1 {
		floodFill(x-1, y, depth+1)
	}
	if y + 1 < capHeight {
		floodFill(x, y+1, depth+1)
	}
	if y - 1 > -1 {
		floodFill(x, y-1, depth+1)
	}
	if depth == 0 {
		// We have finished processing a cluster. Increment the cluster count.
		clusterI++

		// Initiate the next cell in the cluster.
		cluster = append(cluster, cell{dist: 1<<63 - 1})
	}
}

func process(conn net.Conn) {
	for a := 0; a < capWidth; a += 4 {
		for b := 0; b < capHeight; b += 4 {
			floodFill(a, b, 0)
		}
	}

	// Remove last empty cell.
	cluster = cluster[:len(cluster)-1]
}

func init() {
	// Initiate first cell.
	cluster = make([]cell, 1)
	cluster[0].dist = 1<<63 - 1
}

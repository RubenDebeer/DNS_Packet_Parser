package main

import (
	"errors"
)

func main() {

}

var errEOB error = errors.New(" End Of Buffer")

type packetBuffer struct {
	ByteBuffer     [512]byte
	ReaderPosition int
}

// Initialize the Buffer

// Load a copy of the buffer

// Get the byte based on a Position

// Step over a byte

// Change the Position in a Buffer

// Read a Single Byte

// get a Single byte without Moving the buffer position

// Get a range of Bytes

// Read 2 Bytes

// Read 4  bytes and stepping four steps forward

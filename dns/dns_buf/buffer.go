package dnsbuf

// The goal of this Module is to provide helper functions to read a Packet

import (
	"errors"
)

var errByteOutOfRange error = errors.New("ByteOut of range.")

type ByteBuffer struct {
	buffer         [512]byte
	readerPosition int
}

// Why Update the Buffer in Place? So that we  don't have to pass a copy of the buffer every time we do a operation.
func New() *ByteBuffer {
	return &ByteBuffer{
		buffer:         [512]byte{},
		readerPosition: 0,
	}
}

// Get the byte based on a Position
func (b *ByteBuffer) Pos() int { return b.readerPosition }

// Step over the Buffer a specific number of steps.
func (b *ByteBuffer) Step(steps int) {
	b.readerPosition += steps
}

// Change the Position in a Buffer
func (b *ByteBuffer) ChangePosition(position int) {
	b.readerPosition = position
}

// Read a Single Byte
func (b *ByteBuffer) ReadByte() (returnByte byte, ReturnError error) { // I added the Explicit return because the Naked return makes me feel Dirty,(Like my code is dirty).
	if b.readerPosition > len(b.buffer) {
		returnByte = 0
		ReturnError = errByteOutOfRange
		return returnByte, ReturnError
	}
	returnByte = b.buffer[b.readerPosition]
	ReturnError = nil

	b.readerPosition++

	return returnByte, ReturnError
}

// get a Single byte without Moving the buffer position
func (b *ByteBuffer) GetByte() (returnByte byte, ReturnError error)

// Get a range of Bytes

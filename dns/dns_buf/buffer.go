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
// I added the Explicit return because the Naked return makes me feel Dirty,(Like my code is dirty).
func (b *ByteBuffer) ReadByte() (returnByte byte, returnError error) {
	if b.readerPosition > len(b.buffer) {
		returnByte = 0
		returnError = errByteOutOfRange
		return returnByte, returnError
	}

	returnByte = b.buffer[b.readerPosition]
	returnError = nil

	b.readerPosition++

	return returnByte, returnError
}

// get a Single byte without Moving the buffer position
func (b *ByteBuffer) GetByte() (returnByte byte, returnError error) {

	if b.readerPosition > len(b.buffer) {
		returnByte = 0
		returnError = errByteOutOfRange
		return returnByte, returnError
	}
	returnByte = b.buffer[b.readerPosition]
	returnError = nil

	return returnByte, returnError
}

// Get a range of Bytes
func (b *ByteBuffer) GetByteRange(startPos, copyLength int) (byteRange []byte, byteRangeError error) {

	if startPos > len(b.buffer) || startPos < 0 || startPos+copyLength > len(b.buffer) {
		byteRange = nil
		byteRangeError = errByteOutOfRange
		return byteRange, byteRangeError
	}

	byteRange = b.buffer[:]
	byteRangeError = nil

	return byteRange, byteRangeError
}

// Read 2 bytes as Big-Endian
func (b *ByteBuffer) ReadU16() (TwoBytes uint16, errRead error) {
	b1, err := b.ReadByte()

	if err != nil {
		return 0, err
	}

	b2, err := b.ReadByte()
	if err != nil {
		return 0, err
	}

	TwoBytes = (uint16(b1) << 8) | uint16(b2)
	errRead = nil

	return TwoBytes, errRead
}

// Read 4 bytes as Big-Endian
func (b *ByteBuffer) ReadU32() (FourBytes uint32, errRead error) {
	b1, err := b.ReadByte()
	if err != nil {
		return 0, err
	}

	b2, err := b.ReadByte()
	if err != nil {
		return 0, err
	}

	b3, err := b.ReadByte()
	if err != nil {
		return 0, err
	}

	b4, err := b.ReadByte()
	if err != nil {
		return 0, err
	}

	//               32 - 8 == 24          14 - 8 == 16         16 - 8 == 8
	FourBytes = (uint32(b1) << 24) | (uint32(b2) << 16) | (uint32(b3) << 8) | uint32(b4)
	errRead = nil

	return FourBytes, errRead
}

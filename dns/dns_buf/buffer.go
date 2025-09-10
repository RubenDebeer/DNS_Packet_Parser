package dnsbuf

// The goal of this Module is to provide helper functions to read a Packet

import (
	"errors"
	"fmt"
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
func (b *ByteBuffer) ReadByte() (returnByte byte) {
	if b.readerPosition > len(b.buffer) {
		returnByte = 0
		return returnByte
	}

	returnByte = b.buffer[b.readerPosition]

	b.readerPosition++

	return returnByte
}

// Get a Single byte without Moving the buffer position
func (b *ByteBuffer) GetByte(position int) (returnByte byte) {

	if position > len(b.buffer) {
		returnByte = 0
		return returnByte
	}

	returnByte = b.buffer[position]
	return returnByte
}

// Get a range of Bytes
func (b *ByteBuffer) GetByteRange(startPos, copyLength int) (byteRange []byte) {

	if startPos > len(b.buffer) || startPos < 0 || startPos+copyLength > len(b.buffer) {
		byteRange = nil
		return byteRange
	}

	byteRange = b.buffer[startPos : startPos+copyLength]

	return byteRange
}

// Read 2 bytes as Big-Endian
func (b *ByteBuffer) ReadU16() (TwoBytes uint16) {
	b1 := b.ReadByte()

	b2 := b.ReadByte()

	TwoBytes = (uint16(b1) << 8) | uint16(b2)

	return TwoBytes
}

// Read 4 bytes as Big-Endian
func (b *ByteBuffer) ReadU32() (FourBytes uint32, errRead error) {
	b1 := b.ReadByte()

	b2 := b.ReadByte()

	b3 := b.ReadByte()

	b4 := b.ReadByte()

	//               32 - 8 == 24          14 - 8 == 16         16 - 8 == 8
	FourBytes = (uint32(b1) << 24) | (uint32(b2) << 16) | (uint32(b3) << 8) | uint32(b4)
	errRead = nil

	return FourBytes, errRead
}

func (b *ByteBuffer) ReadQName(domainName string) (errWrite error) {

	// --> [6][google][3][com][0]
	// --> [0x06][0x67 0x6F 0x6F 0x67 0x6C 0x65][0x03][0x63 0x6F 0x6D][0x00]
	//[len][bytes][len][bytes]...[0]

	// --> [0x06][0x67 0x6F 0x6F 0x67 0x6C 0x65] | [0x03][0x63 0x6F 0x6D] | [0x00]

	// Create a local position
	// Read the Original buffer from current position for the length provided in the message
	// Create a While loop
	// Get the length
	// Read the content from the original buffer

	// Read the first byte

	var localPos int = b.Pos()
	var doneReading bool = false

	for doneReading != false {

		lengthByte := b.GetByte(localPos)

		if lengthByte == 0 {
			break
		}

		label := b.GetByteRange(localPos, int(lengthByte))

		if label != nil {
			domainName += string(label) + "."
		}

		localPos++

		fmt.Println(label)

	}

}

package dnsbuf

// The goal of this Module is to provide helper functions to read a Packet

import (
	"fmt"
	"strings"
)

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

func (b *ByteBuffer) ReadQname() (string, error) {
	const maxJumps = 10

	localPos := 0
	jumped := false
	jumps := 0

	var out strings.Builder
	delimiter := ""

	for {
		// Vaidation
		if jumps > maxJumps {
			return " ", fmt.Errorf("who , we are trying to jump the gun there max is %d", maxJumps)
		}
		if localPos >= len(b.buffer) {
			return "", fmt.Errorf("whow, out of range")
		}

		length := b.Pos()

		// If the byte at Pos is 11 i.e a flag that indicates a comprtession we need to move to the actual value

		// Compressed Lable Case
		if (length & 0xC0) == 0xC0 {
			if localPos >= len(b.buffer) {
				return "", fmt.Errorf("whow, out of range")
			}

			if !jumped {
				b.ChangePosition(localPos + 1)
			}

			b2 := b.GetByte(localPos + 1)
			offset := int((uint16(length&0x3F) << 8) | uint16(b2))

			localPos = offset
			jumped = true
			jumps++
			continue
		}

		// Normal Lable case
		localPos++
		if length == 0 {
			break
		}

		if localPos+(length) > len(b.buffer) {
			return "", fmt.Errorf("whow, out of range")
		}

		out.WriteString(delimiter)
		lable := b.GetByteRange(localPos, int(length))
		out.WriteString(strings.ToLower(string(lable)))

		delimiter = "."
		localPos = int(length)
	}

	if !jumped {
		b.ChangePosition(localPos)
	}

	return out.String(), nil
}

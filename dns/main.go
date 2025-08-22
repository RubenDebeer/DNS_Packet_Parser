package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}

// So thinking we need some Data Structure to store the Packets
type Packet struct {
	PacketBuffer [512]byte
	position     int
}

// we want to initialize the packet to zero .
func new_buf() *Packet {
	return &Packet{
		PacketBuffer: [512]byte{},
		position:     0,
	}
}

// This new for me but this is called a "Receiver" and it's like a method of a class.
func (p *Packet) Load(internalCopy []byte) {

	packetByteLength := len(internalCopy)

	if packetByteLength > 512 {
		packetByteLength = len(p.PacketBuffer)
	}
	copy(p.PacketBuffer[:], internalCopy[:packetByteLength])
}

// A receiver to get the pozzy
func (p *Packet) Pos() int {
	return p.position
}

// A Receiver to step over the Buffer
func (p *Packet) Step(steps int) error {
	p.position += steps
	return nil
}

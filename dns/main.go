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

// This is new for me but I need to define a "receiver" it's like a method within a "Class" I Think...

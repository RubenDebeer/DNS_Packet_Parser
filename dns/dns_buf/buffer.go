package dnsbuf

// The goal of this Module is to provide helper functions to read a Packet
import (
	"fmt"
	"net"
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

func (b *ByteBuffer) Load(p []byte) int {
	n := copy(b.buffer[:], p)
	b.readerPosition = 0
	return n
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
func (b *ByteBuffer) ReadInt16() (TwoBytes uint16) {
	b1 := b.ReadByte()

	b2 := b.ReadByte()

	TwoBytes = (uint16(b1) << 8) | uint16(b2)

	return TwoBytes
}

// Read 4 bytes as Big-Endian
func (b *ByteBuffer) ReadInt32() (FourBytes uint32, errRead error) {
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

// Result Code
type ResultCode int

const (
	NOERROR  ResultCode = 0
	FORMERR  ResultCode = 1
	SERVFAIL ResultCode = 2
	NXDOMAIN ResultCode = 3
	NOTIMP   ResultCode = 4
	REFUSED  ResultCode = 5
)

func (rescode ResultCode) String() string {
	switch rescode {
	case NOERROR:
		return "NOERROR"
	case FORMERR:
		return "FORMERR"
	case SERVFAIL:
		return "SERVFAIL"
	case NXDOMAIN:
		return "NXDOMAIN"
	case NOTIMP:
		return "NOTIMP"
	case REFUSED:
		return "REFUSED"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", int(rescode))
	}
}

// Given a Number retun a Result Code
func FromNum(num int) ResultCode {
	switch num {
	case 1:
		return FORMERR
	case 2:
		return SERVFAIL
	case 3:
		return NXDOMAIN
	case 4:
		return NOTIMP
	case 5:
		return REFUSED
	case 0:
		fallthrough
	default:
		return NOERROR
	}
}

// Header
type DnsHeader struct {
	id     int16
	qr     bool
	opcode int16

	flag_authoritative_awnser bool
	flag_truncated_message    bool
	flag_recursion_desired    bool
	flag_recursion_available  bool

	// Z Flag 3 Bits
	flag_reserved_z          bool
	flag_athenticated_data_z bool
	flag_checking_disabled_z bool

	// 4 Bits
	flag_response_code ResultCode

	question_count  int16
	awnser_count    int16
	authority_count int16
	aditional_count int16
}

func NewDnsHeader() DnsHeader {
	return DnsHeader{
		id:     0,
		qr:     false,
		opcode: 0,

		flag_authoritative_awnser: false,
		flag_truncated_message:    false,
		flag_recursion_desired:    false,
		flag_recursion_available:  false,
		flag_reserved_z:           false,
		flag_athenticated_data_z:  false,
		flag_checking_disabled_z:  false,
		flag_response_code:        NOERROR,

		question_count:  0,
		awnser_count:    0,
		authority_count: 0,
		aditional_count: 0,
	}
}

func (h *DnsHeader) Read(packetBuffer *ByteBuffer) error {
	h.id = int16(packetBuffer.ReadInt16())

	flags := packetBuffer.ReadInt16()
	upperByte := byte(flags >> 8)
	lowerByte := byte(flags & 0x00FF)
	// Uppper Lower --> Big Endian-Ness Most Significant Bit in the Least significant memmory address.

	//Bit   7  6   5  4  3  2  1  0
	//Field QR OP          AA TC RD
	h.flag_recursion_desired = (upperByte & (1 << 0)) != 0
	h.flag_truncated_message = (upperByte & (1 << 1)) != 0
	h.flag_authoritative_awnser = (upperByte & (1 << 2)) != 0
	h.opcode = int16((upperByte >> 3) & 0x0F)
	h.qr = (upperByte & (1 << 7)) != 0

	h.flag_response_code = FromNum(int(lowerByte & 0x0F))
	h.flag_checking_disabled_z = (lowerByte & (1 << 4)) != 0
	h.flag_athenticated_data_z = (lowerByte & (1 << 5)) != 0
	h.flag_reserved_z = (lowerByte & (1 << 6)) != 0
	h.flag_recursion_available = (lowerByte & (1 << 7)) != 0

	// Section counts
	h.question_count = int16(packetBuffer.ReadInt16())
	h.awnser_count = int16(packetBuffer.ReadInt16())
	h.authority_count = int16(packetBuffer.ReadInt16())
	h.aditional_count = int16(packetBuffer.ReadInt16())

	return nil
}

// Query Type
type QueryType int

const (
	UNKNOWN QueryType = 0
	A       QueryType = 1
	CNAME   QueryType = 2
)

// DNS Question
type DnsQuestion struct {
	QName  string
	QType  QueryType
	QClass uint16
}

func NewDnsQuestion(name string, q_type QueryType) DnsQuestion {
	return DnsQuestion{
		QName: name,
		QType: q_type,
	}
}

func (dns_query *DnsQuestion) Read(packetBuffer *ByteBuffer) error {

	dns_query.QName, _ = packetBuffer.ReadQname()
	dns_query.QType = QueryType(packetBuffer.ReadInt16())
	dns_query.QClass = packetBuffer.ReadInt16()

	return nil
}

// DNS Record
type DnsRecord interface {
	Type() QueryType
	Domain() string
	TTL() uint32
	fmt.Stringer
}

type UnknownRecord struct {
	DomainName string
	QType      uint16
	DataLen    uint16
	TTLVal     uint32
}

func (r *UnknownRecord) Type() QueryType {
	return UNKNOWN
}

func (r *UnknownRecord) Domain() string {
	return r.DomainName
}

func (r *UnknownRecord) TTL() uint32 {
	return r.TTLVal
}

func (r *UnknownRecord) String() string {
	return fmt.Sprintf("UNKNOWN{domain=%s qtype=%d data_len=%d ttl=%d}",
		r.DomainName, r.QType, r.DataLen, r.TTLVal)
}

// All of this Shit, just Because Enums are non existent in Go
type ARecord struct {
	DomainName string
	Addr       net.IP // IPv4 this is also a length of 4
	TTLVal     uint32
}

func (r *ARecord) Type() QueryType { return A }
func (r *ARecord) Domain() string  { return r.DomainName }
func (r *ARecord) TTL() uint32     { return r.TTLVal }
func (r *ARecord) String() string {
	return fmt.Sprintf("A{domain=%s addr=%s ttl=%d}", r.DomainName, r.Addr.String(), r.TTLVal)
}

func ReadRecord(packetBuffer *ByteBuffer) (DnsRecord, error) {

	name, err := packetBuffer.ReadQname()

	if err != nil {
		return nil, fmt.Errorf("reading the resource record went wonk: %w", err)
	}

	qType := packetBuffer.ReadInt16()
	_ = packetBuffer.ReadInt16()

	ttl, err := packetBuffer.ReadInt32()
	if err != nil {
		return nil, fmt.Errorf("read rr ttl: %w", err)
	}
	rdlen := packetBuffer.ReadInt16()

	switch QueryType(qType) {
	case A:
		if rdlen != 4 {
			_ = packetBuffer.GetByteRange(packetBuffer.Pos(), int(rdlen))
			packetBuffer.Step(int(rdlen))
			return &UnknownRecord{
				DomainName: name,
				QType:      qType,
				DataLen:    rdlen,
				TTLVal:     ttl,
			}, fmt.Errorf("a record with invalid rdlen=%d", rdlen)
		}

		// Read bytes to Build up  IPV4 address boy
		b0 := packetBuffer.ReadByte()
		b1 := packetBuffer.ReadByte()
		b2 := packetBuffer.ReadByte()
		b3 := packetBuffer.ReadByte()

		ip := net.IPv4(b0, b1, b2, b3).To4()
		return &ARecord{
			DomainName: name,
			Addr:       ip,
			TTLVal:     ttl,
		}, nil

	default:
		packetBuffer.Step(int(rdlen))
		return &UnknownRecord{
			DomainName: name,
			QType:      qType,
			DataLen:    rdlen,
			TTLVal:     ttl,
		}, nil
	}
}

type DnsPacket struct {
	Header      DnsHeader
	Questions   []DnsQuestion
	Answers     []DnsRecord
	Authorities []DnsRecord
	Resources   []DnsRecord
}

func NewDnsPacket() DnsPacket {
	return DnsPacket{
		Header:      NewDnsHeader(),
		Questions:   make([]DnsQuestion, 0),
		Answers:     make([]DnsRecord, 0),
		Authorities: make([]DnsRecord, 0),
		Resources:   make([]DnsRecord, 0),
	}
}

func ReadPacket(buffer *ByteBuffer) (DnsPacket, error) {

	packet := NewDnsPacket()

	if err := (&packet.Header).Read(buffer); err != nil {
		return DnsPacket{}, err
	}

	for i := 0; i < int(packet.Header.question_count); i++ {
		q := DnsQuestion{}
		if err := (&q).Read(buffer); err != nil {
			return DnsPacket{}, err
		}
		packet.Questions = append(packet.Questions, q)
	}

	for i := 0; i < int(packet.Header.awnser_count); i++ {
		rec, err := ReadRecord(buffer)
		if err != nil {
			return DnsPacket{}, err
		}
		packet.Answers = append(packet.Answers, rec)
	}

	for i := 0; i < int(packet.Header.authority_count); i++ {
		rec, err := ReadRecord(buffer)
		if err != nil {
			return DnsPacket{}, err
		}
		packet.Authorities = append(packet.Authorities, rec)
	}

	for i := 0; i < int(packet.Header.aditional_count); i++ {
		rec, err := ReadRecord(buffer)
		if err != nil {
			return DnsPacket{}, err
		}
		packet.Resources = append(packet.Resources, rec)
	}

	return packet, nil
}

// EBML specification:
//   https://matroska-org.github.io/libebml/specs.html
package ebml

// TODO: format long comment into line comment
// https://matroska.org/technical/specs/index.html

import (
	"math"
)

type Class int

const (
	ClassA Class = 1 + iota
	ClassB
	ClassC
	ClassD
)

func CodedInt(buf []byte) (Class, int64) {
	if buf[0] == 0 {
		return 0, 0
	}
	mask, value := int64(0x80), int64(buf[0])
	class := ClassA
	for i := 1; i < 8 && mask > value; i++ {
		value = (value << 8) + int64(buf[i])
		mask <<= 7
		class++
	}
	return class, value - mask
}

func PutCodedInt(buf []byte, n int64) Class {
	i := 0
	for m := n; m > ((1 << 7) - 1); m >>= 8 {
		i++
	}
	j := i + 1
	mask := byte(0x80)
	for i > 0 {
		buf[i] = byte(n & 0xFF)
		mask >>= 1
		n >>= 8
		i--
	}
	buf[0] = mask | byte(n)
	return Class(j)
}

type Reader struct {
	buf        []byte
	start, pos int
	miss       bool
}

func NewReader(buf []byte) *Reader {
	return &Reader{buf: buf}
}

func (p *Reader) Init(buf []byte) {
	*p = Reader{buf: buf}
}

// Bytes returnes the underlying buffer of element.
func (p *Reader) Bytes() []byte { return p.buf }

// Tell tells the position of the cursor calculated from the entire buffer.
func (p *Reader) Tell() int { return p.start + p.pos }

// AtEOS reports whether reaching the end or not.
func (p *Reader) AtEOS() bool { return p.pos >= len(p.buf) }

// Failed reports whether the recent operation had been failed due to lacked buffer or not.
func (p *Reader) Failed() bool { return p.miss }

func (p *Reader) ReadByte() byte {
	b := p.buf[p.pos]
	p.pos++
	return b
}

// ReadId reads Element ID.
func (p *Reader) ReadId() int32 {
	mask, value := int32(0x80), int32(p.ReadByte())
	for i := 0; i < 7 && !p.AtEOS() && mask > value; i++ {
		value = (value << 8) + int32(p.ReadByte())
		mask <<= 7
	}
	if mask > value {
		p.miss = true
		return 0
	}
	return value
}

// ReadSize read Data Size.
func (p *Reader) ReadSize() int64 {
	mask, value := int64(0x80), int64(p.ReadByte())
	for i := 0; i < 7 && !p.AtEOS() && mask > value; i++ {
		value = (value << 8) + int64(p.ReadByte())
		mask <<= 7
	}
	if mask > value {
		p.miss = true
		return 0
	}
	return value - mask
}

// ReadInt reads signed integer.
func (p *Reader) ReadInt() int64 {
	size := p.ReadSize()
	return p.ReadSizedInt(int(size))

}

// ReadUint reads unsigned integer.
func (p *Reader) ReadUint() uint64 {
	size := p.ReadSize()
	return p.ReadSizedUint(int(size))
}

func (p *Reader) ReadSizedInt(size int) int64 {
	if size < 1 || size > 8 {
		panic("invalid integer size")
	}
	var value int64
	i := 0
	for !p.AtEOS() && i < size {
		value = (value << 8) + int64(p.ReadByte())
		i++
	}
	if i < size {
		p.miss = true
		return 0
	}
	return value
}

func (p *Reader) ReadSizedUint(size int) uint64 {
	if size < 1 || size > 8 {
		panic("invalid integer size")
	}
	var value uint64
	i := 0
	for !p.AtEOS() && i < size {
		value = (value << 8) + uint64(p.ReadByte())
		i++
	}
	if i < size {
		p.miss = true
		return 0
	}
	return value
}

// ReadFloat reads float. The returned type is either float32 (32 bits) or float64 (64 bits).
func (p *Reader) ReadFloat() interface{} {
	size := p.ReadSize()
	switch size {
	case 4:
		return math.Float32frombits(uint32(p.ReadSizedUint(4)))
	case 8:
		return math.Float64frombits(p.ReadSizedUint(8))
	case 0:
		return nil
	default:
		panic("invalid float size")
	}
}

func (p *Reader) ReadBytes() []byte {
	b := p.ReadBinary()
	i := len(b)
	for b[i-1] == 0 {
		i--
	}
	return b[:i]
}

func (p *Reader) ReadString() string { return string(p.ReadBytes()) }

func (p *Reader) ReadBinary() []byte {
	size := int(p.ReadSize())
	binary := p.buf[p.pos : p.pos+(size)]
	p.pos += (size)
	return binary
}

func (p *Reader) PeekClass() Class {
	class := ClassA
	bits := p.buf[p.pos]
	mask := byte(0x80)
	for i := 0; i < 7; i++ {
		if bits&mask > 0 {
			return class
		}
		mask >>= 1
		class++
	}
	return 0
}

func (p *Reader) PeekId() int32 {
	pos := p.pos
	id := p.ReadId()
	p.pos = pos
	return id
}

func (p *Reader) PeekSize() int64 {
	pos := p.pos
	size := p.ReadSize()
	p.pos = pos
	return size
}

// Skip skips the following element.
func (p *Reader) Skip() {
	size := int(p.ReadSize())
	p.pos += size
}

func (p *Reader) Dive() *Reader {
	size := int(p.ReadSize())
	hi := p.pos + size
	if len(p.buf) < hi {
		hi = len(p.buf)
	}
	subReader := Reader{
		buf:   p.buf[p.pos:hi],
		start: p.pos,
		miss:  p.miss,
	}
	p.pos += size
	return &subReader
}

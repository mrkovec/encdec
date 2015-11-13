package encdec

import (
	"math"
	"encoding"
	"encoding/binary"
	"errors"
	"io"
)
var (
	ErrEncode =  errors.New("encdec: encoding error")
	ErrDecode =  errors.New("encdec: decoding error")
	ErrNoDecData =  errors.New("encdec: nothing to decode")
	ErrDecodeNotEnoughtData =  errors.New("encdec: not enought data to decode")
)

// simple encoder
// streams encoded data into []byte buffer
type Enc struct {
	err error
	buf64 [2*binary.MaxVarintLen64]byte
	encbuf []byte
	lng int
}

func NewEnc() *Enc {
	return &Enc{
		err: nil,
		lng: 0,
		encbuf: make([]byte, 0, 1024)}
}

//  implement io.WriterTo
func (e *Enc) WriteTo(w io.Writer) (int64, error) { 
	if e.err != nil {
		return 0, e.err
	}
	e.lng, e.err = w.Write(e.encbuf)
	if e.lng < len(e.encbuf) {
		e.err = ErrEncode
	}
	return int64(e.lng), e.err
}

//  discards all encoded data
func (e *Enc) Reset() { 
	if e.err != nil {
		return 
	}
	e.encbuf = e.encbuf[0:0]
}

// encode a encoding.BinaryMarshaler into buffer 
func (e *Enc) Marshaler(x encoding.BinaryMarshaler) { 
	if e.err != nil {
		return 
	}
	var buf []byte
	buf, e.err = x.MarshalBinary()
	if e.err != nil {
		return 
	}
	e.ByteSlice(buf)
}

// encode a float64 into buffer
func (e *Enc) Float64(x float64) {
	if e.err != nil {
		return 
	}
	e.Uint64(math.Float64bits(x))
}

// encode a int64 into buffer
func (e *Enc) Int64(x int64) []byte {
	if e.err != nil {
		return nil
	}
    defer func(e *Enc) {
        if r := recover(); r != nil {
            e.err = ErrEncode
        }
    }(e)

	e.lng = binary.PutVarint(e.buf64[:], x)
	if e.lng == 0 {
		e.err  = ErrEncode
		return nil
	}
	e.encbuf = append(e.encbuf, byte(e.lng))
	e.encbuf = append(e.encbuf, e.buf64[:e.lng]...)
	return e.buf64[:e.lng]
}

// encode a uint64 into buffer
func (e *Enc) Uint64(x uint64) []byte {
	if e.err != nil {
		return nil
	}
    defer func(e *Enc) {
        if r := recover(); r != nil {
            e.err = ErrEncode
        }
    }(e)

	e.lng = binary.PutUvarint(e.buf64[:], x)
	if e.lng == 0 {
		e.err  = ErrEncode
		return nil
	}
	e.encbuf = append(e.encbuf, byte(e.lng))
	e.encbuf = append(e.encbuf, e.buf64[:e.lng]...)
	return e.buf64[:e.lng]
}

// encode a slice of bytes into buffer
func (e *Enc) ByteSlice(x []byte)  {
	if e.err != nil {
		return 
	}
	if x == nil {
		e.err  = ErrEncode
		return
	}
	e.lng = len(x)
	if e.lng > 0  && e.lng < 256  {
		e.encbuf = append(e.encbuf, byte(e.lng))
	} else {
		e.encbuf = append(e.encbuf, byte(0))
		e.Uint64(uint64(e.lng))
	}
	if e.lng > 0 {
		e.encbuf = append(e.encbuf, x...)	
	}
}

// encode a byte into buffer
func (e *Enc) Byte(x byte) {
	if e.err != nil {
		return 
	}
	e.encbuf = append(e.encbuf, x)
}

// returns byte slice of encoded data
func (e Enc) Bytes() []byte {
	if e.err != nil {
		return nil
	}
	return e.encbuf
}

//  returns encoding error if any
func (e Enc) Error() error {
	return e.err
}
// returns actual length of encoded data
func (e Enc) Len() int {
	return len(e.encbuf)
}


// simple decoder
// read encoded data from []byte buffer
type Dec struct {
	err error
	i int
	lng int
	lst int
	decbuf []byte
}

func NewDec(b []byte) (d *Dec) {
	d = &Dec{
		err: nil, 
		i: 0,
		lng: 0,
		lst: 0,
		decbuf: b}
	if b == nil {
		d.err = ErrDecode
	}  
	return
}

// reads data from a io.Reader
func (d *Dec) ReadFrom(r io.Reader) (int64, error) {
	if d.err != nil {
		return 0, d.err
	}
	buf := make([]byte,256)
	n := 0
	tn := 0
	for {
		tn, d.err = r.Read(buf)
		n = n + tn
		if d.err != nil && d.err != io.EOF {
			return int64(n), d.err
		}
		if tn > 0 {
			d.decbuf = append(d.decbuf, buf[:tn]...)
		}
		if d.err != nil && d.err == io.EOF {
			d.err = nil
			return int64(n), nil
		}
	}

}

// resets decoder to initial state
func (d *Dec) Reset() { 
	if d.err != nil {
		return 
	}
	d.i = 0
}

// decode a encoding.BinaryUnmarshaler from buffer
func (d *Dec) Unmarshaler(x encoding.BinaryUnmarshaler) { 
	d.havedata()
	if d.err != nil {
		return 
	}
	if x == nil {
		d.err = ErrDecode
		return
	}
	d.err = x.UnmarshalBinary(d.ByteSlice()) 
}

// decode a float64 from buffer
func (d *Dec) Float64() float64 {
	d.havedata()
	if d.err != nil {
		return 0.0
	}
	return math.Float64frombits(d.Uint64())
}

// decode a int64 from buffer
func (d *Dec) Int64() int64 {
	d.havedata()
	if d.err != nil {
		return 0
	}
	d.lng = int(d.decbuf[d.i])
	if d.lng <= 0 {
		d.err = ErrDecode
		return 0
	}
	d.i++
	d.lst = d.i + d.lng 
	if d.lst > len(d.decbuf) {
		d.err = ErrDecodeNotEnoughtData
		return 0
	}
	var x int64
	if d.lst == len(d.decbuf) {
		x, d.i =  binary.Varint(d.decbuf[d.i:])
	} else {
		x, d.i =  binary.Varint(d.decbuf[d.i:d.lst])
	}
	if d.i <= 0 {
		d.err = ErrDecode
		return 0
	}
	d.i = d.lst
	return x
}

// decode a uint64 from buffer
func (d *Dec) Uint64() uint64 {
	d.havedata()
	if d.err != nil {
		return 0
	}
	d.lng = int(d.decbuf[d.i])
	if d.lng <= 0 {
		d.err = ErrDecode
		return 0
	}
	d.i++
	d.lst = d.i + d.lng 
	if d.lst > len(d.decbuf) {
		d.err = ErrDecodeNotEnoughtData
		return 0
	}
	var x uint64
	var i int
	if d.lst == len(d.decbuf) {
		x, i =  binary.Uvarint(d.decbuf[d.i:])
	} else {
		x, i =  binary.Uvarint(d.decbuf[d.i:d.lst])
	}
	if i <= 0 {
		d.err = ErrDecode
		return 0
	}
	d.i = d.lst
	return x
}

// decode a slice of bytes from buffer
func (d *Dec) ByteSlice() []byte {
	d.havedata()
	if d.err != nil {
		return nil
	}
	b := d.decbuf[d.i]
	d.i++
	if b > 0 && b <= 255 {
		d.lng = int(b)
	} else {
		d.lng = int(d.Uint64())
	}
	
	if d.lng == 0 {
		return []byte{}
	}
	d.lst = d.i + d.lng 
	
	if d.lst > len(d.decbuf) {
		d.err = ErrDecodeNotEnoughtData
		return nil
	}
	if d.lst == len(d.decbuf) {
		buf := d.decbuf[d.i:]
		d.i = d.lst
		return  buf
	} 
	buf := d.decbuf[d.i:d.lst]
	d.i = d.lst
	return buf
}

// decode a byte from buffer
func (d *Dec) Byte() byte {
	d.havedata()
	if d.err != nil {
		return 0
	}
	b := d.decbuf[d.i]
	d.i++
	return b
}

// skips next encoded entity
func (d *Dec) Skip() {
	d.havedata()
	if d.err != nil {
		return 
	}	

	b := d.decbuf[d.i]
	d.i++
	if b > 0 && b <= 255 {
		d.lng = int(b)
	} else {
		d.lng = int(d.Uint64())
	}
	if d.lng < 0 {
		d.err = ErrDecode
		return 
	}
	d.lst = d.i + d.lng 
	if d.lst > len(d.decbuf) {
		d.err = ErrDecode
		return
	}
	d.i = d.lst
}

//  returns decoding error if any
func (d Dec) Error() error {
	return d.err
}

// length of undecoded buffer
func (d Dec) Len() int {
	return len(d.decbuf) - d.Pos()
}

//  actual decoding position in buffer
func (d Dec) Pos() int {
	return int(d.i)
}

func (d *Dec) havedata() {
	if d.err == nil && d.i >= len(d.decbuf) {
		d.err = ErrNoDecData	
	}
}




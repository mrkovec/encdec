/*
  Package encdec contains a simple binary encoder/decoder designed for fast marshalling/unmarshalling of fixed array of types into binary stream.
*/
package encdec

import (
	"encoding"
	"encoding/binary"
	"errors"
	"io"
	"math"
)

var (
	errEncode               = errors.New("encdec: encoding error")
	errDecode               = errors.New("encdec: decoding error")
	errNoDecData            = errors.New("encdec: nothing to decode")
	errDecodeNotEnoughtData = errors.New("encdec: not enought data to decode")
)

//  Enc is a simple encoder
//  streams encoded data into []byte buffer
type Enc struct {
	err    error
	buf64  [binary.MaxVarintLen64]byte
	encbuf []byte
	lng    int
}

func NewEnc() *Enc {
	return &Enc{
		err:    nil,
		lng:    0,
		encbuf: make([]byte, 0, 1024)}
}

//  WriteTo implements io.WriterTo
func (e *Enc) WriteTo(w io.Writer) (int64, error) {
	if e.err != nil {
		return 0, e.err
	}
	e.lng, e.err = w.Write(e.encbuf)
	if e.lng < len(e.encbuf) {
		e.err = errEncode
	}
	return int64(e.lng), e.err
}

//  Reset discards all encoded data
func (e *Enc) Reset() {
	e.err = nil
	e.encbuf = e.encbuf[0:0]
}

//  Marshaler encodes a encoding.BinaryMarshaler into buffer
func (e *Enc) Marshaler(x encoding.BinaryMarshaler) {
	if e.err != nil {
		return
	}
	if x == nil {
		e.err = errEncode
		return
	}
	var buf []byte
	buf, e.err = x.MarshalBinary()
	if e.err != nil {
		return
	}
	e.ByteSlice(buf)
}

//  Float64 encodes a float64 into buffer
func (e *Enc) Float64(x float64) {
	if e.err != nil {
		return
	}
	e.Uint64(math.Float64bits(x))
}

//  Int64 encodes a int64 into buffer
func (e *Enc) Int64(x int64) []byte {
	if e.err != nil {
		return nil
	}
	// defer func(e *Enc) {
	// 	if r := recover(); r != nil {
	// 		e.err = errEncode
	// 	}
	// }(e)

	e.lng = binary.PutVarint(e.buf64[:], x)
	// if e.lng == 0 {
	// 	e.err = errEncode
	// 	return nil
	// }
	e.encbuf = append(e.encbuf, byte(e.lng))
	e.encbuf = append(e.encbuf, e.buf64[:e.lng]...)
	return e.buf64[:e.lng]
}

//  Uint64 encodes a uint64 into buffer
func (e *Enc) Uint64(x uint64) []byte {
	if e.err != nil {
		return nil
	}
	// defer func(e *Enc) {
	// 	if r := recover(); r != nil {
	// 		e.err = errEncode
	// 	}
	// }(e)

	e.lng = binary.PutUvarint(e.buf64[:], x)
	// if e.lng == 0 {
	// 	e.err = errEncode
	// 	return nil
	// }
	e.encbuf = append(e.encbuf, byte(e.lng))
	e.encbuf = append(e.encbuf, e.buf64[:e.lng]...)
	return e.buf64[:e.lng]
}

//  ByteSlice encodes a slice of bytes into buffer
func (e *Enc) ByteSlice(x []byte) {
	if e.err != nil {
		return
	}
	if x == nil {
		e.err = errEncode
		return
	}
	e.lng = len(x)
	// if e.lng > 0 && e.lng < 256 {
	// 	e.encbuf = append(e.encbuf, byte(e.lng))
	// } else {
	// 	e.encbuf = append(e.encbuf, byte(0))
	// 	e.Uint64(uint64(e.lng))
	// }
	e.Uint64(uint64(e.lng))
	if e.lng > 0 {
		e.encbuf = append(e.encbuf, x...)
	}
}

// Byte encodes a byte into buffer
// func (e *Enc) Byte(x byte) {
// 	if e.err != nil {
// 		return
// 	}
// 	e.encbuf = append(e.encbuf, byte(1), x)
// }

//  Bytes returns byte slice of encoded data
func (e Enc) Bytes() []byte {
	if e.err != nil {
		return nil
	}
	return e.encbuf
}

//  Error returns encoding error if any
func (e Enc) Error() error {
	return e.err
}

//  Len returns actual length of encoded data
func (e Enc) Len() int {
	return len(e.encbuf)
}

//  Dec is a simple decoder
//  reads encoded data from []byte buffer
type Dec struct {
	err    error
	i      int
	lng    int
	lst    int
	decbuf []byte
}

func NewDec(b []byte) (d *Dec) {
	d = &Dec{
		err:    nil,
		i:      0,
		lng:    0,
		lst:    0,
		decbuf: b}
	if b == nil {
		d.err = errDecode
	}
	return
}

//  ReadFrom reads data from a io.Reader
func (d *Dec) ReadFrom(r io.Reader) (int64, error) {
	if d.err != nil {
		return 0, d.err
	}
	buf := make([]byte, 256)
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

//  Reset resets decoder to initial state
func (d *Dec) Reset() {
	d.err = nil
	d.i = 0
}

//  Unmarshaler decodes a encoding.BinaryUnmarshaler from buffer
func (d *Dec) Unmarshaler(x encoding.BinaryUnmarshaler) {
	if d.err != nil {
		return
	}
	if d.i >= len(d.decbuf) || d.i < 0 /*overflow*/ {
		d.err = errNoDecData
		return
	}
	if x == nil {
		d.err = errDecode
		return
	}

	d.err = x.UnmarshalBinary(d.ByteSlice())
}

//  Float64 decodes a float64 from buffer
func (d *Dec) Float64() float64 {
	if d.err != nil {
		return 0.0
	}
	if d.i >= len(d.decbuf) || d.i < 0 /*overflow*/ {
		d.err = errNoDecData
		return 0.0
	}
	return math.Float64frombits(d.Uint64())
}

//  Int64 decodes a int64 from buffer
func (d *Dec) Int64() int64 {
	if d.err != nil {
		return 0
	}
	if d.i >= len(d.decbuf) || d.i < 0 /*overflow*/ {
		d.err = errNoDecData
		return 0
	}
	d.lng = int(d.decbuf[d.i])
	// if d.lng <= 0 {
	// 	d.err = errDecode
	// 	return 0
	// }
	d.i++
	d.lst = d.i + d.lng
	if d.lst > len(d.decbuf) {
		d.err = errDecodeNotEnoughtData
		return 0
	}
	var x int64
	if d.lst == len(d.decbuf) {
		x, d.i = binary.Varint(d.decbuf[d.i:])
	} else {
		x, d.i = binary.Varint(d.decbuf[d.i:d.lst])
	}
	if d.i <= 0 {
		d.err = errDecode
		return 0
	}
	d.i = d.lst
	return x
}

//  Uint64 decodes a uint64 from buffer
func (d *Dec) Uint64() uint64 {
	if d.err != nil {
		return 0
	}
	if d.i >= len(d.decbuf) || d.i < 0 /*overflow*/ {
		d.err = errNoDecData
		return 0
	}
	d.lng = int(d.decbuf[d.i])
	// if d.lng <= 0 {
	// 	d.err = errDecode
	// 	return 0
	// }
	d.i++
	d.lst = d.i + d.lng
	if d.lst > len(d.decbuf) {
		d.err = errDecodeNotEnoughtData
		return 0
	}
	var x uint64
	var i int
	if d.lst == len(d.decbuf) {
		x, i = binary.Uvarint(d.decbuf[d.i:])
	} else {
		x, i = binary.Uvarint(d.decbuf[d.i:d.lst])
	}
	if i <= 0 {
		d.err = errDecode
		return 0
	}
	d.i = d.lst
	return x
}

//  ByteSlice decodes a slice of bytes from buffer
func (d *Dec) ByteSlice() []byte {
	if d.err != nil {
		return nil
	}
	if d.i >= len(d.decbuf) || d.i < 0 /*overflow*/ {
		d.err = errNoDecData
		return nil
	}
	// b := d.decbuf[d.i]
	// d.i++
	// if b > 0 && b <= 255 {
	// 	d.lng = int(b)
	// } else {
	// 	d.lng = int(d.Uint64())
	// 	if d.lng < 0 {
	// 		d.err = errDecode
	// 		return nil
	// 	}
	// }
	d.lng = int(d.Uint64())
	if d.lng < 0 {
		d.err = errDecode
		return nil
	}
	if d.lng == 0 {
		return []byte{}
	}
	d.lst = d.i + d.lng
	if d.lst < 0 {
		d.err = errDecode
		return nil
	}
	if d.lst > len(d.decbuf) {
		d.err = errDecodeNotEnoughtData
		return nil
	}
	if d.lst == len(d.decbuf) {
		buf := d.decbuf[d.i:]
		d.i = d.lst
		return buf
	}
	buf := d.decbuf[d.i:d.lst]
	d.i = d.lst
	return buf
}

// Byte decodes a byte from buffer
// func (d *Dec) Byte() byte {
// 	if d.err != nil {
// 		return 0
// 	}
// 	if d.i >= len(d.decbuf) || d.i < 0 /*overflow*/ {
// 		d.err = errNoDecData
// 		return 0
// 	}
// 	d.i++
// 	if d.i >= len(d.decbuf) || d.i < 0 /*overflow*/ {
// 		d.err = errNoDecData
// 		return 0
// 	}
// 	b := d.decbuf[d.i]
// 	d.i++
// 	return b
// }

// Skip skips next encoded entity
// func (d *Dec) Skip() {
// 	if d.err != nil {
// 		return
// 	}
// 	if d.i >= len(d.decbuf) || d.i < 0 /*overflow*/ {
// 		d.err = errNoDecData
// 		return
// 	}
// 	b := d.decbuf[d.i]
// 	d.i++
// 	if b > 0 && b <= 255 {
// 		d.lng = int(b)
// 	} else {
// 		d.lng = int(d.Uint64())
// 	}
// 	if d.lng < 0 {
// 		d.err = errDecode
// 		return
// 	}
// 	d.lst = d.i + d.lng
// 	if d.lst > len(d.decbuf) {
// 		d.err = errDecode
// 		return
// 	}
// 	d.i = d.lst
// }

//  Error returns decoding error if any
func (d Dec) Error() error {
	return d.err
}

//  Len length of undecoded buffer
func (d Dec) Len() int {
	return len(d.decbuf) - d.Pos()
}

//  Pos returns actual decoding position in buffer
func (d Dec) Pos() int {
	return int(d.i)
}

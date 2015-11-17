package encdec

import (
	"bytes"
	"encoding/gob"
	"testing"
	"testing/quick"
	"time"
)

var e, g interface{}

func TestEncDecByte(t *testing.T) {
	enc := NewEnc()
	d := []byte{5}
	enc.ByteSlice(d)
	if enc.Error() != nil {
		t.Error(enc.Error())
	}
	dec := NewDec(enc.Bytes())
	bd := dec.ByteSlice()
	if dec.Error() != nil {
		t.Error(dec.Error())
	}
	e, g = d[0], bd[0]
	if e != g {
		t.Errorf("expected: %v (type %T) and got: %v (type %T)", e, e, g, g)
	}
}
func TestEncDecUint64(t *testing.T) {
	enc := NewEnc()
	d := uint64(5)
	enc.Uint64(d)
	if enc.Error() != nil {
		t.Error(enc.Error())
	}
	dec := NewDec(enc.Bytes())
	bd := dec.Uint64()
	if dec.Error() != nil {
		t.Error(dec.Error())
	}
	e, g = d, bd
	if e != g {
		t.Errorf("expected: %v (type %T) and got: %v (type %T)", e, e, g, g)
	}
}
func TestEncDecInt64(t *testing.T) {
	enc := NewEnc()
	d := int64(5)
	enc.Int64(d)
	if enc.Error() != nil {
		t.Error(enc.Error())
	}
	dec := NewDec(enc.Bytes())
	bd := dec.Int64()
	if dec.Error() != nil {
		t.Error(dec.Error())
	}
	e, g = d, bd
	if e != g {
		t.Errorf("expected: %v (type %T) and got: %v (type %T)", e, e, g, g)
	}
}
func TestEncDecFloat64(t *testing.T) {
	enc := NewEnc()
	d := float64(5)
	enc.Float64(d)
	if enc.Error() != nil {
		t.Error(enc.Error())
	}
	dec := NewDec(enc.Bytes())
	bd := dec.Float64()
	if dec.Error() != nil {
		t.Error(dec.Error())
	}
	e, g = d, bd
	if e != g {
		t.Errorf("expected: %v (type %T) and got: %v (type %T)", e, e, g, g)
	}
}
func TestEncDecByteSlice(t *testing.T) {
	enc := NewEnc()
	d := []byte{1, 2, 3}
	enc.ByteSlice(d)
	if enc.Error() != nil {
		t.Error(enc.Error())
	}
	dec := NewDec(enc.Bytes())
	bd := dec.ByteSlice()
	if dec.Error() != nil {
		t.Error(dec.Error())
	}
	if !bytes.Equal(d, bd) {
		t.Errorf("expected: %v (type %T) and got: %v (type %T)", d, d, bd, bd)
	}
}
func TestEncDecMarshalUnmarshal(t *testing.T) {
	enc := NewEnc()
	d := time.Now()
	enc.Marshaler(d)
	if enc.Error() != nil {
		t.Error(enc.Error())
	}
	dec := NewDec(enc.Bytes())
	var bd time.Time
	dec.Unmarshaler(&bd)
	if dec.Error() != nil {
		t.Error(dec.Error())
	}
	if !d.Equal(bd) {
		t.Errorf("expected: %v (type %T) and got: %v (type %T)", d, d, bd, bd)
	}
}
func TestQuickEncDec(t *testing.T) {
	if err := quick.Check(func(sequence []byte, x uint64, y int64, f float64, buf []byte) bool {
		enc := NewEnc()
		ti := time.Now()
		for _, s := range sequence {
			switch {
				case s <= 50:
					enc.Uint64(x)
				case s > 50 && s <= 102:
					enc.Int64(y)
				case s > 102 && s <= 153:
					enc.Float64(f)
				case s > 153 && s <= 204:
					enc.ByteSlice(buf)
				default:
					enc.Marshaler(&ti)
			}
		}

		dec := NewDec(enc.Bytes())
		if enc.Error() != nil || dec.Error() != nil || enc.Len() != dec.Len() || dec.Pos() != 0 {
			return false
		}
		for _, s := range sequence {
			switch {
				case s <= 50:
					xd := dec.Uint64()
					if dec.Error() != nil || xd != x {
						return false
					}					
				case s > 50 && s <= 102:
					yd := dec.Int64()
					if dec.Error() != nil || yd != y {
						return false
					}					
				case s > 102 && s <= 153:
					fd := dec.Float64()
					if dec.Error() != nil || fd != f {
						return false
					}					
				case s > 153 && s <= 204:
					bufd := dec.ByteSlice()
					if dec.Error() != nil || !bytes.Equal(buf, bufd) {
						return false
					}					
				default:
					var td time.Time
					dec.Unmarshaler(&td)
					if dec.Error() != nil || !ti.Equal(td) {
						return false
					}					
			}
		}
		return true
	}, nil); err != nil {
		t.Error(err)
	}
}

// func TestQuickEncDec(t *testing.T) {
// 	if err := quick.Check(func(b byte, x uint64, y int64, f float64, buf []byte, str string) bool {

// 		enc := NewEnc()
// 		enc.ByteSlice([]byte{b})
// 		enc.ByteSlice([]byte{b})
// 		enc.ByteSlice([]byte{b})
// 		enc.Reset() //clear encoded data
// 		enc.ByteSlice([]byte{b})
// 		enc.Uint64(x)
// 		enc.Int64(y)
// 		enc.Float64(f)
// 		enc.ByteSlice(buf)
// 		enc.ByteSlice([]byte(str))
// 		t := time.Now()
// 		enc.Marshaler(t)

// 		dec := NewDec(enc.Bytes())
// 		if enc.Error() != nil || dec.Error() != nil || enc.Len() != dec.Len() || dec.Pos() != 0 {
// 			return false
// 		}

// 		bd := dec.ByteSlice()[0]
// 		xd := dec.Uint64()
// 		yd := dec.Int64()
// 		dec.Reset() //start from begining
// 		bd = dec.ByteSlice()[0]
// 		xd = dec.Uint64()
// 		yd = dec.Int64()
// 		fd := dec.Float64()
// 		bufd := dec.ByteSlice()
// 		strd := string(dec.ByteSlice())
// 		var td time.Time
// 		dec.Unmarshaler(&td)

// 		if enc.Error() != nil || dec.Error() != nil || b != bd || x != xd || y != yd || f != fd || !bytes.Equal(buf, bufd) || str != strd || !t.Equal(td) || dec.Len() != 0 || dec.Pos() != enc.Len() {
// 			return false
// 		}

// 		var buffer bytes.Buffer
// 		enc.WriteTo(&buffer)  //send encoded data to buffer
// 		dec.ReadFrom(&buffer) //fill decoder from buffer
// 		bd = dec.ByteSlice()[0]
// 		dec.Skip()
// 		dec.Skip()
// 		dec.Skip()
// 		dec.Skip()
// 		dec.Skip()
// 		dec.Unmarshaler(&td)

// 		if dec.Error() != nil || b != bd || !t.Equal(td) {
// 			return false
// 		}
// 		return true
// 	}, nil); err != nil {
// 		t.Error(err)
// 	}
// }

func TestEncDecErrorPropagation(t *testing.T) {
	enc := NewEnc()
	enc.err = errEncode
	// enc.Byte(1)
	enc.ByteSlice([]byte{1})
	enc.Uint64(1)
	enc.Int64(1)
	enc.Float64(1)
	enc.Marshaler(time.Now())
	enc.Reset()
	enc.Marshaler(nil)
	enc.Reset()
	enc.ByteSlice(nil)

	dec := NewDec(enc.Bytes())
	// _ = dec.Byte()
	// dec.Skip()
	_ = dec.ByteSlice()
	_ = dec.Uint64()
	_ = dec.Int64()
	_ = dec.Float64()
	var v time.Time
	dec.Unmarshaler(&v)
	dec.Reset()
	dec.Unmarshaler(nil)

	var buffer bytes.Buffer
	enc.WriteTo(&buffer)
	dec.ReadFrom(&buffer)

	if enc.Error() == nil || dec.Error() == nil || enc.Len() > 0 || dec.Len() > 0 || dec.Pos() > 0 || buffer.Len() > 0 {
		// fmt.Println(enc.Error(), dec.Error())
		t.Error(nil)
		return
	}

}

// simple values enc/dec
func BenchmarkBasicEncodeGob(b *testing.B) {
	var (
		network bytes.Buffer
		err     error
	)
	enc := gob.NewEncoder(&network)
	for i := 0; i < b.N; i++ {
		t := newTestType()
		err = enc.Encode(&t)
		if err != nil {
			b.Error(err)
			return
		}
	}
}
func BenchmarkBasicDecodeGob(b *testing.B) {
	var (
		network bytes.Buffer
		err     error
		v       testType
	)
	enc := gob.NewEncoder(&network)
	for i := 0; i < b.N; i++ {
		t := newTestType()
		err = enc.Encode(&t)
		if err != nil {
			b.Error(err)
			return
		}
	}
	b.ResetTimer()
	dec := gob.NewDecoder(&network)
	for i := 0; i < b.N; i++ {
		err = dec.Decode(&v)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkBasicEncodeEncDec(b *testing.B) {
	enc := NewEnc()
	for i := 0; i < b.N; i++ {
		t := newTestType()
		enc.Marshaler(&t)
		if enc.Error() != nil {
			b.Error(enc.Error())
			return
		}
	}
}

func BenchmarkBasicDecodeEncDec(b *testing.B) {
	var (
		v testType
	)
	enc := NewEnc()
	for i := 0; i < b.N; i++ {
		t := newTestType()
		enc.Marshaler(&t)
		if enc.Error() != nil {
			b.Error(enc.Error())
			return
		}
	}
	b.ResetTimer()
	dec := NewDec(enc.Bytes())
	for i := 0; i < b.N; i++ {
		dec.Unmarshaler(&v)
		if dec.Error() != nil {
			b.Error(dec.Error())
			return
		}
	}
}

// slice enc/dec
func BenchmarkSliceEncodeGob(b *testing.B) {
	var (
		network bytes.Buffer
		err     error
		v       = []string{"a", "ab", "abc", "abcd"}
	)
	enc := gob.NewEncoder(&network)
	for i := 0; i < b.N; i++ {
		err = enc.Encode(v)
		if err != nil {
			b.Error(err)
			return
		}
	}
}
func BenchmarkSliceDecodeGob(b *testing.B) {

	var (
		network bytes.Buffer
		err     error
		v       = []string{"a", "ab", "abc", "abcd"}
	)
	enc := gob.NewEncoder(&network)
	for i := 0; i < b.N; i++ {
		err = enc.Encode(v)
		if err != nil {
			b.Error(err)
			return
		}
	}
	b.ResetTimer()
	dec := gob.NewDecoder(&network)
	for i := 0; i < b.N; i++ {
		v = make([]string, 0)
		err = dec.Decode(&v)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkSliceEncodeEncDec(b *testing.B) {
	var (
		v = []string{"a", "ab", "abc", "abcd"}
	)
	enc := NewEnc()
	for i := 0; i < b.N; i++ {
		enc.Uint64(uint64(len(v)))
		for _, j := range v {
			enc.ByteSlice([]byte(j))
		}
		if enc.Error() != nil {
			b.Error(enc.Error())
			return
		}
	}
}

func BenchmarkSliceDecodeEncDec(b *testing.B) {

	var (
		v = []string{"a", "ab", "abc", "abcd"}
	)
	enc := NewEnc()
	for i := 0; i < b.N; i++ {
		enc.Uint64(uint64(len(v)))
		for _, j := range v {
			enc.ByteSlice([]byte(j))
		}
		if enc.Error() != nil {
			b.Error(enc.Error())
			return
		}
	}
	b.ResetTimer()
	dec := NewDec(enc.Bytes())
	for i := 0; i < b.N; i++ {
		l := int(dec.Uint64())
		v = make([]string, l)
		for j := 0; j < l; j++ {
			v[j] = string(dec.ByteSlice())
		}
		if dec.Error() != nil {
			b.Error(dec.Error())
			return
		}
	}
}

// map enc/dec
func BenchmarkMapEncodeGob(b *testing.B) {
	var (
		network bytes.Buffer
		err     error
		v       = map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	)
	enc := gob.NewEncoder(&network)
	for i := 0; i < b.N; i++ {
		err = enc.Encode(v)
		if err != nil {
			b.Error(err)
			return
		}
	}
}
func BenchmarkMapDecodeGob(b *testing.B) {

	var (
		network bytes.Buffer
		err     error
		v       = map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	)
	enc := gob.NewEncoder(&network)
	for i := 0; i < b.N; i++ {
		err = enc.Encode(v)
		if err != nil {
			b.Error(err)
			return
		}
	}
	b.ResetTimer()
	dec := gob.NewDecoder(&network)
	for i := 0; i < b.N; i++ {
		v = make(map[string]int)
		err = dec.Decode(&v)
		if err != nil {
			b.Error(err)
			return
		}
	}
}
func BenchmarkMapEncodeEncDec(b *testing.B) {
	var (
		v = map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	)
	enc := NewEnc()
	for i := 0; i < b.N; i++ {
		enc.Uint64(uint64(len(v)))
		for k, v := range v {
			enc.ByteSlice([]byte(k))
			enc.Uint64(uint64(v))
		}
		if enc.Error() != nil {
			b.Error(enc.Error())
			return
		}
	}
}
func BenchmarkMapDecodeEncDec(b *testing.B) {

	var (
		v = map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	)
	enc := NewEnc()
	for i := 0; i < b.N; i++ {
		enc.Uint64(uint64(len(v)))
		for k, v := range v {
			enc.ByteSlice([]byte(k))
			enc.Uint64(uint64(v))
		}
		if enc.Error() != nil {
			b.Error(enc.Error())
			return
		}
	}
	b.ResetTimer()
	dec := NewDec(enc.Bytes())
	for i := 0; i < b.N; i++ {
		l := int(dec.Uint64())
		v = make(map[string]int)
		for j := 0; j < l; j++ {
			v[string(dec.ByteSlice())] = int(dec.Uint64())
		}
		if dec.Error() != nil {
			b.Error(dec.Error())
			return
		}
	}
}

type testType struct {
	A int
	B float64
	C string
	D time.Time
}

func newTestType() testType {
	return testType{123456, 0.123456, "abcdefg", time.Now()}
}
func (t *testType) MarshalBinary() ([]byte, error) {
	enc := NewEnc()
	enc.Int64(int64(t.A))
	enc.Float64(t.B)
	enc.ByteSlice([]byte(t.C))
	enc.Marshaler(&t.D)
	return enc.Bytes(), enc.Error()
}
func (t *testType) UnmarshalBinary(data []byte) error {
	dec := NewDec(data)
	t.A = int(dec.Int64())
	t.B = dec.Float64()
	t.C = string(dec.ByteSlice())
	dec.Unmarshaler(&t.D)
	return dec.Error()
}

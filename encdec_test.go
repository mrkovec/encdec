package goflat

import (
	"testing"
	"testing/quick"
	"bytes"
 	"time"
 	//"errors"
 	//"encoding"
 	// "encoding/json"
 	// "encoding/gob"
 	// "fmt"
)
var nilerr error

func TestQuickEncDec(t *testing.T) {
	if err := quick.Check(func(b byte, x uint64, y int64, f float64, buf []byte, str string) bool {
		// fmt.Printf("\n***************************** %v\n", b)
		var buffer bytes.Buffer

		enc := NewEnc()
		enc.Byte(b)
		enc.Byte(b)
		enc.Byte(b)
		enc.Reset() //clear encoded data
		enc.Byte(b)
		enc.Uint64(x)
		enc.Int64(y)
		enc.Float64(f)
		enc.ByteSlice(buf)
		enc.ByteSlice([]byte(str))
		t := time.Now()
		enc.Marshaler(t)
		

		dec := NewDec(enc.Bytes())
		bd := dec.Byte()
		xd := dec.Uint64()
		yd := dec.Int64()
		dec.Reset() //start from begining
		bd = dec.Byte()
		xd = dec.Uint64()		
		yd = dec.Int64()
		fd := dec.Float64()
		bufd := dec.ByteSlice()
		strd := string(dec.ByteSlice())
		var td time.Time
		dec.Unmarshaler(&td)

		if enc.Error() != nil || dec.Error() != nil || b != bd || x != xd || y != yd || f != fd || !bytes.Equal(buf, bufd) || str != strd || !t.Equal(td) {
			return false
		}
		enc.WriteTo(&buffer) //send encoded data to buffer
		dec.ReadFrom(&buffer) //fill decoder from buffer
		bd = dec.Byte()
		dec.Skip()
		dec.Skip()
		dec.Skip()
		dec.Skip()
		dec.Skip()
		dec.Unmarshaler(&td)

		if dec.Error() != nil || b != bd || !t.Equal(td) {
			// fmt.Println(enc.Error(), dec.Error(), b, bd)
			return false
		}

		return true
	}, nil); err != nil {
		t.Error(err)
	}
}



// func TestEncDecErrors(t *testing.T) {
// 	enc := NewEnc()
// 	enc.ByteSlice(nil)
// 	if ErrDecode != enc.Error() {
// 		t.Errorf("expected: %v (type %T) and got: %v (type %T)", ErrDecode, ErrDecode, enc.Error(), enc.Error())
// 	}	
// 	enc.Byte(byte(5))
// 	enc.Uint64(uint64(5))
// 	enc.ByteSlice([]byte{1,2,3})
// 	enc.Marshaler(time.Now())
// 	if ErrDecode != enc.Error() {
// 		t.Errorf("expected: %v (type %T) and got: %v (type %T)", ErrDecode, ErrDecode, enc.Error(), enc.Error())
// 	}	
// 	if nil != enc.Bytes() {
// 		t.Errorf("expected: nil and got: %v (type %T)", enc.Bytes(), enc.Bytes())
// 	}	
	
// 	dec := NewDec(enc.Bytes())
// 	if ErrDecode != dec.Error() {
// 		t.Errorf("expected: %v (type %T) and got: %v (type %T)", ErrDecode, ErrDecode, dec.Error(), dec.Error())
// 	}	
// 	dec.Byte()
// 	dec.Uint64()
// 	dec.ByteSlice()

// 	enc = NewEnc()
// 	enc.ByteSlice([]byte{1,2,3})


// 	// dec := NewDec(enc.Bytes())
// 	// dec.Byte()
	
// 	// e, g := ErrDecode, dec.Error()
// 	// if e != g {
// 	// 	t.Errorf("expected: %v (type %v) and got: %v (type %v)", e, reflect.TypeOf(e), g, reflect.TypeOf(g))
// 	// }		
	

// }
// func TestMarshUnmarshHeader(t *testing.T) {
// 	if err := quick.Check(func(o uint64, l uint64) bool {
// 		h := header{offset:o, dataLen:l, headerLen:1}
// 		d, err := h.MarshalBinary() 
// 		if err != nil {
// 			return false
// 		}
// 		g := header{}
// 		if err := g.UnmarshalBinary(d); err != nil {
// 			return false
// 		}
// 		if !h.Ok() || !g.Ok() || g.offset != h.offset || g.dataLen != h.dataLen {
// 			return false
// 		}
// 		return true
// 	}, nil); err != nil {
// 		t.Error(err)
// 	}
// }
// func TestMarshUnmarshRow(t *testing.T) {
// 	if err := quick.Check(func(o uint64, d []byte) bool {
// 		if len(d) == 0 {
// 			return true
// 		}
		
// 		data, err := newRow(newHeader(o, d), d).MarshalBinary()
// 		if err != nil {
// 			return false
// 		}
// 		r := &row{}
// 		if err = r.UnmarshalBinary(data); err != nil {
// 			return false
// 		}
// 		if !r.Ok() || r.h.offset != o || r.h.dataLen != uint64(len(d)) || !bytes.Equal(r.d, d) {
// 			return false
// 		}
// 		return true
// 	}, nil); err != nil {
// 		t.Error(err)
// 	}
// }


// var (
// 	buf []byte
// 	bbuf bytes.Buffer
// 	err error

// 	v testType
// 	vs []string = []string{"a", "ab", "abc", "abcd"}
// )
// // simple values enc/dec
// func BenchmarkBasicEncodeJson(b *testing.B) {
// 	bbuf.Reset()
// 	enc := json.NewEncoder(&bbuf)
// 	for i := 0; i < b.N; i++ {

// 		err = enc.Encode(newTestType())
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 	}
// }
// func BenchmarkBasicDecodeJson(b *testing.B) {
	
// 	// var (
// 	// 	network bytes.Buffer
// 	// 	err error
// 	// 	v testType
// 	// )
// 	// enc := json.NewEncoder(&network)
// 	// for i := 0; i < b.N; i++ {

// 	// 	err = enc.Encode(newTestType())
// 	// 	if err != nil {
// 	// 		b.Error(err)
// 	// 		return
// 	// 	}
// 	// }
// 	// b.ResetTimer()
// 	dec := json.NewDecoder(&bbuf)
// 	for i := 0; i < b.N; i++ {
// 		err = dec.Decode(&v)
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 	}

// }

// func BenchmarkBasicEncodeGob(b *testing.B) {
// 	bbuf.Reset()
// 	enc := gob.NewEncoder(&bbuf)
// 	for i := 0; i < b.N; i++ {
// 		err = enc.Encode(newTestType())
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 	}
// }
// func BenchmarkBasicDecodeGob(b *testing.B) {
	
// 	// var (
// 	// 	network bytes.Buffer
// 	// 	err error
// 	// 	v testType
// 	// )
// 	// enc := gob.NewEncoder(&network)
// 	// for i := 0; i < b.N; i++ {

// 	// 	err = enc.Encode(newTestType())
// 	// 	if err != nil {
// 	// 		b.Error(err)
// 	// 		return
// 	// 	}
// 	// }
// 	// b.ResetTimer()
// 	dec := gob.NewDecoder(&bbuf)
// 	for i := 0; i < b.N; i++ {
// 		err = dec.Decode(&v)
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 	}
// }

// func BenchmarkBasicEncodeEncDec(b *testing.B) {
// 	enc := NewEnc()
// 	for i := 0; i < b.N; i++ {
// 		v := newTestType()
// 		enc.Byte(v.A)
// 		enc.Uint64(uint64(v.B))
// 		enc.Float64(v.C)
// 		enc.ByteSlice([]byte(v.D))
// 		enc.Marshaler(v.E)
// 		if enc.Error() != nil {
// 			b.Error(enc.Error())
// 			return
// 		}
// 	}
// 	buf = enc.Bytes()
// }
 
// func BenchmarkBasicDecodeEncDec(b *testing.B) {
// 	var (
// 		//v testType
// 		t time.Time
// 	)
// 	// enc := NewEnc()
// 	// for i := 0; i < b.N; i++ {
// 	// 	v = newTestType()
// 	// 	enc.Byte(v.A)
// 	// 	enc.Uint64(uint64(v.B))
// 	// 	enc.Float64(v.C)
// 	// 	enc.ByteSlice([]byte(v.D))
// 	// 	enc.Marshaler(v.E)
// 	// 	if enc.Error() != nil {
// 	// 		b.Error(enc.Error())
// 	// 		return
// 	// 	}
// 	// }
// 	// b.ResetTimer()
// 	dec := NewDec(buf)
// 	for i := 0; i < b.N; i++ {
// 		v.A = dec.Byte()
// 		v.B = int(dec.Uint64())
// 		v.C = dec.Float64()
// 		v.D = string(dec.ByteSlice())
// 		dec.Unmarshaler(&t)
// 		v.E = t
// 		if dec.Error() != nil {
// 			b.Error(dec.Error())
// 			return
// 		}
// 	}
// }

// // slice enc/dec
// func BenchmarkSliceEncodeJson(b *testing.B) {
// 	var (
// 		network bytes.Buffer
// 		err error		
// 		v []string = []string{"a", "ab", "abc", "abcd"}
// 	)
// 	enc := json.NewEncoder(&network)
// 	for i := 0; i < b.N; i++ {
// 		err = enc.Encode(v)
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 	}
// }
// func BenchmarkSliceDecodeJson(b *testing.B) {
	
// 	var (
// 		network bytes.Buffer
// 		err error
// 		v []string = []string{"a", "ab", "abc", "abcd"}
// 	)
// 	enc := json.NewEncoder(&network)
// 	for i := 0; i < b.N; i++ {
// 		err = enc.Encode(v)
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 	}
// 	b.ResetTimer()
// 	dec := json.NewDecoder(&network)
// 	for i := 0; i < b.N; i++ {
// 		v = make([]string,0)
// 		err = dec.Decode(&v)
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 	}

// }
// func BenchmarkSliceEncodeGob(b *testing.B) {
// 	var (
// 		network bytes.Buffer
// 		err error		
// 		v []string = []string{"a", "ab", "abc", "abcd"}
// 	)
// 	enc := gob.NewEncoder(&network)
// 	for i := 0; i < b.N; i++ {
// 		err = enc.Encode(v)
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 	}
// }
// func BenchmarkSliceDecodeGob(b *testing.B) {
	
// 	var (
// 		network bytes.Buffer
// 		err error
// 		v []string = []string{"a", "ab", "abc", "abcd"}
// 	)
// 	enc := gob.NewEncoder(&network)
// 	for i := 0; i < b.N; i++ {
// 		err = enc.Encode(v)
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 	}
// 	b.ResetTimer()
// 	dec := gob.NewDecoder(&network)
// 	for i := 0; i < b.N; i++ {
// 		v = make([]string,0)
// 		err = dec.Decode(&v)
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 	}
// }

// func BenchmarkSliceEncodeEncDec(b *testing.B) {
// 	var (
// 		v []string = []string{"a", "ab", "abc", "abcd"}
// 	)
 
// 	enc := NewEnc()
// 	for i := 0; i < b.N; i++ {
// 		enc.Uint64(uint64(len(v)))
// 		for _, j := range v {
// 			enc.ByteSlice([]byte(j))
// 		}
// 		if enc.Error() != nil {
// 			b.Error(enc.Error())
// 			return
// 		}
// 	}
// }
 
// func BenchmarkSliceDecodeEncDec(b *testing.B) {
	
// 	var (
// 		v []string = []string{"a", "ab", "abc", "abcd"}	
// 	)
// 	enc := NewEnc()
// 	for i := 0; i < b.N; i++ {
// 		enc.Uint64(uint64(len(v)))
// 		for _, j := range v {
// 			enc.ByteSlice([]byte(j))
// 		}
// 		if enc.Error() != nil {
// 			b.Error(enc.Error())
// 			return
// 		}
// 	}
// 	b.ResetTimer()
// 	dec := NewDec(enc.Bytes())
// 	for i := 0; i < b.N; i++ {
// 		l := int(dec.Uint64())
// 		v = make([]string,l)
// 		for j := 0; j < l; j++ {
// 			v[j] = string(dec.ByteSlice())
// 		}
// 		if dec.Error() != nil {
// 			b.Error(dec.Error())
// 			return
// 		}
// 	}
// } 


// // map enc/dec
// func BenchmarkMapEncodeJson(b *testing.B) {
// 	var (
// 		network bytes.Buffer
// 		err error		
// 		v map[string]int = map[string]int{"a":1, "b":2, "c":3, "d":4}
// 	)
// 	enc := json.NewEncoder(&network)
// 	for i := 0; i < b.N; i++ {
// 		err = enc.Encode(v)
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 	}
// }
// func BenchmarkMapDecodeJson(b *testing.B) {
	
// 	var (
// 		network bytes.Buffer
// 		err error
// 		v map[string]int = map[string]int{"a":1, "b":2, "c":3, "d":4}
// 	)
// 	enc := json.NewEncoder(&network)
// 	for i := 0; i < b.N; i++ {
// 		err = enc.Encode(v)
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 	}
// 	b.ResetTimer()
// 	dec := json.NewDecoder(&network)
// 	for i := 0; i < b.N; i++ {
// 		v = make(map[string]int)
// 		err = dec.Decode(&v)
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 	}

// }
// func BenchmarkMapEncodeGob(b *testing.B) {
// 	var (
// 		network bytes.Buffer
// 		err error		
// 		v map[string]int = map[string]int{"a":1, "b":2, "c":3, "d":4}
// 	)
// 	enc := gob.NewEncoder(&network)
// 	for i := 0; i < b.N; i++ {
// 		err = enc.Encode(v)
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 	}
// }
// func BenchmarkMapDecodeGob(b *testing.B) {
	
// 	var (
// 		network bytes.Buffer
// 		err error
// 		v map[string]int = map[string]int{"a":1, "b":2, "c":3, "d":4}
// 	)
// 	enc := gob.NewEncoder(&network)
// 	for i := 0; i < b.N; i++ {
// 		err = enc.Encode(v)
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 	}
// 	b.ResetTimer()
// 	dec := gob.NewDecoder(&network)
// 	for i := 0; i < b.N; i++ {
// 		v = make(map[string]int)
// 		err = dec.Decode(&v)
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 	}
// }
// func BenchmarkMapEncodeEncDec(b *testing.B) {
// 	var (
// 		v map[string]int = map[string]int{"a":1, "b":2, "c":3, "d":4}
// 	)
//  	enc := NewEnc()
// 	for i := 0; i < b.N; i++ {
// 		enc.Uint64(uint64(len(v)))
// 		for k, v := range v {
// 			enc.ByteSlice([]byte(k))
// 			enc.Uint64(uint64(v))
// 		}
// 		if enc.Error() != nil {
// 			b.Error(enc.Error())
// 			return
// 		}
// 	}
// }
// func BenchmarkMapDecodeEncDec(b *testing.B) {
	
// 	var (
// 		v map[string]int = map[string]int{"a":1, "b":2, "c":3, "d":4}
// 	)
// 	enc := NewEnc()
// 	for i := 0; i < b.N; i++ {
// 		enc.Uint64(uint64(len(v)))
// 		for k, v := range v {
// 			enc.ByteSlice([]byte(k))
// 			enc.Uint64(uint64(v))
// 		}
// 		if enc.Error() != nil {
// 			b.Error(enc.Error())
// 			return
// 		}
// 	}
// 	b.ResetTimer()
// 	dec := NewDec(enc.Bytes())
// 	for i := 0; i < b.N; i++ {
// 		l := int(dec.Uint64())
// 		v = make(map[string]int)
// 		for j := 0; j < l; j++ {
// 			v[string(dec.ByteSlice())] = int(dec.Uint64())
// 		}
// 		if dec.Error() != nil {
// 			b.Error(dec.Error())
// 			return
// 		}
// 	}
// }
// type testType struct {
// 	A byte
// 	B int
// 	C float64
// 	D string
// 	E time.Time
// 	// F []int
// 	// F []time.Time
// }
// func newTestType() testType {
// 	return testType{1, 123456, 0.123456, "abcdefg", time.Now()/*, []int{1, 2, 3, 4, 5, 6}*//*,  []time.Time{time.Now(),time.Now(),time.Now()}*/}
// }

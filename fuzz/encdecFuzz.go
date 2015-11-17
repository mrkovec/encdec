package encdec

import (
	"bytes"
	"errors"
	"github.com/mrkovec/encdec"
	"reflect"
	"time"
)

func Fuzz(data []byte) int {
	var v fuzzTestType

	dec := encdec.NewDec(data)
	dec.Unmarshaler(&v)
	if dec.Error() != nil {
		return 0
	}
	enc := encdec.NewEnc()
	enc.Marshaler(&v)
	if enc.Error() != nil {
		panic(enc.Error())
	}
	var v1 fuzzTestType
	dec = encdec.NewDec(enc.Bytes())
	dec.Unmarshaler(&v1)
	if dec.Error() != nil {
		panic(dec.Error())
	}
	if !reflect.DeepEqual(v, v1) {
		panic("not equal")
	}
	return 1
}

type fuzzTestType struct {
	a uint64
	b int64
	c float64
	d []byte
	e time.Time
}

func (t *fuzzTestType) MarshalBinary() ([]byte, error) {
	enc := encdec.NewEnc()
	enc.Uint64(t.a)
	enc.Uint64(t.a)
	enc.Int64(t.b)
	enc.Int64(t.b)
	enc.Float64(t.c)
	enc.Float64(t.c)
	enc.ByteSlice(t.d)
	enc.ByteSlice(t.d)
	enc.Marshaler(&t.e)
	enc.Marshaler(&t.e)
	return enc.Bytes(), enc.Error()
}

func (t *fuzzTestType) UnmarshalBinary(data []byte) error {
	unmErr := errors.New("unmarshal error")

	dec := encdec.NewDec(data)
	t.a = dec.Uint64()
	if t.a != dec.Uint64() {
		return unmErr
	}
	t.b = dec.Int64()
	if t.b != dec.Int64() {
		return unmErr
	}
	t.c = dec.Float64()
	if t.c != dec.Float64() {
		return unmErr
	}
	t.d = dec.ByteSlice()
	if !bytes.Equal(t.d, dec.ByteSlice()) {
		return unmErr
	}
	dec.Unmarshaler(&t.e)
	var et time.Time
	dec.Unmarshaler(&et)
	if !t.e.Equal(et) {
		return unmErr
	}
	return dec.Error()
}

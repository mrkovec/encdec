package encdec

import (
	"github.com/dvyukov/go-fuzz/examples/fuzz"
	"github.com/mrkovec/encdec"
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
	if !fuzz.DeepEqual(v, v1) {
		panic("not equal")
	}
	return 1
}

type fuzzTestType struct {
	A byte
	B uint64
	C int64
	D float64
	E []byte
	F time.Time
}

func (t *fuzzTestType) MarshalBinary() ([]byte, error) {
	enc := encdec.NewEnc()
	enc.Byte(t.A)
	enc.Byte(t.A)
	enc.Uint64(t.B)
	enc.Uint64(t.B)
	enc.Int64(t.C)
	enc.Int64(t.C)
	enc.Float64(t.D)
	enc.Float64(t.D)
	enc.ByteSlice(t.E)
	enc.ByteSlice(t.E)
	enc.Marshaler(&t.F)
	enc.Marshaler(&t.F)
	return enc.Bytes(), enc.Error()
}
func (t *fuzzTestType) UnmarshalBinary(data []byte) error {
	dec := encdec.NewDec(data)
	dec.Skip()
	t.A = dec.Byte()
	dec.Skip()
	t.B = dec.Uint64()
	dec.Skip()
	t.C = dec.Int64()
	dec.Skip()
	t.D = dec.Float64()
	dec.Skip()
	t.E = dec.ByteSlice()
	dec.Skip()
	dec.Unmarshaler(&t.F)
	return dec.Error()
}

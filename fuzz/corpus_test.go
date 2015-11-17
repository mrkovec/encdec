package encdec

import (
	"fmt"
	"github.com/mrkovec/encdec"
	"os"
	"testing"
	"testing/quick"
	"time"
	// "reflect"
)

func TestFuzzCrashers(t *testing.T) {
	var v fuzzTestType
	var crashers = [][]byte{
		[]byte("\x00\f\xeb\xe4\xcb\xee\xff\xf7\xff\xfd\xff\x0100"),
		[]byte("\x14\x00\x12\xfb\xff\xff\xff\xff\xff\xff\xff\u007f000000000"),
		[]byte("\x00\x12\xfe\xff\xff\xff\xff\xff\xff\xff\u007f000000000"),
		[]byte("\x03\x0100")}
	for _, c := range crashers {
		dec := encdec.NewDec(c)
		dec.Unmarshaler(&v)
	}
	// var errors = [][]byte{
	// 	[]byte("\x02\xb6\x01\t000000000\t000000000\t000000000\t000000000\n0000000000\n0000000000\x010W\x02000000000000000000000000000000000000000000000000000" +
	// 			"00000000000000000000000000000000000\x000\n0000000000\x01\x0f\x0100000000000000")}
	// for _, e := range errors {
	// 	dec := encdec.NewDec(e)
	// 	dec.Unmarshaler(&v)
	// 	if dec.Error() != nil {
	// 		panic(dec.Error())
	// 	}
	// 	enc := encdec.NewEnc()
	// 	enc.Marshaler(&v)
	// 	if enc.Error() != nil {
	// 		panic(enc.Error())
	// 	}
	// 	fmt.Printf("%v\n%v\n%v\n", e, v, enc.Bytes())
	// 	var v1 fuzzTestType
	// 	dec = encdec.NewDec(enc.Bytes())
	// 	dec.Unmarshaler(&v1)
	// 	if dec.Error() != nil {
	// 		panic(dec.Error())
	// 	}
	// 	if !reflect.DeepEqual(v, v1) {
	// 		panic("not equal")
	// 	}
	// }	
}

func testGenerateCorpus(t *testing.T) {
	i := 0
	if err := quick.Check(func(x1 uint64, y1 int64, f1 float64, buf1 []byte) bool {
		fi, _ := os.Create(fmt.Sprintf("corpus/init%v", i))
		i++

		v := fuzzTestType{x1, y1, f1, buf1, time.Now()}
		enc := encdec.NewEnc()
		enc.Marshaler(&v)
		fi.Write(enc.Bytes())
		fi.Close()

		return true
	}, nil); err != nil {
		t.Error(err)
	}
}

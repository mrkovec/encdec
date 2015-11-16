package encdec

import (
	"fmt"
	"github.com/mrkovec/encdec"
	"os"
	"testing"
	"testing/quick"
	"time"
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
}

func testGenerateCorpus(t *testing.T) {
	i := 0
	if err := quick.Check(func(b1 byte, x1 uint64, y1 int64, f1 float64, buf1 []byte) bool {
		fi, _ := os.Create(fmt.Sprintf("corpus/%v", i))
		i++

		v := fuzzTestType{b1, x1, y1, f1, buf1, time.Now()}
		enc := encdec.NewEnc()
		enc.Marshaler(&v)
		fi.Write(enc.Bytes())
		fi.Close()

		return true
	}, nil); err != nil {
		t.Error(err)
	}
}

# encdec
[![GoDoc](https://godoc.org/github.com/mrkovec/encdec?status.svg)](https://godoc.org/github.com/mrkovec/encdec)
[![Coverage Status](https://coveralls.io/repos/mrkovec/encdec/badge.svg?branch=master)](https://coveralls.io/r/mrkovec/encdec?branch=master)

**encdec** contains a simple binary encoder/decoder designed for fast marshalling/unmarshalling of fixed array of types into binary stream.

Functionality outline:
```go
    //setup
    import "github.com/mrkovec/encdec"
    
    type user struct {
        name string
        age int
        registered time.Time
    }
    func (u *user) MarshalBinary() ([]byte, error) {
        enc := encdec.NewEnc()
        enc.ByteSlice([]byte(u.name))
        enc.Int64(int64(u.age))
        enc.Marshaler(u.registered)
        return enc.Bytes(), enc.Error()
    }
    func (u *user) UnmarshalBinary(data []byte) error {
        dec := encdec.NewDec(data)
        u.name = string(dec.ByteSlice())
        u.age = int(dec.Int64())
        dec.Unmarshaler(&u.registered)
        return dec.Error()
    }
    
    var users []user
    
    //main functionality
    //encode data to byte slice
    enc := encdec.NewEnc()
    enc.Uint64(uint64(len(users)))
    for _, u := range users {
        enc.Marshaler(&u)
    }
    if enc.Error() != nil {
        panic(enc.Error())
    }
    slice := enc.Bytes()
    
    //decode data from byte slice
    dec := encdec.NewDec(slice)
    l := int(dec.Uint64())
    users = make([]user, l)
    for i := 0; i < l; i++ {
        dec.Unmarshaler(&users[i])
    }    
    if dec.Error() != nil {
        panic(dec.Error())
    }    
```
Alternatively encdec can write to or read from arbitrary io.Reader/io.Writer
```go
    //encode
    enc := encdec.NewEnc()
    ...
    network bytes.Buffer
    enc.WriteTo(network)
    
    //decode
    dec := encdec.NewDec([]byte{})
    dec.ReadFrom(network)
    ...
```
For more examples look in GoDoc or in test/benchmark files.
package igbinary

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/zarken-go/igbinary/igcode"
	"reflect"
)

type decodeTest struct {
	expected interface{}
	out      interface{}
	data     []byte
	hex      string
	errStr   string
}

func (t decodeTest) String() string {
	return fmt.Sprintf("expected=%#v, out=%#v", t.expected, t.out)
}

func (Suite *DecodeSuite) TestTypes() {
	Tests := []decodeTest{
		{expected: uint8(10), out: new(uint8), data: []byte{igcode.PosInt8, 10}},
		{expected: uint8(255), out: new(uint8), data: []byte{igcode.PosInt16, 0, 255}},
		{expected: uint8(0), out: new(uint8), data: []byte{igcode.PosInt16, 1, 0},
			errStr: `igbinary: Decode(unsigned: int 256 out of range [0:255])`},
		{expected: uint8(0), out: new(uint8), data: []byte{igcode.NegInt16, 0, 1},
			errStr: `igbinary: Decode(unsigned: int -1 out of range [0:255])`},
		{expected: uint8(0), out: new(uint8), hex: `0800`, errStr: `unexpected EOF`},

		{expected: uint16(10), out: new(uint16), data: []byte{igcode.PosInt8, 10}},
		{expected: uint16(65535), out: new(uint16), data: []byte{igcode.PosInt16, 0xff, 0xff}},
		{expected: uint16(0), out: new(uint16), data: []byte{igcode.PosInt32, 0, 1, 0, 0},
			errStr: `igbinary: Decode(unsigned: int 65536 out of range [0:65535])`},
		{expected: uint16(0), out: new(uint16), data: []byte{igcode.NegInt16, 0, 1},
			errStr: `igbinary: Decode(unsigned: int -1 out of range [0:65535])`},
		{expected: uint16(0), out: new(uint16), hex: `0800`, errStr: `unexpected EOF`},

		{expected: uint32(10), out: new(uint32), hex: `060a`},
		{expected: uint32(0xffffffff), out: new(uint32), hex: `0affffffff`},
		{expected: uint32(0), out: new(uint32), data: []byte{igcode.PosInt64, 0, 0, 0, 1, 0, 0, 0, 0},
			errStr: `igbinary: Decode(unsigned: int 4294967296 out of range [0:4294967295])`},
		{expected: uint32(0), out: new(uint32), data: []byte{igcode.NegInt32, 0, 0, 0, 1},
			errStr: `igbinary: Decode(unsigned: int -1 out of range [0:4294967295])`},
		{expected: uint32(0), out: new(uint32), hex: `0800`, errStr: `unexpected EOF`},

		{expected: uint64(10), out: new(uint64), data: []byte{igcode.PosInt8, 10}},
		{expected: uint64(0xffffffffffffffff), out: new(uint64), hex: `20ffffffffffffffff`},
		{expected: uint64(0), out: new(uint64), data: []byte{igcode.NegInt32, 0, 0, 0, 1},
			errStr: `igbinary: Decode(unsigned: int -1 out of range [0:18446744073709551615])`},
		{expected: uint64(0), out: new(uint64), hex: `0800`, errStr: `unexpected EOF`},

		{expected: uint(0), out: new(uint), data: []byte{igcode.PosInt8, 0}},
		{expected: uint(0xffffffffffffffff), out: new(uint), hex: `20ffffffffffffffff`},
		{expected: uint(0), out: new(uint), hex: `0701`,
			errStr: `igbinary: Decode(unsigned: int -1 out of range [0:18446744073709551615])`},
		{expected: uint(0), out: new(uint), hex: `0800`, errStr: `unexpected EOF`},

		{expected: 0, out: new(int), data: []byte{igcode.PosInt8, 0}},
		{expected: 0x7fffffffffffffff, out: new(int), hex: `207fffffffffffffff`},
		{expected: 0, out: new(int), hex: `208000000000000000`,
			errStr: `igbinary: Decode(signed: int 9223372036854775808 out of range [-9223372036854775808:9223372036854775807])`},
		{expected: 0, out: new(int), hex: `218000000000000001`,
			errStr: `igbinary: Decode(signed: int -9223372036854775809 out of range [-9223372036854775808:9223372036854775807])`},
		{expected: 0, out: new(int), hex: `0800`, errStr: `unexpected EOF`},

		{expected: 0, out: new(int), hex: ``, errStr: `EOF`},
		{expected: 0, out: new(int), hex: `06`, errStr: `EOF`},
		{expected: 0, out: new(int), hex: `0900`, errStr: `unexpected EOF`},
		{expected: 0, out: new(int), hex: `0b0100`, errStr: `unexpected EOF`},
		{expected: 0, out: new(int), hex: `61`, errStr: `igbinary: Decode(readInteger unexpected code 'a')`},
	}

	for _, Test := range Tests {
		if Test.hex != `` && Test.data == nil {
			var err error
			if Test.data, err = hex.DecodeString(Test.hex); err != nil {
				Suite.T().Fatalf(`cannot decode hex test data: %s - (%s)`, err, Test.hex)
			}
		}

		buffer := bytes.NewBuffer(Test.data)
		decoder := NewDecoder(buffer)
		err := decoder.Decode(Test.out)

		Suite.assertNilOrError(err, Test.errStr)
		Suite.Equal(indirect(Test.expected), indirect(Test.out))
		// ensure nothing is left in the buffer
		if buffer.Len() > 0 {
			Suite.T().Fatalf("unread data in the buffer: %q (%s)", buffer.Bytes(), Test)
		}

		buffer = bytes.NewBuffer(Test.data)
		decoder = NewDecoder(buffer)
		decodeDest := reflect.ValueOf(Test.out).Elem()
		decoderF := getDecoder(decodeDest.Type())
		err = decoderF(decoder, decodeDest)

		Suite.assertNilOrError(err, Test.errStr)
		Suite.Equal(indirect(Test.expected), indirect(Test.out))
		// ensure nothing is left in the buffer
		if buffer.Len() > 0 {
			Suite.T().Fatalf("unread data in the buffer: %q (%s)", buffer.Bytes(), Test)
		}
	}
}

func indirect(viface interface{}) interface{} {
	v := reflect.ValueOf(viface)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.IsValid() {
		return v.Interface()
	}
	return nil
}

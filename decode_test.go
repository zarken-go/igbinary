package igbinary

import (
	"encoding/hex"
	"github.com/stretchr/testify/suite"
	"github.com/zarken-go/igbinary/igcode"
	"testing"
)

type DecodeSuite struct {
	suite.Suite
}

func (Suite *DecodeSuite) TestStrings() {
	Suite.assertUnmarshalString(`foobar`, `1106666f6f626172`)
	Suite.assertUnmarshalString(`foobar`, `120006666f6f626172`)
	Suite.assertUnmarshalString(`foobar`, `1300000006666f6f626172`)
	Suite.assertUnmarshalString(``, `0d`)
	Suite.assertUnmarshalString(``, `00`)
}

func (Suite *DecodeSuite) assertUnmarshalString(Expected string, Hex string) {
	var dest string
	b, err := hex.DecodeString(Hex)
	if Suite.Nil(err) {
		Suite.Nil(Unmarshal(b, &dest))
		Suite.Equal(Expected, dest)
	}
}

func (Suite *DecodeSuite) TestDecodeInt8() {
	Suite.assertUnmarshalInt8(0, []byte{igcode.PosInt8, 0}, ``)
	Suite.assertUnmarshalInt8(0, []byte{igcode.NegInt8, 0}, ``)
	Suite.assertUnmarshalInt8(127, []byte{igcode.PosInt8, 0x7f}, ``)
	Suite.assertUnmarshalInt8(-128, []byte{igcode.NegInt16, 0, 128}, ``)
	Suite.assertUnmarshalInt8(-64, []byte{igcode.NegInt32, 0, 0, 0, 64}, ``)

	Suite.assertUnmarshalInt8(0, []byte{igcode.PosInt16, 0x00, 0x80},
		`igbinary: Decode(signed: int 128 out of range [-128:127])`)
	Suite.assertUnmarshalInt8(0, []byte{igcode.NegInt16, 0x00, 0x81},
		`igbinary: Decode(signed: int -129 out of range [-128:127])`)
	Suite.assertUnmarshalInt8(0, []byte{igcode.PosInt16, 0x00}, `unexpected EOF`)
}

func (Suite *DecodeSuite) assertUnmarshalInt8(expected int8, data []byte, errStr string) {
	var v int8
	err := Unmarshal(data, &v)
	Suite.assertNilOrError(err, errStr)
	Suite.Equal(expected, v)
}

func (Suite *DecodeSuite) TestDecodeInt16() {
	Suite.assertUnmarshalInt16(0, []byte{igcode.PosInt8, 0}, ``)
	Suite.assertUnmarshalInt16(0, []byte{igcode.NegInt8, 0}, ``)
	Suite.assertUnmarshalInt16(32767, []byte{igcode.PosInt16, 0x7f, 0xff}, ``)
	Suite.assertUnmarshalInt16(-32768, []byte{igcode.NegInt16, 0x80, 0x00}, ``)

	Suite.assertUnmarshalInt16(0, []byte{igcode.PosInt16, 0x80, 0x00},
		`igbinary: Decode(signed: int 32768 out of range [-32768:32767])`)
	Suite.assertUnmarshalInt16(0, []byte{igcode.NegInt16, 0x80, 0x01},
		`igbinary: Decode(signed: int -32769 out of range [-32768:32767])`)
	Suite.assertUnmarshalInt16(0, []byte{igcode.PosInt16, 0x00}, `unexpected EOF`)
}

func (Suite *DecodeSuite) assertUnmarshalInt16(expected int16, data []byte, errStr string) {
	var v int16
	err := Unmarshal(data, &v)
	Suite.assertNilOrError(err, errStr)
	Suite.Equal(expected, v)
}

func (Suite *DecodeSuite) TestDecodeInt32() {
	Suite.assertUnmarshalInt32(0, []byte{igcode.PosInt8, 0}, ``)
	Suite.assertUnmarshalInt32(0, []byte{igcode.NegInt8, 0}, ``)
	Suite.assertUnmarshalInt32(2147483647, []byte{igcode.PosInt32, 0x7f, 0xff, 0xff, 0xff}, ``)
	Suite.assertUnmarshalInt32(-2147483648, []byte{igcode.NegInt32, 0x80, 0x00, 0x00, 0x00}, ``)

	Suite.assertUnmarshalInt32(0, []byte{igcode.PosInt32, 0x80, 0x00, 0x00, 0x00},
		`igbinary: Decode(signed: int 2147483648 out of range [-2147483648:2147483647])`)
	Suite.assertUnmarshalInt32(0, []byte{igcode.NegInt32, 0x80, 0x00, 0x00, 0x01},
		`igbinary: Decode(signed: int -2147483649 out of range [-2147483648:2147483647])`)
	Suite.assertUnmarshalInt32(0, []byte{igcode.PosInt16, 0x00}, `unexpected EOF`)
}

func (Suite *DecodeSuite) assertUnmarshalInt32(expected int32, data []byte, errStr string) {
	var v int32
	err := Unmarshal(data, &v)
	Suite.assertNilOrError(err, errStr)
	Suite.Equal(expected, v)
}

func (Suite *DecodeSuite) TestDecodeInt64() {
	Suite.assertUnmarshalInt64(0, []byte{igcode.PosInt8, 0}, ``)
	Suite.assertUnmarshalInt64(0, []byte{igcode.NegInt8, 0}, ``)
	Suite.assertUnmarshalInt64(9223372036854775807, []byte{igcode.PosInt64, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, ``)
	Suite.assertUnmarshalInt64(-9223372036854775808, []byte{igcode.NegInt64, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ``)

	Suite.assertUnmarshalInt64(0, []byte{igcode.PosInt64, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		`igbinary: Decode(signed: int 9223372036854775808 out of range [-9223372036854775808:9223372036854775807])`)
	Suite.assertUnmarshalInt64(0, []byte{igcode.NegInt64, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		`igbinary: Decode(signed: int -9223372036854775809 out of range [-9223372036854775808:9223372036854775807])`)
	Suite.assertUnmarshalInt64(0, []byte{igcode.PosInt16, 0x00}, `unexpected EOF`)
}

func (Suite *DecodeSuite) assertUnmarshalInt64(expected int64, data []byte, errStr string) {
	var v int64
	err := Unmarshal(data, &v)
	Suite.assertNilOrError(err, errStr)
	Suite.Equal(expected, v)
}

func (Suite *DecodeSuite) assertNilOrError(actualErr error, expectedErr string) {
	if expectedErr == `` {
		Suite.Nil(actualErr, `unexpected error: %s`, actualErr)
	} else {
		Suite.EqualError(actualErr, expectedErr)
	}
}

func (Suite *DecodeSuite) TestDecodeUint8() {
	var v uint8
	var err error

	Suite.Nil(Unmarshal([]byte{igcode.PosInt8, 0}, &v))
	Suite.Equal(uint8(0), v)

	Suite.Nil(Unmarshal([]byte{igcode.NegInt8, 0}, &v))
	Suite.Equal(uint8(0), v)

	Suite.Nil(Unmarshal([]byte{igcode.PosInt16, 0, 255}, &v))
	Suite.Equal(uint8(255), v)

	Suite.Nil(Unmarshal([]byte{igcode.PosInt32, 0, 0, 0, 64}, &v))
	Suite.Equal(uint8(64), v)

	err = Unmarshal([]byte{igcode.PosInt32, 0, 0, 1, 0}, &v)
	Suite.EqualError(err, `igbinary: Decode(unsigned: int 256 out of range [0:255])`)
	Suite.Equal(uint8(0), v)

	err = Unmarshal([]byte{igcode.NegInt32, 0, 0, 0, 1}, &v)
	Suite.EqualError(err, `igbinary: Decode(unsigned: int -1 out of range [0:255])`)
	Suite.Equal(uint8(0), v)

	err = Unmarshal([]byte{igcode.NegInt64, 0, 0, 1, 0}, &v)
	Suite.EqualError(err, `unexpected EOF`)
	Suite.Equal(uint8(0), v)
}

func TestDecodeSuite(t *testing.T) {
	suite.Run(t, new(DecodeSuite))
}

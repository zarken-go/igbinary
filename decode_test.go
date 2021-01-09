package igbinary

import (
	"encoding/hex"
	"github.com/stretchr/testify/suite"
	"github.com/zarken-go/igbinary/igcode"
	"io"
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
	var v int8
	var err error

	Suite.Nil(Unmarshal([]byte{igcode.PosInt8, 127}, &v))
	Suite.Equal(int8(127), v)

	Suite.Nil(Unmarshal([]byte{igcode.NegInt16, 0, 128}, &v))
	Suite.Equal(int8(-128), v)

	Suite.Nil(Unmarshal([]byte{igcode.NegInt32, 0, 0, 0, 64}, &v))
	Suite.Equal(int8(-64), v)

	err = Unmarshal([]byte{igcode.NegInt32, 0, 0, 0, 129}, &v)
	Suite.EqualError(err, `igbinary: Decode(signed: int -129 out of range [-128:127])`)
	Suite.Equal(int8(0), v)

	err = Unmarshal([]byte{igcode.PosInt32, 0, 0, 0, 128}, &v)
	Suite.EqualError(err, `igbinary: Decode(signed: int 128 out of range [-128:127])`)
	Suite.Equal(int8(0), v)

	err = Unmarshal([]byte{igcode.NegInt64, 0, 0, 1, 0}, &v)
	Suite.EqualError(io.EOF, `EOF`)
	Suite.Equal(int8(0), v)
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
	Suite.EqualError(io.EOF, `EOF`)
	Suite.Equal(uint8(0), v)
}

func TestDecodeSuite(t *testing.T) {
	suite.Run(t, new(DecodeSuite))
}

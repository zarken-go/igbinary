package igbinary

import (
	"bytes"
	"encoding/hex"
	"github.com/stretchr/testify/suite"
	"testing"
)

type EncodeSuite struct {
	suite.Suite
}

func (Suite *EncodeSuite) TestNils() {
	Suite.assertMarshal(nil, `00`)
}

func (Suite *EncodeSuite) TestBooleans() {
	Suite.assertMarshal(true, `05`)
	Suite.assertMarshal(false, `04`)
}

func (Suite *EncodeSuite) TestIntegers() {
	Suite.assertMarshal(0, `0600`)
	Suite.assertMarshal(1, `0601`)
	Suite.assertMarshal(-1, `0701`)
	Suite.assertMarshal(255, `06ff`)
	Suite.assertMarshal(-255, `07ff`)
	Suite.assertMarshal(1000, `0803e8`)
	Suite.assertMarshal(-1000, `0903e8`)
	Suite.assertMarshal(100000, `0a000186a0`)
	Suite.assertMarshal(-100000, `0b000186a0`)
}

func (Suite *EncodeSuite) TestFloats() {
	Suite.assertMarshal(123.456, `0c405edd2f1a9fbe77`)
}

func (Suite *EncodeSuite) TestStrings() {
	Suite.assertMarshal(`foobar`, `1106666f6f626172`)
	Suite.assertMarshal([]byte(`foobar`), `1106666f6f626172`)
}

func (Suite *EncodeSuite) TestEncodeArrayLen() {
	b := &bytes.Buffer{}
	Encoder := NewEncoder(b)
	Suite.Nil(Encoder.EncodeArrayLen(10))
	Suite.Equal([]byte{0x14, 0xa}, b.Bytes())

	b.Reset()
	Suite.Nil(Encoder.EncodeArrayLen(300))
	Suite.Equal([]byte{0x15, 0x1, 0x2c}, b.Bytes())

	b.Reset()
	Suite.Nil(Encoder.EncodeArrayLen(0xfffff))
	Suite.Equal([]byte{0x16, 0x0, 0xf, 0xff, 0xff}, b.Bytes())

	b.Reset()
	err := Encoder.EncodeArrayLen(0xfffffffff)
	Suite.EqualError(err, `igbinary: Encode(unsupported array length 68719476735)`)
}

func (Suite *EncodeSuite) assertMarshal(v interface{}, expectedHex string) {
	b, err := Marshal(v)
	Suite.Nil(err)
	Suite.Equal(`00000002`+expectedHex, hex.EncodeToString(b))
}

func TestEncodeSuite(t *testing.T) {
	suite.Run(t, new(EncodeSuite))
}

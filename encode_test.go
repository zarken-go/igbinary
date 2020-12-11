package igbinary

import (
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

func (Suite *EncodeSuite) assertMarshal(v interface{}, expectedHex string) {
	b, err := Marshal(v)
	Suite.Nil(err)
	Suite.Equal(`00000002`+expectedHex, hex.EncodeToString(b))
}

func TestEncodeSuite(t *testing.T) {
	suite.Run(t, new(EncodeSuite))
}

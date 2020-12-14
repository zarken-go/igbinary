package igbinary

import (
	"encoding/hex"
	"github.com/stretchr/testify/suite"
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

func TestDecodeSuite(t *testing.T) {
	suite.Run(t, new(DecodeSuite))
}

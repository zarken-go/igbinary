package igbinary

import (
	"fmt"
	"github.com/zarken-go/igbinary/igcode"
	"reflect"
)

func (d *Decoder) DecodeString() (string, error) {
	c, err := d.readCode()
	if err != nil {
		return "", err
	}
	return d.string(c)
}

func (d *Decoder) bytesLen(c byte) (int, error) {
	switch c {
	case igcode.Nil:
		return -1, nil
	case igcode.StringEmpty:
		return 0, nil
	case igcode.String8:
		n, err := d.uint8()
		return int(n), err
	case igcode.String16:
		n, err := d.uint16()
		return int(n), err
	case igcode.String32:
		n, err := d.uint32()
		return int(n), err
	}

	return 0, fmt.Errorf("igbinary: invalid code=%x decoding string/bytes length", c)
}

func (d *Decoder) string(c byte) (string, error) {
	n, err := d.bytesLen(c)
	if err != nil {
		return "", err
	}

	return d.stringWithLen(n)
}

func (d *Decoder) stringWithLen(n int) (string, error) {
	if n <= 0 {
		return "", nil
	}
	b, err := d.readN(n)
	return string(b), err
}

func decodeStringValue(d *Decoder, v reflect.Value) error {
	s, err := d.DecodeString()
	if err != nil {
		return err
	}
	v.SetString(s)
	return nil
}

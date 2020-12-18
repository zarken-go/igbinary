package igbinary

import "github.com/zarken-go/igbinary/igcode"

func (d *Decoder) uint8() (uint8, error) {
	c, err := d.readCode()
	if err != nil {
		return 0, err
	}
	return c, nil
}

func (d *Decoder) uint16() (uint16, error) {
	b, err := d.readN(2)
	if err != nil {
		return 0, err
	}
	return (uint16(b[0]) << 8) | uint16(b[1]), nil
}

func (d *Decoder) uint32() (uint32, error) {
	b, err := d.readN(4)
	if err != nil {
		return 0, err
	}
	n := (uint32(b[0]) << 24) |
		(uint32(b[1]) << 16) |
		(uint32(b[2]) << 8) |
		uint32(b[3])
	return n, nil
}

func (d *Decoder) DecodeInt8() (int8, error) {
	code, value, err := d.readInteger()
	if err != nil {
		return 0, err
	}
	if igcode.IsNegative(code) && value <= 127 {
		return -int8(value), nil
	} else if igcode.IsNegative(code) && value == 128 {
		return -128, nil
	} else if !igcode.IsNegative(code) && value <= 127 {
		return int8(value), nil
	}
	return 0, decodeErrorF(`int8 out of range`)
}
func (d *Decoder) readInteger() (byte, uint64, error) {
	code, err := d.readCode()
	if err != nil {
		return 0, 0, err
	}
	switch code {
	case igcode.PosInt8, igcode.NegInt8:
		b, err := d.readCode()
		if err != nil {
			return code, 0, err
		}
		return code, uint64(b), nil
	case igcode.PosInt16, igcode.NegInt16:
		b, err := d.readN(2)
		if err != nil {
			return code, 0, err
		}
		n := (uint64(b[0]) << 8) |
			uint64(b[1])
		return code, n, nil
	case igcode.PosInt32, igcode.NegInt32:
		b, err := d.readN(4)
		if err != nil {
			return code, 0, err
		}
		n := (uint64(b[0]) << 24) |
			(uint64(b[1]) << 16) |
			(uint64(b[2]) << 8) |
			uint64(b[3])
		return code, n, nil
	default:
		return code, 0, decodeErrorF(`readInteger unexpected code '%c'`, code)
	}
}

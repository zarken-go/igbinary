package igbinary

import (
	"github.com/zarken-go/igbinary/igcode"
)

const uintSize = 32 << (^uint(0) >> 32 & 1) // 32 or 64

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
	value, err := d.decodeSignedInt(0x7f)
	if err != nil {
		return 0, err
	}
	return int8(value), nil
}

func (d *Decoder) DecodeUint8() (uint8, error) {
	value, err := d.decodeUnsignedInt(0xff)
	if err != nil {
		return 0, err
	}
	return uint8(value), nil
}

func (d *Decoder) DecodeInt16() (int16, error) {
	value, err := d.decodeSignedInt(0x7fff)
	if err != nil {
		return 0, err
	}
	return int16(value), nil
}

func (d *Decoder) DecodeUint16() (uint16, error) {
	value, err := d.decodeUnsignedInt(0xffff)
	if err != nil {
		return 0, err
	}
	return uint16(value), nil
}

func (d *Decoder) DecodeInt32() (int32, error) {
	value, err := d.decodeSignedInt(0x7fffffff)
	if err != nil {
		return 0, err
	}
	return int32(value), nil
}

func (d *Decoder) DecodeUint32() (uint32, error) {
	value, err := d.decodeUnsignedInt(0xffffffff)
	if err != nil {
		return 0, err
	}
	return uint32(value), nil
}

func (d *Decoder) DecodeInt64() (int64, error) {
	value, err := d.decodeSignedInt(0x7fffffffffffffff)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (d *Decoder) DecodeUint64() (uint64, error) {
	value, err := d.decodeUnsignedInt(0xffffffffffffffff)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (d *Decoder) DecodeInt() (int, error) {
	var limit uint64 = 0x7fffffff // 32-bit system
	if uintSize == 64 {
		limit = 0x7fffffffffffffff // 64-bit system
	}
	value, err := d.decodeSignedInt(limit)
	if err != nil {
		return 0, err
	}
	return int(value), nil
}

func (d *Decoder) DecodeUint() (uint, error) {
	var limit uint64 = 0xffffffff // 32-bit system
	if uintSize == 64 {
		limit = 0xffffffffffffffff // 64-bit system
	}
	value, err := d.decodeUnsignedInt(limit)
	if err != nil {
		return 0, err
	}
	return uint(value), nil
}

func (d *Decoder) decodeSignedInt(limit uint64) (int64, error) {
	code, value, err := d.readInteger()
	if err != nil {
		return 0, err
	}
	if value <= limit {
		if igcode.IsNegative(code) {
			return -int64(value), nil
		}
		return int64(value), nil
	} else if value == limit+1 && igcode.IsNegative(code) {
		return -int64(limit) - 1, nil
	}
	if igcode.IsNegative(code) {
		return 0, decodeErrorF(`signed: int -%d out of range [-%d:%d]`,
			value, limit+1, limit)
	}
	return 0, decodeErrorF(`signed: int %d out of range [-%d:%d]`,
		value, limit+1, limit)
}

func (d *Decoder) decodeUnsignedInt(limit uint64) (uint64, error) {
	code, value, err := d.readInteger()
	if err != nil {
		return 0, err
	}
	if value == 0 {
		return 0, nil
	}
	if !igcode.IsNegative(code) && value <= limit {
		return value, nil
	}
	if igcode.IsNegative(code) {
		return 0, decodeErrorF(`unsigned: int -%d out of range [0:%d]`,
			value, limit)
	}
	return 0, decodeErrorF(`unsigned: int %d out of range [0:%d]`,
		value, limit)
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

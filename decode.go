package igbinary

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/zarken-go/igbinary/igcode"
	"io"
	"reflect"
)

const (
	disallowUnknownFieldsFlag uint32 = 1 << iota
)

type bufReader interface {
	io.Reader
	io.ByteScanner
}

type Decoder struct {
	r io.Reader
	s io.ByteScanner

	flags uint32

	buf []byte
	rec []byte // accumulates read data if not nil
}

func Unmarshal(data []byte, v interface{}) error {
	d := NewDecoder(bytes.NewReader(data))
	return d.Decode(v)
}

func NewDecoder(r io.Reader) *Decoder {
	d := new(Decoder)
	d.resetReader(r)
	return d
}

func (d *Decoder) resetReader(r io.Reader) {
	if br, ok := r.(bufReader); ok {
		d.r = br
		d.s = br
	} else {
		br := bufio.NewReader(r)
		d.r = br
		d.s = br
	}
}

//nolint:gocyclo
func (d *Decoder) Decode(v interface{}) error {
	var err error
	switch v := v.(type) {
	case *string:
		if v != nil {
			*v, err = d.DecodeString()
			return err
		}
		//case *[]byte:
		//	if v != nil {
		//		return ErrUnsupported // d.decodeBytesPtr(v)
		//	}
	case *int:
		if v != nil {
			*v, err = d.DecodeInt()
			return err
		}
	case *int8:
		if v != nil {
			*v, err = d.DecodeInt8()
			return err
		}
	case *int16:
		if v != nil {
			*v, err = d.DecodeInt16()
			return err
		}
	case *int32:
		if v != nil {
			*v, err = d.DecodeInt32()
			return err
		}
	case *int64:
		if v != nil {
			*v, err = d.DecodeInt64()
			return err
		}
	case *uint:
		if v != nil {
			*v, err = d.DecodeUint()
			return err
		}
	case *uint8:
		if v != nil {
			*v, err = d.DecodeUint8()
			return err
		}
	case *uint16:
		if v != nil {
			*v, err = d.DecodeUint16()
			return err
		}
	case *uint32:
		if v != nil {
			*v, err = d.DecodeUint32()
			return err
		}
	case *uint64:
		if v != nil {
			*v, err = d.DecodeUint64()
			return err
		}
		//case *bool:
		//	if v != nil {
		//		*v, err = d.DecodeBool()
		//		return err
		//	}
		//case *float32:
		//	if v != nil {
		//		*v, err = d.DecodeFloat32()
		//		return err
		//	}
		//case *float64:
		//	if v != nil {
		//		*v, err = d.DecodeFloat64()
		//		return err
		//	}
		//case *[]string:
		//	return ErrUnsupported // d.decodeStringSlicePtr(v)
		//case *map[string]string:
		//	return ErrUnsupported // d.decodeMapStringStringPtr(v)
		//case *map[string]interface{}:
		//	return ErrUnsupported // d.decodeMapStringInterfacePtr(v)
		//case *time.Duration:
		//	if v != nil {
		//		vv, err := d.DecodeInt64()
		//		*v = time.Duration(vv)
		//		return err
		//	}
		//case *time.Time:
		//	return ErrUnsupported
		//	/*if v != nil {
		//		*v, err = d.DecodeTime()
		//		return err
		//	}*/
	}
	//
	//vv := reflect.ValueOf(v)
	//if !vv.IsValid() {
	//	return errors.New("igbinary: Decode(nil)")
	//}
	//if vv.Kind() != reflect.Ptr {
	//	return fmt.Errorf("igbinary: Decode(non-pointer %T)", v)
	//}
	//if vv.IsNil() {
	//	return fmt.Errorf("igbinary: Decode(non-settable %T)", v)
	//}
	//
	//vv = vv.Elem()
	//if vv.Kind() == reflect.Interface {
	//	if !vv.IsNil() {
	//		vv = vv.Elem()
	//		if vv.Kind() != reflect.Ptr {
	//			return fmt.Errorf("igbinary: Decode(non-pointer %s)", vv.Type().String())
	//		}
	//	}
	//}

	return fmt.Errorf(`igbinary: Decode(unsupported type %s)`, reflect.TypeOf(v))
	// return d.DecodeValue(vv)
}

func (d *Decoder) PeekCode() (byte, error) {
	c, err := d.s.ReadByte()
	if err != nil {
		return 0, err
	}
	return c, d.s.UnreadByte()
}

func (d *Decoder) hasNilCode() bool {
	code, err := d.PeekCode()
	return err == nil && code == igcode.Nil
}

func (d *Decoder) DecodeNil() error {
	return d.skipExpected('N', ';')
}

func (d *Decoder) skipExpected(expected ...byte) error {
	for _, e := range expected {
		c, err := d.s.ReadByte()
		if err != nil {
			return err
		}
		if c != e {
			return fmt.Errorf(`igbinary: Decode(expected byte '%c' found '%c')`, e, c)
		}
	}
	return nil
}

func (d *Decoder) readCode() (byte, error) {
	c, err := d.s.ReadByte()
	if err != nil {
		return 0, err
	}
	if d.rec != nil {
		d.rec = append(d.rec, c)
	}
	return c, nil
}

func (d *Decoder) readN(n int) ([]byte, error) {
	var err error
	d.buf, err = readN(d.r, d.buf, n)
	if err != nil {
		return nil, err
	}
	if d.rec != nil {
		// TODO: read directly into d.rec?
		d.rec = append(d.rec, d.buf...)
	}
	return d.buf, nil
}

func (d *Decoder) decodeArrayLen() (int, error) {
	c, err := d.readCode()
	if err != nil {
		return 0, err
	}
	switch c {
	case igcode.Array8:
		v, err := d.uint8()
		return int(v), err
	}

	return 0, fmt.Errorf(`igbinary: Decode(array length code '%c')`, c)
}

// DisallowUnknownFields causes the Decoder to return an error when the destination
// is a struct and the input contains object keys which do not match any
// non-ignored, exported fields in the destination.
func (d *Decoder) DisallowUnknownFields(on bool) {
	if on {
		d.flags |= disallowUnknownFieldsFlag
	} else {
		d.flags &= ^disallowUnknownFieldsFlag
	}
}

package igbinary

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/zarken-go/igbinary/igcode"
	"io"
	"math"
	"reflect"
)

var headerBytes = []byte{0x00, 0x00, 0x00, 0x02}

type writer interface {
	io.Writer
	WriteByte(byte) error
}

type byteWriter struct {
	io.Writer
}

func newByteWriter(w io.Writer) byteWriter {
	return byteWriter{
		Writer: w,
	}
}

func (bw byteWriter) WriteByte(c byte) error {
	_, err := bw.Write([]byte{c})
	return err
}

// Marshal returns the igbinary serialized encoding of v.
func Marshal(v interface{}) ([]byte, error) {
	//enc := GetEncoder()
	enc := NewEncoder(nil)

	var buf bytes.Buffer
	//enc.Reset(&buf)
	enc.resetWriter(&buf)

	if err := enc.EncodeHeader(); err != nil {
		return nil, err
	}

	err := enc.Encode(v)
	b := buf.Bytes()

	// PutEncoder(enc)

	if err != nil {
		return nil, err
	}
	return b, err
}

type Encoder struct {
	w        writer
	buf      []byte
	strings  map[string]uint
	stringID uint
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	e := &Encoder{
		buf: make([]byte, 9),
	}
	e.resetWriter(w)
	return e
}

func (e *Encoder) resetWriter(w io.Writer) {
	if bw, ok := w.(writer); ok {
		e.w = bw
	} else {
		e.w = newByteWriter(w)
	}
}

func (e *Encoder) EncodeHeader() error {
	return e.write(headerBytes)
}

func (e *Encoder) Encode(v interface{}) error {
	switch v := v.(type) {
	case nil:
		return e.EncodeNil()
	case string:
		return e.EncodeString(v)
	case []byte:
		return e.EncodeBytes(v)
	case int:
		return e.EncodeInt64(int64(v))
	case int8:
		return e.EncodeInt64(int64(v))
	case int16:
		return e.EncodeInt64(int64(v))
	case int32:
		return e.EncodeInt64(int64(v))
	case int64:
	//	return e.EncodeInt64(v)
	//case uint:
	//	return e.EncodeUint(uint64(v))
	//case uint64:
	//	return e.encodeUint64Cond(v)
	case bool:
		return e.EncodeBool(v)
	//case float32:
	//	return e.EncodeFloat64(float64(v))
	case float64:
		return e.EncodeFloat64(v)
		//	case time.Duration:
		//		return e.encodeInt64Cond(int64(v))
		//	case time.Time:
		//		return e.EncodeTime(v)
	}

	return fmt.Errorf("igbinary: Encode(unsupported %s)", reflect.TypeOf(v))

	// return e.EncodeValue(reflect.ValueOf(v))
}

func (e *Encoder) EncodeNil() error {
	return e.writeBytes(igcode.Nil)
}

func (e *Encoder) EncodeBool(v bool) error {
	if v {
		return e.writeBytes(igcode.BoolTrue)
	}
	return e.writeBytes(igcode.BoolFalse)
}

func (e *Encoder) EncodeInt64(v int64) error {
	if v >= 0x00 && v <= 0xff {
		return e.writeBytes(igcode.PosInt8, byte(v))
	} else if v >= -0xff && v < 0x00 {
		return e.writeBytes(igcode.NegInt8, byte(v*-1))
	} else if v > 0xff && v <= 0xffff {
		return e.write2(igcode.PosInt16, uint16(v))
	} else if v < -0xff && v >= -0xffff {
		return e.write2(igcode.NegInt16, uint16(v*-1))
	} else if v > 0xffff && v <= 0xffffffff {
		return e.write4(igcode.PosInt32, uint32(v))
	} else if v < -0xffff && v >= -0xffffffff {
		return e.write4(igcode.NegInt32, uint32(v*-1))
	}

	return errors.New(`igbinary: Encode(int out of range)`)
}

func (e *Encoder) EncodeFloat64(v float64) error {
	return e.write8(igcode.Double, math.Float64bits(v))
}

func (e *Encoder) EncodeString(v string) error {
	return e.EncodeBytes([]byte(v))
}

func (e *Encoder) EncodeBytes(v []byte) error {
	// TODO: unnecessary re-conversion back to string.
	str := string(v)
	if e.strings == nil {
		e.strings = make(map[string]uint)
	}
	if ID, ok := e.strings[str]; ok {
		// Encode as ID
		if ID <= 0xff {
			return e.write1(igcode.StringID8, uint8(ID))
		}
		if ID <= 0xffff {
			return e.write2(igcode.StringID16, uint16(ID))
		}
		if ID <= 0xffffffff {
			return e.write4(igcode.StringID32, uint32(ID))
		}
		return fmt.Errorf(`igbinary: Encode(string ID exceeds range)`)
	}

	e.strings[str] = e.stringID
	e.stringID++

	length := len(v)
	if length <= 0xff {
		if err := e.write1(igcode.String8, uint8(length)); err != nil {
			return err
		}
		return e.write(v)
	}
	if length <= 0xffff {
		if err := e.write2(igcode.String16, uint16(length)); err != nil {
			return err
		}
		return e.write(v)
	}
	if length <= 0xffffffff {
		if err := e.write4(igcode.String32, uint32(length)); err != nil {
			return err
		}
		return e.write(v)
	}

	return errors.New(`igbinary: Encode([]byte exceeds capacity)`)
}

func (e *Encoder) write1(code byte, n uint8) error {
	e.buf = e.buf[:2]
	e.buf[0] = code
	e.buf[1] = n
	return e.write(e.buf)
}

func (e *Encoder) write2(code byte, n uint16) error {
	e.buf = e.buf[:3]
	e.buf[0] = code
	e.buf[1] = byte(n >> 8)
	e.buf[2] = byte(n)
	return e.write(e.buf)
}

func (e *Encoder) write4(code byte, n uint32) error {
	e.buf = e.buf[:5]
	e.buf[0] = code
	e.buf[1] = byte(n >> 24)
	e.buf[2] = byte(n >> 16)
	e.buf[3] = byte(n >> 8)
	e.buf[4] = byte(n)
	return e.write(e.buf)
}

func (e *Encoder) write8(code byte, n uint64) error {
	e.buf = e.buf[:9]
	e.buf[0] = code
	e.buf[1] = byte(n >> 56)
	e.buf[2] = byte(n >> 48)
	e.buf[3] = byte(n >> 40)
	e.buf[4] = byte(n >> 32)
	e.buf[5] = byte(n >> 24)
	e.buf[6] = byte(n >> 16)
	e.buf[7] = byte(n >> 8)
	e.buf[8] = byte(n)
	return e.write(e.buf)
}

func (e *Encoder) write(b []byte) error {
	_, err := e.w.Write(b)
	return err
}

func (e *Encoder) writeBytes(b ...byte) error {
	_, err := e.w.Write(b)
	return err
}

func (e *Encoder) EncodeArrayLen(length int) error {
	if length <= 0xff {
		return e.write1(igcode.Array8, uint8(length))
	}
	if length <= 0xffff {
		return e.write2(igcode.Array16, uint16(length))
	}
	if length <= 0xffffffff {
		return e.write4(igcode.Array32, uint32(length))
	}

	return fmt.Errorf(`igbinary: Encode(unsupported array length %d)`, length)
}

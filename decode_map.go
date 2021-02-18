package igbinary

import "reflect"

var (
	mapStringStringPtrType = reflect.TypeOf((*map[string]string)(nil))
)

func decodeMapValue(d *Decoder, v reflect.Value) error {
	n, err := d.DecodeArrayLen()
	if err != nil {
		return err
	}

	typ := v.Type()
	if n == -1 {
		v.Set(reflect.Zero(typ))
		return nil
	}

	if v.IsNil() {
		v.Set(reflect.MakeMap(typ))
	}
	if n == 0 {
		return nil
	}

	if err := d.decodeTypedMapValue(v, n); err != nil {
		return err
	}

	return nil
}

func (d *Decoder) decodeTypedMapValue(v reflect.Value, n int) error {
	typ := v.Type()
	keyType := typ.Key()
	valueType := typ.Elem()

	for i := 0; i < n; i++ {
		mk := reflect.New(keyType).Elem()
		if err := d.DecodeValue(mk); err != nil {
			return err
		}

		mv := reflect.New(valueType).Elem()
		if err := d.DecodeValue(mv); err != nil {
			return err
		}

		v.SetMapIndex(mk, mv)
	}

	return nil
}

func (d *Decoder) decodeMapStringStringPtr(ptr *map[string]string) error {
	size, err := d.DecodeArrayLen()
	if err != nil {
		return err
	}
	if size == -1 {
		*ptr = nil
		return nil
	}

	m := *ptr
	if m == nil {
		*ptr = make(map[string]string, min(size, maxMapSize))
		m = *ptr
	}

	for i := 0; i < size; i++ {
		mk, err := d.DecodeString()
		if err != nil {
			return err
		}
		mv, err := d.DecodeString()
		if err != nil {
			return err
		}
		m[mk] = mv
	}

	return nil
}

func decodeMapStringStringValue(d *Decoder, v reflect.Value) error {
	mptr := v.Addr().Convert(mapStringStringPtrType).Interface().(*map[string]string)
	return d.decodeMapStringStringPtr(mptr)
}

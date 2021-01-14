package igbinary

import (
	"reflect"
)

var (
	stringType = reflect.TypeOf((*string)(nil)).Elem()
)

func getDecoder(typ reflect.Type) decoderFunc {
	if v, ok := typeDecMap.Load(typ); ok {
		return v.(decoderFunc)
	}
	fn := _getDecoder(typ)
	typeDecMap.Store(typ, fn)
	return fn
}

func _getDecoder(typ reflect.Type) decoderFunc {
	kind := typ.Kind()

	if kind == reflect.Ptr {
		if _, ok := typeDecMap.Load(typ.Elem()); ok {
			return ptrDecoderFunc(typ)
		}
	}

	/*
		if typ.Implements(customDecoderType) {
			return decodeCustomValue
		}
		if typ.Implements(unmarshalerType) {
			return unmarshalValue
		}
		if typ.Implements(binaryUnmarshalerType) {
			return unmarshalBinaryValue
		}
		if typ.Implements(textUnmarshalerType) {
			return unmarshalTextValue
		}

		// Addressable struct field value.
		if kind != reflect.Ptr {
			ptr := reflect.PtrTo(typ)
			if ptr.Implements(customDecoderType) {
				return decodeCustomValueAddr
			}
			if ptr.Implements(unmarshalerType) {
				return unmarshalValueAddr
			}
			if ptr.Implements(binaryUnmarshalerType) {
				return unmarshalBinaryValueAddr
			}
			if ptr.Implements(textUnmarshalerType) {
				return unmarshalTextValueAddr
			}
		}
	*/

	switch kind {
	case reflect.Ptr:
		return ptrDecoderFunc(typ)
	//case reflect.Slice:
	//	elem := typ.Elem()
	//	if elem.Kind() == reflect.Uint8 {
	//		return decodeBytesValue
	//	}
	//	if elem == stringType {
	//		return decodeStringSliceValue
	//	}
	//case reflect.Array:
	//	if typ.Elem().Kind() == reflect.Uint8 {
	//		return decodeByteArrayValue
	//	}
	case reflect.Map:
		if typ.Key() == stringType {
			switch typ.Elem() {
			case stringType:
				return decodeMapStringStringValue
				//case interfaceType:
				//	return decodeMapStringInterfaceValue
			}
		}
	}

	return valueDecoders[kind]
}

func ptrDecoderFunc(typ reflect.Type) decoderFunc {
	decoder := getDecoder(typ.Elem())
	return func(d *Decoder, v reflect.Value) error {
		if d.hasNilCode() {
			if !v.IsNil() {
				v.Set(reflect.Zero(v.Type()))
			}
			return d.DecodeNil()
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		return decoder(d, v.Elem())
	}
}

var valueDecoders []decoderFunc

//nolint:gochecknoinits
func init() {
	valueDecoders = []decoderFunc{
		reflect.Bool:          decodeUnsupportedValue, // decodeBoolValue,
		reflect.Int:           decodeIntValue,
		reflect.Int8:          decodeInt8Value,
		reflect.Int16:         decodeInt16Value,
		reflect.Int32:         decodeInt32Value,
		reflect.Int64:         decodeInt64Value,
		reflect.Uint:          decodeUintValue,
		reflect.Uint8:         decodeUint8Value,
		reflect.Uint16:        decodeUint16Value,
		reflect.Uint32:        decodeUint32Value,
		reflect.Uint64:        decodeUint64Value,
		reflect.Float32:       decodeUnsupportedValue, // decodeFloat32Value,
		reflect.Float64:       decodeUnsupportedValue, // decodeFloat64Value,
		reflect.Complex64:     decodeUnsupportedValue,
		reflect.Complex128:    decodeUnsupportedValue,
		reflect.Array:         decodeUnsupportedValue, //   decodeArrayValue,
		reflect.Chan:          decodeUnsupportedValue,
		reflect.Func:          decodeUnsupportedValue,
		reflect.Interface:     decodeUnsupportedValue, //decodeInterfaceValue,
		reflect.Map:           decodeMapValue,
		reflect.Ptr:           decodeUnsupportedValue,
		reflect.Slice:         decodeUnsupportedValue, // decodeSliceValue,
		reflect.String:        decodeStringValue,
		reflect.Struct:        decodeStructValue,
		reflect.UnsafePointer: decodeUnsupportedValue,
	}
}

func decodeUnsupportedValue(_ *Decoder, v reflect.Value) error {
	return decodeErrorF(`unsupported %s`, v.Type())
}

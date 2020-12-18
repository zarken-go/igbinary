package igbinary

import (
	"fmt"
	"reflect"
)

var valueEncoders []encoderFunc

//nolint:gochecknoinits
func init() {
	valueEncoders = []encoderFunc{
		reflect.Bool:          encodeUnsupportedValue, // encodeBoolValue,
		reflect.Int:           encodeUnsupportedValue, // encodeIntValue,
		reflect.Int8:          encodeUnsupportedValue, // encodeIntValue,
		reflect.Int16:         encodeUnsupportedValue, // encodeIntValue,
		reflect.Int32:         encodeUnsupportedValue, // encodeIntValue,
		reflect.Int64:         encodeUnsupportedValue, // encodeIntValue,
		reflect.Uint:          encodeUnsupportedValue, // encodeUintValue,
		reflect.Uint8:         encodeUnsupportedValue, // encodeUintValue,
		reflect.Uint16:        encodeUnsupportedValue, // encodeUintValue,
		reflect.Uint32:        encodeUnsupportedValue, // encodeUintValue,
		reflect.Uint64:        encodeUnsupportedValue, // encodeUintValue,
		reflect.Float32:       encodeUnsupportedValue, // encodeFloat32Value,
		reflect.Float64:       encodeUnsupportedValue, // encodeFloat64Value,
		reflect.Complex64:     encodeUnsupportedValue,
		reflect.Complex128:    encodeUnsupportedValue,
		reflect.Array:         encodeUnsupportedValue, //encodeArrayValue,
		reflect.Chan:          encodeUnsupportedValue,
		reflect.Func:          encodeUnsupportedValue,
		reflect.Interface:     encodeUnsupportedValue, // encodeInterfaceValue,
		reflect.Map:           encodeUnsupportedValue, // encodeMapValue,
		reflect.Ptr:           encodeUnsupportedValue, // encodeUnsupportedValue,
		reflect.Slice:         encodeUnsupportedValue, // encodeSliceValue,
		reflect.String:        encodeUnsupportedValue, // encodeStringValue,
		reflect.Struct:        encodeUnsupportedValue, // encodeStructValue,
		reflect.UnsafePointer: encodeUnsupportedValue,
	}
}

func getEncoder(typ reflect.Type) encoderFunc {
	if v, ok := typeEncMap.Load(typ); ok {
		return v.(encoderFunc)
	}
	fn := _getEncoder(typ)
	typeEncMap.Store(typ, fn)
	return fn
}

func _getEncoder(typ reflect.Type) encoderFunc {
	kind := typ.Kind()

	//if kind == reflect.Ptr {
	//	if _, ok := typeEncMap.Load(typ.Elem()); ok {
	//		return ptrEncoderFunc(typ)
	//	}
	//}

	/*if typ.Implements(customEncoderType) {
		return encodeCustomValue
	}
	if typ.Implements(marshalerType) {
		return marshalValue
	}
	if typ.Implements(binaryMarshalerType) {
		return marshalBinaryValue
	}
	if typ.Implements(textMarshalerType) {
		return marshalTextValue
	}

	// Addressable struct field value.
	if kind != reflect.Ptr {
		ptr := reflect.PtrTo(typ)
		if ptr.Implements(customEncoderType) {
			return encodeCustomValuePtr
		}
		if ptr.Implements(marshalerType) {
			return marshalValuePtr
		}
		if ptr.Implements(binaryMarshalerType) {
			return marshalBinaryValueAddr
		}
		if ptr.Implements(textMarshalerType) {
			return marshalTextValueAddr
		}
	}*/

	/*if typ == errorType {
		return encodeErrorValue
	}*/

	//switch kind {
	//case reflect.Ptr:
	//	return ptrEncoderFunc(typ)
	//	/*case reflect.Slice:
	//		elem := typ.Elem()
	//		if elem.Kind() == reflect.Uint8 {
	//			return encodeByteSliceValue
	//		}
	//		if elem == stringType {
	//			return encodeStringSliceValue
	//		}
	//	case reflect.Array:
	//		if typ.Elem().Kind() == reflect.Uint8 {
	//			return encodeByteArrayValue
	//		}
	//	case reflect.Map:
	//		if typ.Key() == stringType {
	//			switch typ.Elem() {
	//			case stringType:
	//				return encodeMapStringStringValue
	//			case interfaceType:
	//				return encodeMapStringInterfaceValue
	//			}
	//		}*/
	//}

	return valueEncoders[kind]
}

func encodeUnsupportedValue(e *Encoder, v reflect.Value) error {
	return fmt.Errorf("igbinary: Encode(unsupported %s)", v.Type())
}

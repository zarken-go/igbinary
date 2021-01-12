package igcode

const (
	Nil byte = iota
	ArrayRef8
	ArrayRef16
	ArrayRef32
	BoolFalse
	BoolTrue
	PosInt8     // Integer 8-bit Positive
	NegInt8     // Integer 8-bit Negative
	PosInt16    // Integer 16-bit Positive
	NegInt16    // Integer 16-bit Negative
	PosInt32    // Integer 32-bit Positive
	NegInt32    // Integer 32-bit Negative
	Double      // Double
	StringEmpty // Empty String
	StringID8   // String ID
	StringID16  // String ID
	StringID32  // String ID
	String8     // String
	String16    // String
	String32    // String
	Array8      // Array
	Array16     // Array
	Array32     // Array
	Object8     // Object
	Object16    // Object
	Object32    // Object
	ObjectID8   // Object string id
	ObjectID16  // Object string id
	ObjectID32  // Object string id
	ObjectSer8  // Object serialized data
	ObjectSer16 // Object serialized data
	ObjectSer32 // Object serialized data
	PosInt64    // Integer 64-bit Positive
	NegInt64    // Integer 64-bit Positive
	ObjectRef8  // Object reference
	ObjectRef16 // Object reference
	ObjectRef32 // Object reference
	SimpleRef   // Simple reference
)

func IsNegative(c byte) bool {
	switch c {
	case NegInt8, NegInt16, NegInt32, NegInt64:
		return true
	default:
		return false
	}
}

func IsStringID(c byte) bool {
	switch c {
	case StringID8, StringID16, StringID32:
		return true
	default:
		return false
	}
}

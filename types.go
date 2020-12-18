package igbinary

import (
	"fmt"
	"github.com/vmihailenco/tagparser"
	"reflect"
	"sync"
)

type (
	encoderFunc func(*Encoder, reflect.Value) error
	decoderFunc func(*Decoder, reflect.Value) error
)

var (
	typeEncMap sync.Map
	typeDecMap sync.Map
)
var structs = newStructCache()

type structCache struct {
	m sync.Map
}

type structCacheKey struct {
	tag string
	typ reflect.Type
}

func newStructCache() *structCache {
	return new(structCache)
}

func (m *structCache) Fields(typ reflect.Type, tag string) *fields {
	key := structCacheKey{tag: tag, typ: typ}

	if v, ok := m.m.Load(key); ok {
		return v.(*fields)
	}

	fs := getFields(typ, tag)
	m.m.Store(key, fs)

	return fs
}

type field struct {
	name  string
	index []int
	// omitEmpty bool
	encoder encoderFunc
	decoder decoderFunc
}

func newFields(typ reflect.Type) *fields {
	return &fields{
		Type: typ,
		Map:  make(map[string]*field, typ.NumField()),
		List: make([]*field, 0, typ.NumField()),
	}
}

var (
	defaultStructTag = `php`
)

func getFields(typ reflect.Type, fallbackTag string) *fields {
	fs := newFields(typ)

	// var omitEmpty bool
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)

		tagStr := f.Tag.Get(defaultStructTag)
		if tagStr == "" && fallbackTag != "" {
			tagStr = f.Tag.Get(fallbackTag)
		}

		tag := tagparser.Parse(tagStr)
		if tag.Name == "-" {
			continue
		}

		field := &field{
			name:  tag.Name,
			index: f.Index,
			// omitEmpty: omitEmpty || tag.HasOption("omitempty"),
		}

		field.encoder = getEncoder(f.Type)
		field.decoder = getDecoder(f.Type)

		if field.name == "" {
			field.name = f.Name
		}

		fs.Add(field)
	}

	return fs
}

type fields struct {
	Type reflect.Type
	Map  map[string]*field
	List []*field
	// AsArray bool

	// hasOmitEmpty bool
}

func (fs *fields) Add(field *field) {
	// fs.warnIfFieldExists(field.name)
	fs.Map[field.name] = field
	fs.List = append(fs.List, field)
	//if field.omitEmpty {
	//	fs.hasOmitEmpty = true
	//}
}

func (f *field) DecodeValue(d *Decoder, strct reflect.Value) error {
	v := fieldByIndexAlloc(strct, f.index)
	if f.decoder == nil {
		return fmt.Errorf(`igbinary: could not find decoder for field %s`, f.name)
	}
	return f.decoder(d, v)
}

func fieldByIndex(v reflect.Value, index []int) (_ reflect.Value, ok bool) {
	if len(index) == 1 {
		return v.Field(index[0]), true
	}

	for i, idx := range index {
		if i > 0 {
			if v.Kind() == reflect.Ptr {
				if v.IsNil() {
					return v, false
				}
				v = v.Elem()
			}
		}
		v = v.Field(idx)
	}

	return v, true
}

func fieldByIndexAlloc(v reflect.Value, index []int) reflect.Value {
	if len(index) == 1 {
		return v.Field(index[0])
	}

	/*
		for i, idx := range index {
			if i > 0 {
				var ok bool
				v, ok = indirectNil(v)
				if !ok {
					return v
				}
			}
			v = v.Field(idx)
		}

		return v
	*/

	panic(`unsupported`)
}

func (f *field) EncodeValue(e *Encoder, strct reflect.Value) error {
	v, ok := fieldByIndex(strct, f.index)
	if !ok {
		return e.EncodeNil()
	}
	return f.encoder(e, v)
}

func (fs *fields) OmitEmpty(strct reflect.Value) []*field {
	//if !fs.hasOmitEmpty {
	return fs.List
	//}

	/*fields := make([]*field, 0, len(fs.List))

	for _, f := range fs.List {
		if !f.Omit(strct) {
			fields = append(fields, f)
		}
	}

	return fields
	*/
}

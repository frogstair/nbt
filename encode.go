package nbt

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"reflect"
)

// Encode encodes a given value into its corresponding nbt representation
// Maps and structs are interpreted as compounds
func Encode(v interface{}, name string) []byte {
	var b bytes.Buffer

	m, ok := v.(map[string]interface{})
	if ok {
		writeMap(&b, m, name)
	} else {
		panic("Cannot decode struct, not implemented")
	}

	return b.Bytes()
}

func EncodeCompress(v interface{}, name string) []byte {
	var b bytes.Buffer

	m, ok := v.(map[string]interface{})
	if ok {
		writeMap(&b, m, name)
	} else {
		panic("Cannot decode struct, not implemented")
	}

	var bc bytes.Buffer
	gz := gzip.NewWriter(&bc)
	gz.Write(b.Bytes())
	gz.Close()

	return bc.Bytes()
}

func writeMap(b *bytes.Buffer, v map[string]interface{}, name string) {

	b.WriteByte(TagCompound)
	if name == "" {
		b.Write([]byte{0, 0})
	} else {
		writeName(b, name)
	}

	for k, el := range v {
		writeType(b, el, k)
	}

	b.WriteByte(TagEnd)
}

func writeByte(b *bytes.Buffer, s byte, name string) {
	b.WriteByte(TagByte)
	writeName(b, name)
	b.WriteByte(s)
}

func writeShort(b *bytes.Buffer, s int16, name string) {
	b.WriteByte(TagShort)
	writeName(b, name)
	binary.Write(b, binary.BigEndian, s)
}

func writeInt(b *bytes.Buffer, s int32, name string) {
	b.WriteByte(TagInt)
	writeName(b, name)
	binary.Write(b, binary.BigEndian, s)
}

func writeLong(b *bytes.Buffer, s int64, name string) {
	b.WriteByte(TagLong)
	writeName(b, name)
	binary.Write(b, binary.BigEndian, s)
}

func writeFloat(b *bytes.Buffer, s float32, name string) {
	b.WriteByte(TagFloat)
	writeName(b, name)
	binary.Write(b, binary.BigEndian, s)
}

func writeDouble(b *bytes.Buffer, s float64, name string) {
	b.WriteByte(TagDouble)
	writeName(b, name)
	binary.Write(b, binary.BigEndian, s)
}

func writeByteArray(b *bytes.Buffer, s []byte, name string) {
	b.WriteByte(TagByteArray)
	writeName(b, name)

	l := len(s)
	binary.Write(b, binary.BigEndian, int32(l))
	binary.Write(b, binary.BigEndian, s)
}

func writeString(b *bytes.Buffer, s string, name string) {
	b.WriteByte(TagString)
	writeName(b, name)
	writeName(b, s)
}

func writeList(b *bytes.Buffer, s []interface{}, name string) {
	writeName(b, name)
	if len(s) == 0 {
		b.Write([]byte{0, 0, 0, 0, 0})
		return
	}

	t := getType(s[0])
	b.WriteByte(t)

	binary.Write(b, binary.BigEndian, int32(len(s)))

	for _, el := range s {
		if t == TagString {
			writeName(b, el.(string))
		} else if t == TagByteArray {
			e := el.([]byte)
			l := len(e)
			binary.Write(b, binary.BigEndian, int32(l))
			binary.Write(b, binary.BigEndian, e)
		} else if t == TagCompound {
			e := el.(map[string]interface{})
			for k, el := range e {
				writeType(b, el, k)
			}
			b.WriteByte(TagEnd)
		} else if t == TagList {
			e := reflect.ValueOf(el)
			arr := make([]interface{}, e.Len())
			for i := 0; i < e.Len(); i++ {
				arr[i] = e.Index(i).Interface()
			}
			writeList(b, arr, "")
		}
		binary.Write(b, binary.BigEndian, el)
	}
}

func writeName(b *bytes.Buffer, s string) error {
	if s == "" {
		return nil
	}
	l := int16(len(s))
	err := binary.Write(b, binary.BigEndian, l)
	if err != nil {
		return err
	}
	_, err = b.WriteString(s)
	if err != nil {
		return err
	}
	return nil
}

func writeType(b *bytes.Buffer, el interface{}, name string) {
	t := getType(el)
	switch t {
	case TagByte:
		writeByte(b, el.(byte), name)
	case TagShort:
		writeShort(b, el.(int16), name)
	case TagInt:
		writeInt(b, el.(int32), name)
	case TagLong:
		writeLong(b, el.(int64), name)
	case TagFloat:
		writeFloat(b, el.(float32), name)
	case TagDouble:
		writeDouble(b, el.(float64), name)
	case TagByteArray:
		writeByteArray(b, el.([]byte), name)
	case TagString:
		writeString(b, el.(string), name)
	case TagList:
		b.WriteByte(TagList)
		s := reflect.ValueOf(el)
		arr := make([]interface{}, s.Len())
		for i := 0; i < s.Len(); i++ {
			arr[i] = s.Index(i).Interface()
		}
		writeList(b, arr, name)
	case TagCompound:
		writeMap(b, el.(map[string]interface{}), name)
	default:
		panic(fmt.Errorf("invalid type supplied %T", el))
	}
}

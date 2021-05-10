package nbt

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"reflect"
)

// Encode encodes a given value into its corresponding nbt representation
func Encode(v interface{}, name string) []byte {
	var b bytes.Buffer

	m, okc := v.(C)
	var mm map[string]interface{}
	if okc {
		mm = map[string]interface{}(m)
	} else {
		mm = v.(map[string]interface{})
	}

	writeMap(&b, mm, name)

	if mm == nil {
		panic("Cannot decode struct, not implemented")
	}

	return b.Bytes()
}

// EncodeCompress encodes a given value into its corresponding nbt representation
// and used gzip compression on it afterwards
func EncodeCompress(v interface{}, name string) []byte {
	var b bytes.Buffer

	m, okc := v.(C)
	var mm map[string]interface{} = nil
	if okc {
		mm = map[string]interface{}(m)
	} else {
		mm = v.(map[string]interface{})
	}

	writeMap(&b, mm, name)

	if mm == nil {
		panic("Cannot decode struct, not implemented")
	}

	var bc bytes.Buffer
	gz := gzip.NewWriter(&bc)
	gz.Write(b.Bytes())
	gz.Close()

	return bc.Bytes()
}

func writeMap(b *bytes.Buffer, v map[string]interface{}, name string) error {

	b.WriteByte(tagCompound)
	if name == "" {
		b.Write([]byte{0, 0})
	} else {
		writeName(b, name)
	}

	for k, el := range v {
		err := writeType(b, el, k)
		if err != nil {
			return err
		}
	}

	b.WriteByte(tagEnd)
	return nil
}

func writeByte(b *bytes.Buffer, s int8, name string) {
	b.WriteByte(tagByte)
	writeName(b, name)
	binary.Write(b, binary.BigEndian, s)
}

func writeShort(b *bytes.Buffer, s int16, name string) {
	b.WriteByte(tagShort)
	writeName(b, name)
	binary.Write(b, binary.BigEndian, s)
}

func writeInt(b *bytes.Buffer, s int32, name string) {
	b.WriteByte(tagInt)
	writeName(b, name)
	binary.Write(b, binary.BigEndian, s)
}

func writeLong(b *bytes.Buffer, s int64, name string) {
	b.WriteByte(tagLong)
	writeName(b, name)
	binary.Write(b, binary.BigEndian, s)
}

func writeFloat(b *bytes.Buffer, s float32, name string) {
	b.WriteByte(tagFloat)
	writeName(b, name)
	binary.Write(b, binary.BigEndian, s)
}

func writeDouble(b *bytes.Buffer, s float64, name string) {
	b.WriteByte(tagDouble)
	writeName(b, name)
	binary.Write(b, binary.BigEndian, s)
}

func writeByteArray(b *bytes.Buffer, s []int8, name string) {
	b.WriteByte(tagByteArray)
	writeName(b, name)

	l := len(s)
	binary.Write(b, binary.BigEndian, int32(l))
	binary.Write(b, binary.BigEndian, s)
}

func writeString(b *bytes.Buffer, s string, name string) {
	b.WriteByte(tagString)
	writeName(b, name)
	writeName(b, s)
}

func writeList(b *bytes.Buffer, s []interface{}, name string) {
	writeName(b, name)
	if len(s) == 0 {
		b.Write([]byte{0, 0, 0, 0, 0})
		return
	}

	t, _ := getType(s[0])
	b.WriteByte(t)

	binary.Write(b, binary.BigEndian, int32(len(s)))

	for _, el := range s {
		if t == tagString {
			writeName(b, el.(string))
		} else if t == tagByteArray {
			e := el.([]int8)
			l := len(e)
			binary.Write(b, binary.BigEndian, int32(l))
			binary.Write(b, binary.BigEndian, e)
		} else if t == tagCompound {
			m, okc := el.(C)
			var mm map[string]interface{}
			if okc {
				mm = map[string]interface{}(m)
			} else {
				mm = el.(map[string]interface{})
			}
			for k, el := range mm {
				writeType(b, el, k)
			}
			b.WriteByte(tagEnd)
		} else if t == tagList {
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

func writeType(b *bytes.Buffer, el interface{}, name string) error {
	t, g := getType(el)
	switch t {
	case tagByte:
		if g {
			writeByte(b, int8(el.(int)), name)
		} else {
			writeByte(b, el.(int8), name)
		}
	case tagShort:
		if g {
			writeShort(b, int16(el.(int)), name)
		} else {
			writeShort(b, el.(int16), name)
		}
	case tagInt:
		if g {
			writeInt(b, int32(el.(int)), name)
		} else {
			writeInt(b, el.(int32), name)
		}
	case tagLong:
		if g {
			writeLong(b, int64(el.(int)), name)
		} else {
			writeLong(b, el.(int64), name)
		}
	case tagFloat:
		writeFloat(b, el.(float32), name)
	case tagDouble:
		writeDouble(b, el.(float64), name)
	case tagByteArray:
		writeByteArray(b, el.([]int8), name)
	case tagString:
		writeString(b, el.(string), name)
	case tagList:
		b.WriteByte(tagList)
		s := reflect.ValueOf(el)
		arr := make([]interface{}, s.Len())
		for i := 0; i < s.Len(); i++ {
			arr[i] = s.Index(i).Interface()
		}
		writeList(b, arr, name)
	case tagCompound:
		m, okc := el.(C)
		var mm map[string]interface{}
		if okc {
			mm = map[string]interface{}(m)
		} else {
			mm = el.(map[string]interface{})
		}
		writeMap(b, mm, name)
	default:
		return fmt.Errorf("invalid type supplied %T", el)
	}
	return nil
}

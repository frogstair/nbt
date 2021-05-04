package nbt

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

var errSyntax = fmt.Errorf("invalid byte sequence")

func DecodeBytes(data []byte, v interface{}) (err error) {
	b := bufio.NewReader(bytes.NewBuffer(data))

	m, ok := v.(*C)
	if ok {
		var mp interface{}
		_, mp, _, err = readNamedNext(b)
		(*m) = mp.(C)
	} else {
		panic("Cannot decode struct, not implemented")
	}

	return err
}

func DecodeCompressedBytes(data []byte, v interface{}) (err error) {
	gr, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer gr.Close()
	data, err = ioutil.ReadAll(gr)
	if err != nil {
		return err
	}

	return DecodeBytes(data, v)
}

func DecodeStream(r io.Reader, v interface{}) error {
	b := bufio.NewReader(r)

	var err error
	m, ok := v.(*C)
	if ok {
		var mp interface{}
		_, mp, _, err = readNamedNext(b)
		(*m) = mp.(C)
	} else {
		panic("Cannot decode struct, not implemented")
	}

	return err
}

func DecodeCompressedStream(r io.Reader, v interface{}) error {

	gr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gr.Close()
	data, err := ioutil.ReadAll(gr)
	if err != nil {
		return err
	}

	return DecodeBytes(data, v)
}

func readMap(r *bufio.Reader) (C, error) {

	m := make(C)

	for {
		n, val, typ, err := readNamedNext(r)
		if err != nil {
			return nil, err
		}
		if typ == tagEnd {
			return m, nil
		}
		m[n] = val
	}
}

func readString(r *bufio.Reader) (string, error) {
	length, err := readShort(r)
	if err != nil {
		return "", err
	}

	b := make([]byte, length)
	_, err = r.Read(b)
	if err != nil {
		return "", errSyntax
	}

	s := string(b)
	return s, nil
}

func readByte(r *bufio.Reader) (b byte, err error) {
	err = binary.Read(r, binary.BigEndian, &b)
	return
}

func readShort(r *bufio.Reader) (s int16, err error) {
	err = binary.Read(r, binary.BigEndian, &s)
	return
}

func readInt(r *bufio.Reader) (i int32, err error) {
	err = binary.Read(r, binary.BigEndian, &i)
	return
}

func readLong(r *bufio.Reader) (i int64, err error) {
	err = binary.Read(r, binary.BigEndian, &i)
	return
}

func readFloat(r *bufio.Reader) (f float32, err error) {
	err = binary.Read(r, binary.BigEndian, &f)
	return
}

func readDouble(r *bufio.Reader) (d float64, err error) {
	err = binary.Read(r, binary.BigEndian, &d)
	return
}

func readByteArray(r *bufio.Reader) (b []byte, err error) {
	l, err := readInt(r)
	if err != nil {
		return
	}

	b = make([]byte, l)
	err = binary.Read(r, binary.BigEndian, &b)
	return
}

func readList(r *bufio.Reader) (l []interface{}, err error) {
	t, err := r.ReadByte()
	if err != nil {
		return
	}

	len, err := readInt(r)
	if err != nil {
		return
	}

	l = make([]interface{}, len)

	for i := int32(0); i < len; i++ {
		l[i], err = readUnnamedNext(r, t)
		if err != nil {
			return
		}
	}

	return
}

func readNamedNext(r *bufio.Reader) (name string, v interface{}, t byte, err error) {
	t, err = r.ReadByte()
	if err != nil {
		err = errSyntax
		return
	}
	if t == 0 {
		v = nil
		err = nil
		name = ""
		return
	}

	name, err = readString(r)
	if err != nil {
		err = errSyntax
		return
	}

	v, err = readUnnamedNext(r, t)
	return
}

func readUnnamedNext(r *bufio.Reader, t byte) (v interface{}, err error) {
	switch t {
	case tagByte:
		v, err = readByte(r)
	case tagShort:
		v, err = readShort(r)
	case tagInt:
		v, err = readInt(r)
	case tagLong:
		v, err = readLong(r)
	case tagFloat:
		v, err = readFloat(r)
	case tagDouble:
		v, err = readDouble(r)
	case tagByteArray:
		v, err = readByteArray(r)
	case tagString:
		v, err = readString(r)
	case tagList:
		v, err = readList(r)
	case tagCompound:
		v, err = readMap(r)
	}

	if err != nil {
		err = errSyntax
	}

	return
}

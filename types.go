package nbt

import "reflect"

const (
	tagEnd       = 0
	tagByte      = 1
	tagShort     = 2
	tagInt       = 3
	tagLong      = 4
	tagFloat     = 5
	tagDouble    = 6
	tagByteArray = 7
	tagString    = 8
	tagList      = 9
	tagCompound  = 10
)

// C is a wrapper for map[string]interface{}
type C map[string]interface{}

func getType(v interface{}) (byte, bool) {
	switch v := v.(type) {
	case int:
		switch {
		case v < 127 && v > -128:
			return tagByte, true
		case v < 32767 && v > -32768:
			return tagShort, true
		case v < 2147483647 && v > -2147483648:
			return tagInt, true
		case v < 9223372036854775807 && v > -9223372036854775807:
			return tagLong, true
		default:
			return tagEnd, false
		}
	case int8:
		return tagByte, false
	case int16:
		return tagShort, false
	case int32:
		return tagInt, false
	case int64:
		return tagLong, false
	case float32:
		return tagFloat, false
	case float64:
		return tagDouble, false
	case []byte:
		return tagByteArray, false
	case string:
		return tagString, false
	case C:
		return tagCompound, false
	case map[string]interface{}:
		return tagCompound, false
	default:
		val := reflect.TypeOf(v)
		if val.Kind() == reflect.Slice {
			return tagList, false
		} else {
			return tagEnd, false
		}
	}
}

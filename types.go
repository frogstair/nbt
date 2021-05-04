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

func getType(v interface{}) byte {
	switch v := v.(type) {
	case byte:
		return tagByte
	case int16:
		return tagShort
	case int32:
		return tagInt
	case int64:
		return tagLong
	case float32:
		return tagFloat
	case float64:
		return tagDouble
	case []byte:
		return tagByteArray
	case string:
		return tagString
	case map[string]interface{}:
		return tagCompound
	default:
		val := reflect.TypeOf(v)
		if val.Kind() == reflect.Slice {
			return tagList
		} else {
			return tagEnd
		}
	}
}

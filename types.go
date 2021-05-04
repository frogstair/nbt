package nbt

import "reflect"

const (
	TagEnd       = 0
	TagByte      = 1
	TagShort     = 2
	TagInt       = 3
	TagLong      = 4
	TagFloat     = 5
	TagDouble    = 6
	TagByteArray = 7
	TagString    = 8
	TagList      = 9
	TagCompound  = 10
)

func getType(v interface{}) byte {
	switch v := v.(type) {
	case byte:
		return TagByte
	case int16:
		return TagShort
	case int32:
		return TagInt
	case int64:
		return TagLong
	case float32:
		return TagFloat
	case float64:
		return TagDouble
	case []byte:
		return TagByteArray
	case string:
		return TagString
	case map[string]interface{}:
		return TagCompound
	default:
		val := reflect.TypeOf(v)
		if val.Kind() == reflect.Slice {
			return TagList
		} else {
			return TagEnd
		}
	}
}

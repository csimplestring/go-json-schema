package schema
import (
	"encoding/json"
	"strings"
	"fmt"
	"reflect"
)

type JsonType string

const (
	JsonTypeInteger = JsonType("integer")
	JsonTypeNumber  = JsonType("number")
	JsonTypeString  = JsonType("string")
	JsonTypeArray   = JsonType("array")
	JsonTypeObject  = JsonType("object")
	JsonTypeBoolean = JsonType("boolean")
	JsonTypeNull    = JsonType("null")
)

func ParseType(v interface{}) (JsonType, error) {
	switch v.(type) {
	case bool:
		return JsonTypeBoolean, nil
	case int8, int16, int32, int, int64:
		return JsonTypeInteger, nil
	case float32, float64:
		return JsonTypeNumber, nil
	case string:
		return JsonTypeString, nil
	case []interface{}:
		return JsonTypeArray, nil
	case map[string]interface{}:
		return JsonTypeObject, nil
	case nil:
		return JsonTypeNull, nil
	case json.Number:
		if strings.Contains(v.(json.Number).String(), ".") {
			return JsonTypeNumber, nil
		} else {
			return JsonTypeInteger, nil
		}
	default:
		return JsonType(""), fmt.Errorf("Unsupported json type %s", reflect.TypeOf(v).String())
	}
}

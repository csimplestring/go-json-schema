package schema
import "encoding/json"

//
//type Schema struct {
//	ID string
//	Schema string
//	Title string
//	Description string
//}

// firstly, assume the schema passed in is valid, after the schema validation is finished, you can validate the schema
// against it.

type Schema map[string]interface{}

func (s Schema) getFloat64Value(key string) (value float64, exist bool) {
	value = 0
	exist = false

	if v, ok := s[key]; ok {
		value, _ = v.(json.Number).Float64()
		exist = true
	}

	return
}

func (s Schema) getBoolValue(key string) (value bool, exist bool) {
	value = false
	exist = false

	if v, ok := s[key]; ok {
		value, _ = v.(bool)
		exist = true
	}

	return
}

func (s Schema) getIntValue(key string) (value int, exist bool) {
	value = 0
	exist = false

	if v, ok := s[key]; ok {
		i64, _ := v.(json.Number).Int64()
		value = int(i64)
		exist = true
	}

	return
}

func (s Schema) getStringValue(key string) (value string, exist bool) {
	value = ""
	exist = false

	if v, ok := s[key]; ok {
		value, _ = v.(string)
		exist = true
	}

	return
}

// validation keywords for numeric

func (s Schema) MultipleOf() (divided float64, exist bool) {
	return s.getFloat64Value("multipleOf")
}

func (s Schema) Maximum() (max float64, exist bool) {
	return s.getFloat64Value("maximum")
}

func (s Schema) Minimum() (max float64, exist bool) {
	return s.getFloat64Value("minimum")
}

func (s Schema) ExclusiveMinimum() bool {
	if v, ok := s.getBoolValue("exclusiveMinimum"); ok {
		return v
	}
	return false
}

func (s Schema) ExclusiveMaximum() bool {
	if v, ok := s.getBoolValue("exclusiveMaximum"); ok {
		return v
	}
	return false
}

// validation keywords for string

func (s Schema) MaxLength() (maxLen int, exist bool) {
	return s.getIntValue("maxLength")
}

func (s Schema) MinLength() (minLen int, exist bool) {
	return s.getIntValue("minLength")
}

func (s Schema) Pattern() (pattern string, exist bool) {
	return s.getStringValue("pattern")
}
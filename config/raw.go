package config

import (
	"code.cloudfoundry.org/bytefmt"
	"fmt"
	"github.com/go-yaml/yaml"
	"io"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type Raw map[string]interface{}

func Parse(reader io.Reader) (Raw, error) {
	var out map[string]interface{}
	if err := yaml.NewDecoder(reader).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c Raw) Sub(key string) Raw {
	val := c[key]
	if val == nil {
		return nil
	}
	if reflect.TypeOf(val).Kind() == reflect.Map {
		switch v := val.(type) {
		case map[interface{}]interface{}:
			var sub = map[string]interface{}{}
			for key, elem := range v {
				if s, ok := key.(string); ok {
					sub[s] = elem
				}
			}
			return sub
		case map[string]interface{}:
			return v
		}
	}
	return nil
}

func (c Raw) Has(key string) bool {
	_, exists := c[key]
	return exists
}

func (c Raw) String(key string) string {
	return asString(c[key])
}

func (c Raw) StringSlice(key string) []string {
	val := c[key]
	if val == nil {
		return nil
	}
	if s, ok := val.([]string); ok {
		return s
	}
	if s, ok := val.([]interface{}); ok {
		slice := make([]string, 0, len(s))
		for _, raw := range s {
			if elem := asString(raw); elem != "" {
				slice = append(slice, elem)
			}
		}
		return slice
	}
	return nil
}

func (c Raw) Bool(key string) bool {
	return asBool(c[key])
}

func (c Raw) Uint64(key string) uint64 {
	return asUint64(c[key])
}

func (c Raw) Int64(key string) int64 {
	return asInt64(c[key])
}

func (c Raw) Bytes(key string) uint64 {
	val := c[key]
	if val == nil {
		return 0
	}
	s, ok := val.(string)
	if !ok {
		return asUint64(val)
	}

	s = strings.ToUpper(s)
	i := strings.IndexFunc(s, unicode.IsLetter)

	if i >= 0 {
		s = strings.Replace(s, " ", "", -1)
		bytes, err := bytefmt.ToBytes(s)
		if err == nil {
			return bytes
		}
	}
	parsed, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return parsed
}

func asString(val interface{}) string {
	if val == nil {
		return ""
	}
	if s, ok := val.(string); ok {
		return s
	}
	if s, ok := val.(fmt.Stringer); ok && s != nil {
		return s.String()
	}
	return fmt.Sprintf("%v", val)
}

func asUint64(val interface{}) uint64 {
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int8:
		return uint64(v)
	case int16:
		return uint64(v)
	case int32:
		return uint64(v)
	case int64:
		return uint64(v)
	case int:
		return uint64(v)
	case uint8:
		return uint64(v)
	case uint16:
		return uint64(v)
	case uint32:
		return uint64(v)
	case uint64:
		return v
	case uint:
		return uint64(v)
	case string:
		i, err := strconv.ParseUint(v, 10, 64)
		if err == nil {
			return i
		}
	}
	return 0
}

func asInt64(val interface{}) int64 {
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case int:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	case uint:
		return int64(v)
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return i
		}
	}
	return 0
}

func asBool(val interface{}) bool {
	if val == nil {
		return false
	}
	if b, ok := val.(bool); ok {
		return b
	}
	if s, ok := val.(string); ok {
		b, err := strconv.ParseBool(s)
		if err == nil {
			return b
		}
	}
	return false
}

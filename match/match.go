package match

import (
	"reflect"
	"regexp"
)

type Matcher func(interface{}) bool

func Match(v1, v2 interface{}) bool {
	matcher, isMatcher := v1.(Matcher)
	if isMatcher {
		if !matcher(v2) {
			return false
		}
	} else {
		if !reflect.DeepEqual(v1, v2) {
			return false
		}
	}
	return true
}

func Any() Matcher {
	return func(_ interface{}) bool {
		return true
	}
}

func Regexp(re *regexp.Regexp) Matcher {
	return func(val interface{}) bool {
		switch val := val.(type) {
		case []byte:
			return re.Match(val)
		case string:
			return re.MatchString(val)
		default:
			return false
		}
	}
}

func Keys(m1 map[interface{}]interface{}) Matcher {
	t1 := reflect.TypeOf(m1)
	v1 := reflect.ValueOf(m1)
	return func(m2 interface{}) bool {
		t2 := reflect.TypeOf(m2)
		if t2.Kind() != reflect.Map {
			return false
		}
		keyT := t1.Key()
		valT := t1.Elem()
		isMatcher := reflect.TypeOf(Any()) == valT
		if t2.Key() != keyT || (!isMatcher && t2.Elem() != valT) {
			return false
		}
		v2 := reflect.ValueOf(m2)
		for _, key := range v1.MapKeys() {
			i1 := v1.MapIndex(key).Interface()
			i2 := v2.MapIndex(key).Interface()
			if !Match(i1, i2) {
				return false
			}
		}
		return true
	}
}

func Contains(v interface{}) Matcher {
	return func(s interface{}) bool {
		if reflect.TypeOf(s).Kind() != reflect.Slice {
			return false
		}
		for _, sv := range s.([]interface{}) {
			if Match(v, sv) {
				return true
			}
		}
		return false
	}
}

func Fields(s1 interface{}) Matcher {
	t1 := reflect.TypeOf(s1)
	if t1.Kind() != reflect.Struct {
		panic("Fields matcher must be called with a struct")
	}
	v1 := reflect.ValueOf(s1)
	return func(s2 interface{}) bool {
		t2 := reflect.TypeOf(s2)
		if t2.Kind() != reflect.Struct {
			return false
		}
		v2 := reflect.ValueOf(s2)
		for i := 0; i < t1.NumField(); i++ {
			name := t1.Field(i).Name
			_, hasName := t2.FieldByName(name)
			if !hasName {
				return false
			}
			i1 := v1.FieldByName(name).Interface()
			i2 := v2.FieldByName(name).Interface()
			if !Match(i1, i2) {
				return false
			}
		}
		return true
	}
}

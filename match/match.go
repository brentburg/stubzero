package match

import (
	"reflect"
	"regexp"
)

type Matcher func(interface{}) bool

func Match(v1, v2 interface{}) bool {
	if reflect.TypeOf(v1) == reflect.TypeOf(Any) {
		args := []reflect.Value{reflect.ValueOf(v2)}
		return reflect.ValueOf(v1).Call(args)[0].Bool()
	}
	return reflect.DeepEqual(v1, v2)
}

var Any Matcher = func(_ interface{}) bool {
	return true
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

func Keys(m1 interface{}) Matcher {
	t1 := reflect.TypeOf(m1)
	if t1.Kind() != reflect.Map {
		panic("Keys matcher must be called with a map")
	}
	v1 := reflect.ValueOf(m1)
	return func(m2 interface{}) bool {
		t2 := reflect.TypeOf(m2)
		if t2.Kind() != reflect.Map {
			return false
		}
		v2 := reflect.ValueOf(m2)
		for _, key := range v1.MapKeys() {
			i1 := v1.MapIndex(key).Interface()
			key2 := v2.MapIndex(key)
			if !key2.IsValid() {
				return false
			}
			i2 := key2.Interface()
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
		sv := reflect.ValueOf(s)
		for i := 0; i < sv.Len(); i++ {
			if Match(v, sv.Index(i).Interface()) {
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

func Custom(m func(interface{}) bool) Matcher {
	return m
}

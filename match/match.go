package match

import (
	"reflect"
	"regexp"
)

type Matcher func(interface{}) bool

var Any Matcher = func(_ interface{}) bool {
	return true
}

func Match(v1, v2 interface{}) bool {
	if reflect.TypeOf(v1) == reflect.TypeOf(Any) {
		args := []reflect.Value{reflect.ValueOf(v2)}
		return reflect.ValueOf(v1).Call(args)[0].Bool()
	}
	return reflect.DeepEqual(v1, v2)
}

func Regexp(exp interface{}) Matcher {
	re, ok := exp.(*regexp.Regexp)
	if !ok {
		re = regexp.MustCompile(exp.(string))
	}
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

func Key(k, v interface{}) Matcher {
	return func(m interface{}) bool {
		t := reflect.TypeOf(m)
		if t.Kind() != reflect.Map {
			return false
		}
		mv := reflect.ValueOf(m)
		mkv := mv.MapIndex(reflect.ValueOf(k))
		if !mkv.IsValid() {
			return false
		}
		return Match(v, mkv.Interface())
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

func Field(n string, v interface{}) Matcher {
	return func(s interface{}) bool {
		t := reflect.TypeOf(s)
		if t.Kind() != reflect.Struct {
			return false
		}
		sv := reflect.ValueOf(s)
		_, hasName := t.FieldByName(n)
		if !hasName {
			return false
		}
		return Match(v, sv.FieldByName(n).Interface())
	}
}

func Custom(m func(interface{}) bool) Matcher {
	return m
}

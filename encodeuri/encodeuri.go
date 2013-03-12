package encodeuri

import (
	"reflect"
	"errors"
	"strings"
	"strconv"
	"encoding/base64"
)

func Marshal(i interface{}) (uri string, err error) {
	v := reflect.ValueOf(i)
	t := v.Type()
	if v.Kind() != reflect.Struct {
		err = errors.New("Only Support for Struct")
		return
	}
	
	for i:=0; i<v.NumField(); i++ {
		if isEmptyValue(v.Field(i)) {
			continue
		}
		
		f := t.Field(i)
		tag := f.Tag.Get("uri")
		name, opts := parseTags(tag)
		if name == "-" {
			continue
		}
		if name != "" {
			uri += "/" + name
		}
		value := UriParseValue(v.Field(i))
		if strings.Index(opts, ",encoded") != -1 {
			value = base64.URLEncoding.EncodeToString([]byte(value))
		}
		uri += "/" + value
	}
	return
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func UriParseValue(v reflect.Value) string {
	switch v.Kind() {
		case reflect.String:
			return v.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return string(strconv.AppendInt([]byte{}, v.Int(), 10))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return string(strconv.AppendUint([]byte{}, v.Uint(), 10))
	}
	return ""
}

func parseTags(tag string) (name, opt string) {
	idx := strings.Index(tag, ",")
	if idx == -1 {
		return tag, ""
	}
	return tag[:idx], tag[idx:]
}

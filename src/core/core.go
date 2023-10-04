package core

import "reflect"

func GetRealStruct(p any) *reflect.Value {
	v := reflect.ValueOf(p)

	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return &v
}

func init() {
	Log.SetFormatter()
	Log.SetLevel("debug")
}

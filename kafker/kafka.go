package main

import (
	"reflect"
	"runtime"
)

func main() {

	var mS runtime.MemStats
	runtime.ReadMemStats(&mS)

	v := reflect.ValueOf(mS)
	t := v.Type()

	mappa := map[string]float64{}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name
		switch value := field.Interface().(type) {
		case uint64:
			mappa[fieldName] = float64(value)
		case uint32:
			mappa[fieldName] = float64(value)
		case float64:
			mappa[fieldName] = float64(value)
		case bool:
			if value {
				mappa[fieldName] = 1
			} else {
				mappa[fieldName] = 0
			}
		}
	}

}

package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// ToMap converts a struct to a map using the struct's tags.
//
// ToMap uses tags on struct fields to decide which fields to add to the
// returned map.
// from: https://stackoverflow.com/questions/23589564/function-for-converting-a-struct-to-map-in-golang
func ToMap(in interface{}, tag string) (prometheus.Labels, error) {
	out := make(prometheus.Labels)

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(tag); tagv != "" {

			var value string
			fieldValue := v.Field(i).Interface()
			switch v := fieldValue.(type) { //varType := v.Field(i).Kind(); varType {
			case float64:
				value = fmt.Sprintf("%2.f", v)
			case int64:
				value = fmt.Sprintf("%d", v)
			case time.Time:
				value = fmt.Sprintf("%v", v)
			default:
				value = fmt.Sprintf("%s", v)
			}

			out[tagv] = value
		}
	}
	return out, nil
}

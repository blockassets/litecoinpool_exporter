package main

import (
	"fmt"
	"reflect"
	"strings"
)

func reflectJsonTags(str interface{}) []string {
	val := reflect.ValueOf(str)
	numFields := val.NumField()

	result := make([]string, numFields)

	for i := 0; i < numFields; i++ {
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		result[i] = strings.Split(tag.Get("json"), ",")[0]
	}

	return result
}

func reflectStructName(s interface{}) string {
	return reflect.ValueOf(s).Type().Name()
}

/** Cache struct names */
type StructNameCache map[interface{}]string

var structNameCache = make(StructNameCache)

func lookupStructName(str interface{}) string {
	val, ok := structNameCache[str]
	if !ok {
		val = reflectStructName(str)
		structNameCache[str] = val
	}
	return val
}

/** Cache tag names in structs */
type StructJsonTagCache map[interface{}][]string

var structJsonTagCache = make(StructJsonTagCache)

func lookupStructTags(str interface{}) []string {
	val, ok := structJsonTagCache[str]
	if !ok {
		val = reflectJsonTags(str)
		structJsonTagCache[str] = val
	}
	return val
}

var floatType = reflect.TypeOf(float64(0))

func ConvertToFloat(unk interface{}) (float64, error) {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)

	// Why not supported by the language?
	if v.Type().Name() == "bool" {
		if v.Bool() {
			return 1, nil
		} else {
			return 0, nil
		}
	}

	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(floatType)
	return fv.Float(), nil
}

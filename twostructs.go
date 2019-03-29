package twostructs

import (
	"errors"
	"reflect"
)

type MappingFunction struct {
	In  reflect.Type
	Out reflect.Type
	Fn  reflect.Value
}

type am struct {
	mapFns []MappingFunction
}

func New() am {
	return am{}
}

func (am *am) RegisterMappingFunction(mapFn interface{}) {
	mapFnType := reflect.TypeOf(mapFn)
	mapFnValue := reflect.ValueOf(mapFn)
	if mapFnType.Kind() != reflect.Func ||
		mapFnType.NumIn() != 1 ||
		mapFnType.NumOut() != 1 {
		panic("mapFn must be: function(IN)OUT")
	}

	am.mapFns = append(am.mapFns, MappingFunction{
		In:  mapFnType.In(0),
		Out: mapFnType.Out(0),
		Fn:  mapFnValue,
	})
}

// Struct convert a source (struct) to a destination (struct reference)
func (am *am) Struct(src, dst interface{}) error {
	srcType := reflect.TypeOf(src)
	srcValue := reflect.ValueOf(src)
	dstType := reflect.TypeOf(dst)
	dstValue := reflect.ValueOf(dst)

	if srcType.Kind() != reflect.Struct ||
		dstType.Kind() != reflect.Ptr ||
		dstType.Elem().Kind() != reflect.Struct {
		return errors.New("FromTo only accepts a source struct and destination struct pointer")
	}
	if srcType.NumField() != dstType.Elem().NumField() {
		return errors.New("structs have different lengths")
	}
	for fieldNum := 0; fieldNum < srcType.NumField(); fieldNum++ {
		srcFieldType := srcType.Field(fieldNum).Type
		srcFieldValue := srcValue.Field(fieldNum)
		dstFieldType := dstType.Elem().Field(fieldNum).Type
		dstFieldValue := dstValue.Elem().Field(fieldNum)

		// Possibly dereference a pointer on the source
		if srcFieldType.Kind() == reflect.Ptr && !srcFieldValue.IsNil() {
			srcFieldType = srcFieldType.Elem()
			srcFieldValue = srcFieldValue.Elem()
		}

		if srcFieldType.Kind() == dstFieldType.Kind() {
			dstFieldValue.Set(srcFieldValue)
		} else {
			found := false
			for _, mapFn := range am.mapFns {
				if srcFieldType == mapFn.In &&
					dstFieldType == mapFn.Out {
					mappedValues := mapFn.Fn.Call([]reflect.Value{srcFieldValue})
					dstFieldValue.Set(mappedValues[0])
					found = true
					break
				}
			}
			if !found {
				//				panic(fmt.Sprintf("call RegisterMappingFunction() with func(%s.%s)%s.%s", srcFieldType.PkgPath(), srcFieldType.Name(), dstFieldType.PkgPath(), dstFieldType.Name()))
			}
		}
	}
	return nil
}

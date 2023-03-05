package struct_utils

import "reflect"

// CopyIntersectionStruct assign values between two struct assignments(A and B) that have intersection.
// traverse all elements in A and assign to according part of B if B has that field
func CopyIntersectionStruct(src, dst interface{}) {
	sElement := reflect.ValueOf(src).Elem()
	dElement := reflect.ValueOf(dst).Elem()
	for i := 0; i < dElement.NumField(); i++ {
		dField := dElement.Type().Field(i)
		sValue := sElement.FieldByName(dField.Name)
		if !sValue.IsValid() {
			continue
		}
		value := dElement.Field(i)
		value.Set(sValue)
	}
}

type valueIsEmptyFunc func(v reflect.Value) bool

var valueIsEmptyFuncMap = map[reflect.Kind]valueIsEmptyFunc{
	reflect.Int:        intValueIsEmpty,
	reflect.Int8:       intValueIsEmpty,
	reflect.Int16:      intValueIsEmpty,
	reflect.Int32:      intValueIsEmpty,
	reflect.Int64:      intValueIsEmpty,
	reflect.Uint:       uintValueIsEmpty,
	reflect.Uint8:      uintValueIsEmpty,
	reflect.Uint16:     uintValueIsEmpty,
	reflect.Uint32:     uintValueIsEmpty,
	reflect.Uint64:     uintValueIsEmpty,
	reflect.Complex64:  completeValueIsEmpty,
	reflect.Complex128: completeValueIsEmpty,
	reflect.Array:      sliceValueIsEmpty,
	reflect.Slice:      sliceValueIsEmpty,
	reflect.String:     stringValueIsEmpty,
}

func intValueIsEmpty(v reflect.Value) bool {
	return v.Int() == 0
}

func uintValueIsEmpty(v reflect.Value) bool {
	return v.Uint() == 0
}

func completeValueIsEmpty(v reflect.Value) bool {
	return v.Complex() == 0
}

func stringValueIsEmpty(v reflect.Value) bool {
	return len(v.String()) == 0
}

func sliceValueIsEmpty(v reflect.Value) bool {
	return v.IsNil() || v.Len() == 0
}

// IsStructEmpty
// @Description: reflect.DeepEqual() 的方式判空的话如果slice是 初始化了长度为0 则无法判断，通过该方法slice为0仍然会认定是空
func IsStructEmpty(v interface{}) bool {
	vType := reflect.TypeOf(v)
	value := reflect.ValueOf(v)
	for i := 0; i < value.NumField(); i++ {
		field := vType.Field(i)
		// 结构体变量递归判断
		if field.Type.Kind() == reflect.Struct && field.Type.NumField() > 0 {
			if isEmpty := IsStructEmpty(value.FieldByName(field.Name).Interface()); !isEmpty {
				return false
			}
		}
		if isEmptyFunc, ok := valueIsEmptyFuncMap[field.Type.Kind()]; ok {
			if isEmpty := isEmptyFunc(value.Field(i)); !isEmpty {
				return false
			}
		}
	}
	return true
}

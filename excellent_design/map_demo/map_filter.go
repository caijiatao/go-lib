package map_demo

import "reflect"

func FilterUpdateMapNil(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		vType := reflect.TypeOf(v)
		if vType == nil {
			delete(m, k)
			continue
		}
		if vType.Kind() != reflect.Ptr {
			continue
		}
		if reflect.ValueOf(v).IsNil() {
			delete(m, k)
			continue
		}
	}
	return m
}

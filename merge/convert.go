package merge

import (
	"fmt"
)

// Convert converts map[interface{}]interface{} to map[string]interface{} which
// is mergable by merge.Maps
func Convert(m map[interface{}]interface{}) map[string]interface{} {
	return convert(m)
}

func convert(m map[interface{}]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range m {
		var value interface{}
		switch v2 := v.(type) {
		case map[interface{}]interface{}:
			value = convert(v2)
		case []map[interface{}]interface{}:
			arr := make([]map[string]interface{}, len(v2))
			for i, m := range v2 {
				arr[i] = convert(m)
			}
			value = arr
		default:
			value = v
		}
		key := "null"
		if k != nil {
			key = fmt.Sprint(k)
		}
		res[key] = value
	}
	return res
}

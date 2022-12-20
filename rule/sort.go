// Copyright 2020 Jebbs. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package rule

import "sort"

// sortByFields sort slice elements by specified fields
func sortByFields(slice []interface{}, fields []Field) {
	if len(fields) == 0 {
		return
	}
	sort.Slice(
		slice,
		func(i, j int) bool {
			return getPriority(slice[i], fields) < getPriority(slice[j], fields)
		},
	)
}

func getPriority(v interface{}, fields []Field) float64 {
	value := getField(v, fields)
	if value == nil {
		return 0
	}
	switch num := value.(type) {
	case float64:
		return num
	case float32:
		return float64(num)
	}
	return 0
}

func getField(v interface{}, fields []Field) interface{} {
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	for _, field := range fields {
		if p, ok := m[field.Key]; ok {
			return p
		}
	}
	return nil
}

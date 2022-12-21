// Copyright 2020 Jebbs. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package rule

import "sort"

type meta struct {
	index    int
	priority float64
	value    interface{}
}

// sortByFields sort slice elements by specified fields
func sortByFields(slice []interface{}, fields []Field) {
	if len(slice) == 0 || len(fields) == 0 {
		return
	}
	metas := make([]meta, len(slice))
	for i, v := range slice {
		metas[i] = meta{
			index:    i,
			priority: getPriority(v, fields),
			value:    v,
		}
	}
	sort.Slice(
		metas,
		func(i, j int) bool {
			if metas[i].priority != metas[j].priority {
				return metas[i].priority < metas[j].priority
			}
			return metas[i].index < metas[j].index
		},
	)
	for i, m := range metas {
		slice[i] = m.value
	}
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

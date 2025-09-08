// Copyright 2020 Jebbs. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package jsons

// apply applies rule according to m
func (r *Options) apply(m map[string]interface{}) error {
	if r == nil || (len(r.MergeBy) == 0 && len(r.OrderBy) == 0) {
		return nil
	}
	err := r.sortMergeSlices(m)
	if err != nil {
		return err
	}
	r.removeHelperFields(m)
	return nil
}

// sortMergeSlices enumerates all slices in a map, to sort by order and merge by tag
func (r *Options) sortMergeSlices(target map[string]interface{}) error {
	for key, value := range target {
		if slice, ok := value.([]interface{}); ok {
			sortByFields(slice, r.OrderBy)
			s, err := mergeByFields(r.TypeOverride, slice, r.MergeBy)
			if err != nil {
				return err
			}
			target[key] = s
			for _, item := range s {
				if m, ok := item.(map[string]interface{}); ok {
					r.sortMergeSlices(m)
				}
			}
		} else if field, ok := value.(map[string]interface{}); ok {
			r.sortMergeSlices(field)
		}
	}
	return nil
}

func (r *Options) removeHelperFields(target map[string]interface{}) {
	for key, value := range target {
		if r.shouldDelete(key) {
			delete(target, key)
		} else if slice, ok := value.([]interface{}); ok {
			for _, e := range slice {
				if el, ok := e.(map[string]interface{}); ok {
					r.removeHelperFields(el)
				}
			}
		} else if field, ok := value.(map[string]interface{}); ok {
			r.removeHelperFields(field)
		}
	}
}

// shouldDelete tells if the field should be deleted according to the rules
func (r *Options) shouldDelete(key string) bool {
	for _, field := range r.MergeBy {
		if key != field.Name {
			continue
		}
		return field.Remove
	}
	for _, field := range r.OrderBy {
		if key != field.Name {
			continue
		}
		return field.Remove
	}
	return false
}

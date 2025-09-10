// Copyright 2020 Jebbs. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package jsons

import "github.com/qjebbs/go-jsons/internal/ordered"

// apply applies rule according to m
func (r *options) apply(m *ordered.Map) error {
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
func (r *options) sortMergeSlices(target *ordered.Map) error {
	for key, value := range target.Values {
		if slice, ok := value.([]interface{}); ok {
			sortByFields(slice, r.OrderBy)
			s, err := mergeByFields(r.TypeOverride, slice, r.MergeBy)
			if err != nil {
				return err
			}
			target.Set(key, s)
			for _, item := range s {
				if m, ok := item.(*ordered.Map); ok {
					r.sortMergeSlices(m)
				}
			}
		} else if field, ok := value.(*ordered.Map); ok {
			r.sortMergeSlices(field)
		}
	}
	return nil
}

func (r *options) removeHelperFields(target *ordered.Map) {
	for key, value := range target.Values {
		if r.shouldDelete(key) {
			target.Remove(key)
		} else if slice, ok := value.([]interface{}); ok {
			for _, e := range slice {
				if el, ok := e.(*ordered.Map); ok {
					r.removeHelperFields(el)
				}
			}
		} else if field, ok := value.(*ordered.Map); ok {
			r.removeHelperFields(field)
		}
	}
}

// shouldDelete tells if the field should be deleted according to the rules
func (r *options) shouldDelete(key string) bool {
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

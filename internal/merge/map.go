// Copyright 2020 Jebbs. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package merge

import (
	"fmt"
	"reflect"
)

// Maps merges source maps into target
func Maps(typeOverride bool, target map[string]interface{}, sources ...map[string]interface{}) (err error) {
	for _, source := range sources {
		err = mergeMap(typeOverride, target, source)
		if err != nil {
			return err
		}
	}
	return nil
}

// mergeMap merges source map into target
// it supports only map[string]interface{} type for any children of the map tree
func mergeMap(typeOverride bool, target map[string]interface{}, source map[string]interface{}) (err error) {
	for key, value := range source {
		target[key], err = mergeField(typeOverride, target[key], value)
		if err != nil {
			return fmt.Errorf("field '%s': %s", key, err)
		}
	}
	return nil
}

func mergeField(typeOverride bool, target interface{}, source interface{}) (interface{}, error) {
	if source == nil {
		return target, nil
	}
	if target == nil {
		return source, nil
	}
	if reflect.TypeOf(source) != reflect.TypeOf(target) {
		if !typeOverride {
			return nil, fmt.Errorf("type mismatch, expect %T, incoming %T", target, source)
		}
		return source, nil
	}
	if slice, ok := source.([]interface{}); ok {
		tslice, _ := target.([]interface{})
		tslice = append(tslice, slice...)
		return tslice, nil
	}
	if smap, ok := source.(map[string]interface{}); ok {
		tmap, _ := target.(map[string]interface{})
		err := mergeMap(typeOverride, tmap, smap)
		return tmap, err
	}
	return source, nil
}

package merge

import (
	"fmt"
	"reflect"

	"github.com/qjebbs/go-jsons/internal/ordered"
)

// OrderedMaps merges source ordered maps into target
func OrderedMaps(typeOverride bool, target *ordered.Map, sources ...*ordered.Map) (err error) {
	for _, source := range sources {
		err = mergeOrderedMap(typeOverride, target, source)
		if err != nil {
			return err
		}
	}
	return nil
}

func mergeOrderedMap(typeOverride bool, target *ordered.Map, source *ordered.Map) (err error) {
	for _, sk := range source.Keys {
		if _, exists := target.Values[sk]; !exists {
			target.Keys = append(target.Keys, sk)
		}
	}
	for key, value := range source.Values {
		target.Values[key], err = mergeOrderedField(typeOverride, target.Values[key], value)
		if err != nil {
			return fmt.Errorf("field '%s': %s", key, err)
		}
	}
	return nil
}

func mergeOrderedField(typeOverride bool, target interface{}, source interface{}) (interface{}, error) {
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
	if smap, ok := source.(*ordered.Map); ok {
		tmap, _ := target.(*ordered.Map)
		err := mergeOrderedMap(typeOverride, tmap, smap)
		return tmap, err
	}
	return source, nil
}

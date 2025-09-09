package ordered

import (
	"encoding/json"
	"sort"
)

var _ json.Marshaler = &Map{}
var _ json.Unmarshaler = &Map{}

// Map represents a JSON object that maintains the order of its keys.
type Map struct {
	Values map[string]interface{}
	Keys   []string
}

// New creates a new empty Ordered object.
func New() *Map {
	return &Map{
		Values: make(map[string]interface{}),
		Keys:   []string{},
	}
}

// FromMap creates an Ordered object from a standard map.
func FromMap(m map[string]interface{}) *Map {
	o := New()
	for k, v := range m {
		o.Keys = append(o.Keys, k)
		var value any
		switch v := v.(type) {
		case []interface{}:
			slice := make([]interface{}, len(v))
			for i, e := range v {
				if em, ok := e.(map[string]interface{}); ok {
					slice[i] = FromMap(em)
				} else {
					slice[i] = e
				}
			}
			value = slice
		case map[string]interface{}:
			value = FromMap(v)
		default:
			value = v
		}
		o.Values[k] = value
	}
	// sort.Slice(o.Keys, func(i, j int) bool {
	// 	return o.Keys[i] < o.Keys[j]
	// })
	return o
}

// Remove removes a key from the Ordered object.
func (o *Map) Remove(key string) {
	if _, exists := o.Values[key]; exists {
		delete(o.Values, key)
		for i, k := range o.Keys {
			if k == key {
				o.Keys = append(o.Keys[:i], o.Keys[i+1:]...)
				break
			}
		}
	}
}

// Set sets a key-value pair in the Ordered object.
func (o *Map) Set(key string, value interface{}) {
	if _, exists := o.Values[key]; !exists {
		o.Keys = append(o.Keys, key)
	}
	o.Values[key] = value
}

// Sort sorts the keys of the Ordered object in ascending order.
func (o *Map) Sort(less ...func(a, b string) bool) *Map {
	var fn func(a, b string) bool
	if len(less) > 0 {
		fn = less[0]
	} else {
		fn = func(a, b string) bool {
			return a < b
		}
	}
	sort.Slice(o.Keys, func(i, j int) bool {
		return fn(o.Keys[i], o.Keys[j])
	})
	for _, v := range o.Values {
		if child, ok := v.(*Map); ok {
			child.Sort(fn)
		} else if slice, ok := v.([]interface{}); ok {
			for _, item := range slice {
				if child, ok := item.(*Map); ok {
					child.Sort(fn)
				}
			}
		}
	}
	return o
}

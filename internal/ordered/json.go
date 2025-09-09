package ordered

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// MarshalJSON implements the json.Marshaler interface.
func (o Map) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	for i, k := range o.Keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		// add key
		if err := encoder.Encode(k); err != nil {
			return nil, err
		}
		buf.WriteByte(':')
		// add value
		if err := encoder.Encode(o.Values[k]); err != nil {
			return nil, err
		}
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (o *Map) UnmarshalJSON(data []byte) error {
	if o.Values == nil {
		o.Values = map[string]interface{}{}
	}

	dec := json.NewDecoder(bytes.NewReader(data))

	token, err := dec.Token()
	if err != nil {
		return err
	}
	// json: cannot unmarshal array into Go value of type map[string]interface {}
	if delim, ok := token.(json.Delim); ok {
		if delim == '[' {
			return fmt.Errorf("json: cannot unmarshal array into Go value of type %T", o)
		}
	} else if num, ok := token.(float64); ok {
		return fmt.Errorf("json: cannot unmarshal number %v into Go value of type %T", num, o)
	} else if str, ok := token.(string); ok {
		return fmt.Errorf("json: cannot unmarshal string %q into Go value of type %T", str, o)
	} else if b, ok := token.(bool); ok {
		return fmt.Errorf("json: cannot unmarshal bool %v into Go value of type %T", b, o)
	} else if token == nil {
		return fmt.Errorf("json: cannot unmarshal null into Go value of type %T", o)
	}

	o.Keys = make([]string, 0)
	hasKey := make(map[string]bool)

	for dec.More() {
		// read key
		token, err := dec.Token()
		if err != nil {
			return err
		}
		key := token.(string)

		// duplicate key, remove previous occurrence
		if hasKey[key] {
			for j, k := range o.Keys {
				if k == key {
					copy(o.Keys[j:], o.Keys[j+1:])
					o.Keys = o.Keys[:len(o.Keys)-1]
					break
				}
			}
		} else {
			hasKey[key] = true
		}
		o.Keys = append(o.Keys, key)

		// read value
		var value shallowParser
		if err := dec.Decode(&value); err != nil {
			return err
		}
		v, err := o.processNestedValue(value)
		if err != nil {
			return err
		}
		o.Values[key] = v
	}

	// expect '}'
	token, err = dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := token.(json.Delim); !ok || delim != '}' {
		return fmt.Errorf("expected end of object, got %v", token)
	}

	return nil
}

func (o *Map) processNestedValue(value shallowParser) (interface{}, error) {
	if value.Val != nil {
		var r interface{}
		if err := json.Unmarshal(value.Val, &r); err != nil {
			return nil, err
		}
		return r, nil
	}
	if value.Obj != nil {
		nested := &Map{}
		if err := nested.UnmarshalJSON(value.Obj); err != nil {
			return nil, err
		}
		return nested, nil
	}
	r := make([]interface{}, len(value.Arr))
	for i, item := range value.Arr {
		var sv shallowParser
		if err := json.Unmarshal(item, &sv); err != nil {
			return nil, err
		}
		v, err := o.processNestedValue(sv)
		if err != nil {
			return nil, err
		}
		r[i] = v
	}
	return r, nil
}

var _ json.Unmarshaler = &Map{}

type shallowParser struct {
	Obj json.RawMessage
	Arr []json.RawMessage
	Val json.RawMessage
}

func (s *shallowParser) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	token, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := token.(json.Delim); ok {
		switch delim {
		case '{':
			s.Obj = data
		case '[':
			for dec.More() {
				var raw json.RawMessage
				if err := dec.Decode(&raw); err != nil {
					return err
				}
				s.Arr = append(s.Arr, raw)
			}
			// expect ']'
			_, err = dec.Token()
			return err
		}
		return nil
	}
	s.Val = data
	return nil
}

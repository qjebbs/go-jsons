package ordered_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/qjebbs/go-jsons/internal/ordered"
)

func TestOrderedJSON(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i := 97; i <= 108; i++ {
		buf.WriteString(fmt.Sprintf(`"%s":0,`, string(byte(i))))
	}
	buf.WriteString(`"m":[{`)
	for i := 97; i <= 105; i++ {
		buf.WriteString(fmt.Sprintf(`"%s":0`, string(byte(i))))
		if i < 105 {
			buf.WriteByte(',')
		}
	}
	buf.WriteString(`}],`)
	for i := 110; i <= 122; i++ {
		buf.WriteString(fmt.Sprintf(`"%s":0`, string(byte(i))))
		if i < 122 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte('}')
	want := buf.Bytes()
	// t.Logf("raw: %s", string(want))
	o := &ordered.Map{}
	err := json.Unmarshal(want, o)
	if err != nil {
		t.Errorf("UnmarshalJSON failed: %v", err)
	}
	got, err := json.Marshal(o)
	if err != nil {
		t.Errorf("MarshalJSON failed: %v", err)
	}
	if string(got) != string(want) {
		t.Errorf("MarshalJSON result mismatch, want: %s, got: %s", want, string(got))
	}
}

func TestOrderedFromMap(t *testing.T) {
	m := map[string]interface{}{
		"a": 0,
		"b": 0,
		"c": map[string]interface{}{
			"d": 0,
			"e": 0,
		},
		"f": []interface{}{
			map[string]interface{}{
				"g": 0,
				"h": 0,
			},
			123,
			"string",
		},
	}
	got := ordered.FromMap(m).Sort()
	want := &ordered.Map{
		Keys: []string{"a", "b", "c", "f"},
		Values: map[string]interface{}{
			"a": 0,
			"b": 0,
			"c": &ordered.Map{
				Keys: []string{"d", "e"},
				Values: map[string]interface{}{
					"d": 0,
					"e": 0,
				},
			},
			"f": []interface{}{
				&ordered.Map{
					Keys: []string{"g", "h"},
					Values: map[string]interface{}{
						"g": 0,
						"h": 0,
					},
				},
				123,
				"string",
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FromMap result mismatch, want: %+v, got: %+v", want, got)
	}
}

func TestOrderedSetRemove(t *testing.T) {
	o := ordered.New()
	o.Set("a", 1)
	o.Set("b", 2)
	o.Set("c", 3)
	o.Set("d", 4)
	o.Set("e", 5)
	o.Remove("c")
	want := &ordered.Map{
		Keys:   []string{"a", "b", "d", "e"},
		Values: map[string]interface{}{"a": 1, "b": 2, "d": 4, "e": 5},
	}
	if !reflect.DeepEqual(o, want) {
		t.Errorf("Set or Remove result mismatch, want: %+v, got: %+v", want, o)
	}
}

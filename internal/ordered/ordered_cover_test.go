package ordered_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/qjebbs/go-jsons/internal/ordered"
)

func TestOrderedUnmarshalDupKey(t *testing.T) {
	raw := []byte(`{"a":1,"b":2,"a":3}`)
	o := ordered.New()
	err := json.Unmarshal(raw, o)
	if err != nil {
		t.Errorf("UnmarshalJSON failed: %v", err)
	}
	want := &ordered.Map{
		Keys: []string{"b", "a"},
		Values: map[string]interface{}{
			"a": 3,
			"b": 2,
		},
	}
	if reflect.DeepEqual(o, want) {
		t.Errorf("UnmarshalJSON result mismatch, want: %+v, got: %+v", want, o)
	}
}

func TestOrderedUnmarshalErrors(t *testing.T) {
	o := ordered.New()
	expectError(t, json.Unmarshal([]byte(`[1,2,3]`), o))
	expectError(t, json.Unmarshal([]byte(`1`), o))
	expectError(t, json.Unmarshal([]byte(`true`), o))
	expectError(t, json.Unmarshal([]byte(`"string"`), o))
	expectError(t, json.Unmarshal([]byte(`null`), o))
}

func TestOrderedMarshalErrors(t *testing.T) {
	expectError(t, func() error {
		_, err := json.Marshal(ordered.Map{
			Keys: []string{"a"},
			Values: map[string]interface{}{
				"a": func() {}, // functions cannot be marshaled to JSON
			},
		})
		return err
	}())
}

func expectError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("Expected error, got nil")
		return
	}
}

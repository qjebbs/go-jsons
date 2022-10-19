package merge_test

import (
	"reflect"
	"testing"

	"github.com/qjebbs/go-jsons/merge"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		value map[interface{}]interface{}
		want  map[string]interface{}
	}{
		{
			value: map[interface{}]interface{}{"a": 1},
			want:  map[string]interface{}{"a": 1},
		},
		{
			value: map[interface{}]interface{}{nil: 1},
			want:  map[string]interface{}{"null": 1},
		},
		{
			value: map[interface{}]interface{}{"a": map[interface{}]interface{}{"b": 1}},
			want:  map[string]interface{}{"a": map[string]interface{}{"b": 1}},
		},
		{
			value: map[interface{}]interface{}{"a": []map[interface{}]interface{}{
				{"b": 1},
				{"c": 1},
			}},
			want: map[string]interface{}{"a": []map[string]interface{}{
				{"b": 1},
				{"c": 1},
			}},
		},
	}
	for _, tt := range tests {
		got := merge.Convert(tt.value)
		if !reflect.DeepEqual(tt.want, got) {
			t.Errorf("want:\n%v\n\ngot:\n%v", tt.want, got)
		}
	}
}

package helper_test

import (
	"reflect"
	"testing"

	"github.com/qjebbs/go-jsons/helper"
	"gopkg.in/yaml.v3"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		value string
		want  map[string]interface{}
	}{
		{
			value: `a: 1`,
			want:  map[string]interface{}{"a": 1},
		},
		{
			value: `nil: 1`,
			want:  map[string]interface{}{"nil": 1},
		},
		{
			value: "a:\n  b: \n    c: 1",
			want: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": 1,
					},
				},
			},
		},
		{
			value: "a:\n- b:\n    c: 1\n- d: 1",
			want: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{
						"b": map[string]interface{}{
							"c": 1,
						},
					},
					map[string]interface{}{"d": 1},
				},
			},
		},
	}
	for _, tt := range tests {
		m := make(map[interface{}]interface{})
		yaml.Unmarshal([]byte(tt.value), m)
		got := helper.ConvertYAMLMap(m)
		if !reflect.DeepEqual(tt.want, got) {
			t.Errorf("want:\n%#v\n\ngot:\n%#v", tt.want, got)
		}
	}
}

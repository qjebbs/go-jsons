package jsons

import (
	"reflect"
	"testing"
)

func TestMergeMaps(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name         string
		typeOverride bool
		values       []map[string]interface{}
		want         map[string]interface{}
		wantErr      bool
	}{
		{
			name: "fields_merge",
			values: []map[string]interface{}{
				{"a": 1},
				{"b": 2},
			},
			want: map[string]interface{}{
				"a": 1,
				"b": 2,
			},
		},
		{
			name: "values_merge",
			values: []map[string]interface{}{
				{"a": []interface{}{1, 2}},
				{"a": []interface{}{3, 4}},
			},
			want: map[string]interface{}{
				"a": []interface{}{1, 2, 3, 4},
			},
		},
		{
			name: "not_overwrite_nil",
			values: []map[string]interface{}{
				{"a": 1},
				{"a": nil},
			},
			want: map[string]interface{}{
				"a": 1,
			},
		},
		{
			name: "merge_deep",
			values: []map[string]interface{}{
				{"a": map[string]interface{}{"b": 1}},
				{"a": map[string]interface{}{"b": 2}},
			},
			want: map[string]interface{}{
				"a": map[string]interface{}{"b": 2},
			},
		},
		{
			name:         "type_override",
			typeOverride: true,
			values: []map[string]interface{}{
				{"a": false},
				{"a": 3},
			},
			want: map[string]interface{}{
				"a": 3,
			},
		},
		{
			name: "mismatch_type",
			values: []map[string]interface{}{
				{"a": false},
				{"a": 3},
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := make(map[string]interface{})
			err := mergeMaps(tc.typeOverride, got, tc.values...)
			switch tc.wantErr {
			case true:
				if err == nil {
					t.Fatalf("want err got nil")
				}
			case false:
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(tc.want, got) {
					t.Errorf("want:\n%v\n\ngot:\n%v", tc.want, got)
				}
			}
		})
	}
}

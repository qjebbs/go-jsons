package merge_test

import (
	"reflect"
	"testing"

	"github.com/qjebbs/go-jsons/merge"
)

func TestMap(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		values  []map[string]interface{}
		want    map[string]interface{}
		wantErr bool
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
			err := merge.Maps(got, tc.values...)
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

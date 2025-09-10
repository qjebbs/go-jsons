package merge_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/qjebbs/go-jsons/internal/merge"
)

var testCases = []struct {
	name         string
	typeOverride bool
	values       []string
	want         string
	wantErr      bool
}{
	{
		name: "fields_merge",
		values: []string{
			`{"a": 1}`,
			`{"a": 2}`,
			`{"b": 2}`,
		},
		want: `{"a": 2, "b": 2}`,
	},
	{
		name: "values_merge",
		values: []string{
			`{"a": [1, 2]}`,
			`{"a": [3, 4]}`,
		},
		want: `{"a": [1, 2, 3, 4]}`,
	},
	{
		name: "not_overwrite_nil",
		values: []string{
			`{"a": 1}`,
			`{"a": null}`,
		},
		want: `{"a": 1}`,
	},
	{
		name: "merge_deep",
		values: []string{
			`{"a": {"b": [1, 2]}}`,
			`{"a": {"b": [3 ,4]}}`,
		},
		want: `{"a": {"b": [1, 2, 3, 4]}}`,
	},
	{
		name:         "type_override",
		typeOverride: true,
		values: []string{
			`{"a": false}`,
			`{"a": 3}`,
		},
		want: `{"a": 3}`,
	},
	{
		name: "mismatch_type",
		values: []string{
			`{"a": false}`,
			`{"a": 3}`,
		},
		wantErr: true,
	},
}

func TestMergeMaps(t *testing.T) {
	t.Parallel()
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var want map[string]interface{}
			if tc.want != "" {
				want = convertToMap(t, tc.want)
			}
			items := convertToMaps(t, tc.values)
			got := make(map[string]interface{})
			err := merge.Maps(tc.typeOverride, got, items...)
			switch tc.wantErr {
			case true:
				if err == nil {
					t.Fatalf("want err got nil")
				}
			case false:
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(want, got) {
					t.Errorf("want:\n%v\n\ngot:\n%v", want, got)
				}
			}
		})
	}
}

func convertToMaps(t *testing.T, values []string) []map[string]interface{} {
	t.Helper()
	var maps []map[string]interface{}
	for _, v := range values {
		maps = append(maps, convertToMap(t, v))
	}
	return maps
}

func convertToMap(t *testing.T, value string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(value), &m); err != nil {
		t.Fatalf("convert %q: %s", value, err)
	}
	return m
}

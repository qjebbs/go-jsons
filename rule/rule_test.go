package rule_test

import (
	"reflect"
	"testing"

	"github.com/qjebbs/go-jsons/rule"
)

var testRule = rule.NewRule(
	rule.OrderBy("priority"),
	rule.MergeBy("tag"),
	rule.OrderByAndRemove("_priority"),
	rule.MergeByAndRemove("_tag"),
)

func TestRules(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		value   map[string]interface{}
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "merge_and_sort",
			value: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"_tag": "a", "value": 1},
					map[string]interface{}{"_tag": "b", "_priority": -100, "value": 2},
					map[string]interface{}{"_tag": "a", "value": 0},
				},
			},
			want: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"value": 0},
					map[string]interface{}{"value": 2},
				},
			},
		},
		{
			name: "as_is",
			value: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"value": 0},
					map[string]interface{}{"value": 1},
				},
			},
			want: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"value": 0},
					map[string]interface{}{"value": 1},
				},
			},
		},
		{
			name: "apply_deep",
			value: map[string]interface{}{
				"a": map[string]interface{}{
					"b": []interface{}{
						map[string]interface{}{"_tag": "a.b", "value": 0},
						map[string]interface{}{"_tag": "a.b", "value": 1},
					},
				},
			},
			want: map[string]interface{}{
				"a": map[string]interface{}{
					"b": []interface{}{
						map[string]interface{}{"value": 1},
					},
				},
			},
		},
		{
			name: "invalid_tag",
			value: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"tag": 1, "value": 0},
					map[string]interface{}{"tag": 1, "value": 1},
				},
			},
			want: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"tag": 1, "value": 0},
					map[string]interface{}{"tag": 1, "value": 1},
				},
			},
		},
		{
			name: "sort_slice",
			value: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"priority": float64(1)},
					map[string]interface{}{"priority": float32(0)},
				},
			},
			want: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"priority": float32(0)},
					map[string]interface{}{"priority": float64(1)},
				},
			},
		},
		{
			name: "empty_slice",
			value: map[string]interface{}{
				"a": []interface{}(nil),
			},
			want: map[string]interface{}{
				"a": []interface{}(nil),
			},
		},
		{
			name: "simple_value_slice",
			value: map[string]interface{}{
				"a": []interface{}{1, 2, 3},
			},
			want: map[string]interface{}{
				"a": []interface{}{1, 2, 3},
			},
		},
		{
			name: "merge_fail",
			value: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"_tag": "a", "value": 1},
					map[string]interface{}{"_tag": "a", "value": false},
				},
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := testRule.Apply(tc.value)
			switch tc.wantErr {
			case true:
				if err == nil {
					t.Fatal("want err got nil")
				}
			case false:
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(tc.want, tc.value) {
					t.Fatalf("want:\n%v\n\ngot:\n%v", tc.want, tc.value)
				}
			}
		})
	}
}

func TestNils(t *testing.T) {
	t.Parallel()
	err := (*rule.Rule)(nil).Apply(nil)
	if err != nil {
		t.Fatalf("want nil, got err: %s", err)
	}
	err = testRule.Apply(nil)
	if err != nil {
		t.Fatalf("want nil, got err: %s", err)
	}
}

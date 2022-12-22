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
			name: "sort_slice",
			value: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"priority": float64(12)},
					map[string]interface{}{"priority": float32(11)},
					map[string]interface{}{"priority": int(10)},
					map[string]interface{}{"priority": int8(9)},
					map[string]interface{}{"priority": int16(8)},
					map[string]interface{}{"priority": int32(7)},
					map[string]interface{}{"priority": int64(6)},
					map[string]interface{}{"priority": uint(5)},
					map[string]interface{}{"priority": uint8(4)},
					map[string]interface{}{"priority": uint16(3)},
					map[string]interface{}{"priority": uint32(2)},
					map[string]interface{}{"priority": uint64(1)},
					map[string]interface{}{"priority": "str"},
					map[string]interface{}{"priority": nil},
					map[string]interface{}{},
				},
			},
			want: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"priority": "str"},
					map[string]interface{}{"priority": nil},
					map[string]interface{}{},
					map[string]interface{}{"priority": uint64(1)},
					map[string]interface{}{"priority": uint32(2)},
					map[string]interface{}{"priority": uint16(3)},
					map[string]interface{}{"priority": uint8(4)},
					map[string]interface{}{"priority": uint(5)},
					map[string]interface{}{"priority": int64(6)},
					map[string]interface{}{"priority": int32(7)},
					map[string]interface{}{"priority": int16(8)},
					map[string]interface{}{"priority": int8(9)},
					map[string]interface{}{"priority": int(10)},
					map[string]interface{}{"priority": float32(11)},
					map[string]interface{}{"priority": float64(12)},
				},
			},
		},
		{
			name: "multi_tag_sort",
			value: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"value": 0},
					map[string]interface{}{"_priority": 1, "priority": -1, "value": 1},
				},
			},
			want: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"priority": -1, "value": 1},
					map[string]interface{}{"value": 0},
				},
			},
		},
		{
			name: "sort_then_merge",
			value: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"_tag": "a", "value": 1},
					map[string]interface{}{"_tag": "a", "_priority": 100, "value": 2},
					map[string]interface{}{"_tag": "a", "value": 0},
				},
			},
			want: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"value": 2},
				},
			},
		},
		{
			name: "multi_tag_merge",
			value: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"_tag": "a", "tag": "b", "value": 0},
					map[string]interface{}{"_tag": "c", "tag": "a", "value": 1},
				},
			},
			want: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"tag": "a", "value": 1},
				},
			},
		},
		{
			name: "as_is",
			value: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"tag": "a"},
					map[string]interface{}{"tag": "b"},
					1, false, "str",
				},
			},
			want: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"tag": "a"},
					map[string]interface{}{"tag": "b"},
					1, false, "str",
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
			name: "empty_slice",
			value: map[string]interface{}{
				"a": []interface{}(nil),
			},
			want: map[string]interface{}{
				"a": []interface{}(nil),
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

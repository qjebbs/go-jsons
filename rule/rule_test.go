package rule_test

import (
	"reflect"
	"testing"

	"github.com/qjebbs/go-jsons/rule"
)

var testRule = rule.NewRule(
	rule.OrderByAndRemove("_priority"),
	rule.MergeByAndRemove("_tag"),
)

func TestRules(t *testing.T) {
	tests := []struct {
		value   map[string]interface{}
		want    map[string]interface{}
		wantErr bool
	}{
		{
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
			value: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"_tag": "a", "value": 1},
					map[string]interface{}{"_tag": "a", "value": false},
				},
			},
			wantErr: true,
		},
	}
	for i, tt := range tests {
		err := testRule.Apply(tt.value)
		switch tt.wantErr {
		case true:
			if err == nil {
				t.Errorf("#%d: want err got nil", i)
				continue
			}
		case false:
			if err != nil {
				t.Errorf("#%d: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(tt.want, tt.value) {
				t.Errorf("want:\n%v\n\ngot:\n%v", tt.want, tt.value)
			}
		}
	}
}

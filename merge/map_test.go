package merge_test

import (
	"reflect"
	"testing"

	"github.com/qjebbs/go-jsons/merge"
)

func TestMap(t *testing.T) {
	tests := []struct {
		values  []map[string]interface{}
		want    map[string]interface{}
		wantErr bool
	}{
		{
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
			values: []map[string]interface{}{
				{"a": []interface{}{1, 2}},
				{"a": []interface{}{3, 4}},
			},
			want: map[string]interface{}{
				"a": []interface{}{1, 2, 3, 4},
			},
		},
		{
			values: []map[string]interface{}{
				{"a": false},
				{"a": 3},
			},
			wantErr: true,
		},
	}
	for i, tt := range tests {
		got := make(map[string]interface{})
		err := merge.Maps(got, tt.values...)
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
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want:\n%v\n\ngot:\n%v", tt.want, got)
			}
		}
	}
}

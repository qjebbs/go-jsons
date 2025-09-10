package merge_test

import (
	"reflect"
	"testing"

	"github.com/qjebbs/go-jsons/internal/merge"
	"github.com/qjebbs/go-jsons/internal/ordered"
)

func TestMergeOrdered(t *testing.T) {
	t.Parallel()
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var want *ordered.Map
			if tc.want != "" {
				want = ordered.FromMap(convertToMap(t, tc.want)).Sort()
			}
			items := convertToOrderedMaps(t, tc.values)
			got := ordered.New()
			err := merge.OrderedMaps(tc.typeOverride, got, items...)
			got.Sort()
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

func convertToOrderedMaps(t *testing.T, values []string) []*ordered.Map {
	t.Helper()
	m := convertToMaps(t, values)
	s := make([]*ordered.Map, 0, len(values))
	for _, v := range m {
		m := ordered.FromMap(v)
		s = append(s, m)
	}
	return s
}

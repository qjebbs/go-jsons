package jsons_test

import (
	"testing"

	"github.com/qjebbs/go-jsons"
	"github.com/qjebbs/go-jsons/merge"
	"gopkg.in/yaml.v2"
)

func TestMergeMixFormats(t *testing.T) {
	const FormatYAML jsons.Format = "yaml"
	m := jsons.NewMerger()
	m.RegisterDefaultLoader()
	err := m.RegisterLoader(
		FormatYAML,
		[]string{".yaml", ".yml"},
		func(b []byte) (map[string]interface{}, error) {
			m1 := make(map[interface{}]interface{})
			err := yaml.Unmarshal(b, &m1)
			if err != nil {
				return nil, err
			}
			m2 := merge.Convert(m1)
			return m2, nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	a := []byte(`{"a": 1}`)
	b := []byte(`b: 1`)
	want := []byte(`{"a":1,"b":1}`)
	got, err := m.Merge(a, b)
	if err != nil {
		t.Fatal(err)
	}
	assertJSONEqual(t, want, got)
}

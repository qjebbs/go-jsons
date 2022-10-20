package jsons_test

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/qjebbs/go-jsons"
	"gopkg.in/yaml.v3"
)

func TestMerge(t *testing.T) {
	a := []byte(`{"a":1}`)
	b := []byte(`{"b":1}`)
	c := strings.NewReader(`{"c":1}`)
	want := []byte(`{"a":1,"b":1,"c":1}`)
	got, err := jsons.Merge(a, b, c)
	if err != nil {
		t.Error(err)
	}
	assertJSONEqual(t, want, got)
}

func TestMergeAs(t *testing.T) {
	a := []byte(`{"a":1}`)
	b := []byte(`{"b":[1]}`)
	c := []byte(`{"b":[2]}`)
	want := []byte(`{"a":1,"b":[1,2]}`)
	got, err := jsons.MergeAs(jsons.FormatJSON, a, b, c)
	if err != nil {
		t.Error(err)
	}
	assertJSONEqual(t, want, got)
}

func TestRegisterYAML(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterLoader(
		"yaml",
		[]string{".yaml", ".yml"},
		func(b []byte) (map[string]interface{}, error) {
			m := make(map[string]interface{})
			err := yaml.Unmarshal(b, &m)
			if err != nil {
				return nil, err
			}
			return m, nil
		},
	)
	a := []byte(`a: 1`)
	b := []byte(`b: 1`)
	c := []byte(`c: 1`)
	want := []byte(`{"a":1,"b":1,"c":1}`)
	got, err := m.MergeAs("yaml", a, b, c)
	if err != nil {
		t.Error(err)
	}
	assertJSONEqual(t, want, got)
}

func assertJSONEqual(t *testing.T, want, got []byte) {
	wantMap := make(map[string]interface{})
	gotMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(want), &wantMap)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal([]byte(got), &gotMap)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want:\n%s\n\ngot:\n%s", want, got)
	}
}

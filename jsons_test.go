package jsons_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/qjebbs/go-jsons"
	"github.com/qjebbs/go-jsons/merge"
	"gopkg.in/yaml.v2"
)

func TestMergeBytes(t *testing.T) {
	a := []byte(`{"a":1}`)
	b := []byte(`{"b":1}`)
	c := []byte(`{"c":1}`)
	want := []byte(`{"a":1,"b":1,"c":1}`)
	got, err := jsons.Merge(a, b, c)
	if err != nil {
		t.Error(err)
	}
	assertJSONEqual(t, want, got)
}

func TestMergeBytesAs(t *testing.T) {
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
	jsons.Register(
		"yaml",
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
	a := []byte(`a: 1`)
	b := []byte(`b: 1`)
	c := []byte(`c: 1`)
	want := []byte(`{"a":1,"b":1,"c":1}`)
	got, err := jsons.MergeAs("yaml", a, b, c)
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

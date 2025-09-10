package jsons_test

import (
	"encoding/json"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/qjebbs/go-jsons"
)

func TestMerge(t *testing.T) {
	a := []byte(`{"a":1}`)
	b := [][]byte{
		[]byte(`{"b":1}`),
	}
	c := strings.NewReader(`{"c":1}`)
	d := []io.Reader{
		strings.NewReader(`{"d":1}`),
	}
	want := []byte(`{"a":1,"b":1,"c":1,"d":1}`)
	got, err := jsons.Merge(a, b, c, d)
	if err != nil {
		t.Fatal(err)
	}
	assertJSONEqual(t, want, got)
}

func TestMergeComplex(t *testing.T) {
	a := []byte(`
	{
	  	"array_1": [{
			"tag":"1",
			"array_2": [{
				"tag":"2",
				"array_3.1": ["string",true,false],
				"array_3.2": [1,2,3],
				"number_1": 1,
				"number_2": 1,
				"bool_1": true,
				"bool_2": true
			}]
		}],
		"obj_1": {
			"array_4": [
				{ "order": 3 },
				{ "order": 2 },
				{ "order": 1 }
			]
		}
	}
`)
	b := []byte(`
	{
		"array_1": [{
			"tag":"1",
			"array_2": [{
				"tag":"2",
				"array_3.1": [0,1,null],
				"array_3.2": null,
				"number_1": 0,
				"number_2": 1,
				"bool_1": true,
				"bool_2": false,
				"null_1": null
			}]
		},{
			"tag":"2",
			"order": -1
		}]
	}
`)
	want := []byte(`
	{
	  "array_1": [{
			"tag":"2",
			"order": -1
		},{
		"tag":"1",
		"array_2": [{
			"tag":"2",
			"array_3.1": ["string",true,false,0,1,null],
			"array_3.2": [1,2,3],
			"number_1": 0,
			"number_2": 1,
			"bool_1": true,
			"bool_2": false,
			"null_1": null
		}]
	  }],
	  "obj_1": {
		  "array_4": [
			  { "order": 1 },
			  { "order": 2 },
			  { "order": 3 }
		  ]
	  }
	}
	`)
	m := jsons.NewMerger(
		jsons.WithMergeBy("tag"),
		jsons.WithOrderBy("order"),
	)
	got, err := m.MergeAs(jsons.FormatJSON, a, b)
	if err != nil {
		t.Fatal(err)
	}
	assertJSONEqual(t, want, got)
}

func TestMergeFiles(t *testing.T) {
	f, err := os.CreateTemp("", "jsons-test-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	_, err = f.Write([]byte(`{"a":1}`))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	files := []string{f.Name()}
	bytes := []byte(`{"b":1}`)
	got, err := jsons.Merge(files, bytes)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte(`{"a":1,"b":1}`)
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
	if !reflect.DeepEqual(wantMap, gotMap) {
		t.Errorf("want:\n%s\n\ngot:\n%s", want, got)
	}
}

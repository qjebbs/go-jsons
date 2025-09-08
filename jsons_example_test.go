package jsons_test

import (
	"fmt"
	"io"
	"strings"

	"github.com/qjebbs/go-jsons"
)

func ExampleMerger_Merge() {
	var m = jsons.NewMerger(
		jsons.WithTypeOverride(true),
		// jsons.WithMergeBy("tag"),
		// jsons.WithMergeByAndRemove("_tag"),
		// jsons.WithOrderByAndRemove("_order"),
	)
	a := []byte(`{"a":1}`)
	b := [][]byte{
		[]byte(`{"b":1}`),
	}
	c := strings.NewReader(`{"c":false}`)
	d := []io.Reader{
		strings.NewReader(`{"c":1}`),
	}
	got, err := m.Merge(a, b, c, d)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(got))
	// Output: {"a":1,"b":1,"c":1}
}

package jsons_test

import (
	"fmt"

	"github.com/qjebbs/go-jsons"
)

func Example_merge() {
	a := []byte(`{"a":1}`)
	b := []byte(`{"b":[1]}`)
	c := []byte(`{"b":[2]}`)
	d := []byte(`{"c":null}`)

	got, err := jsons.Merge(a, b, c, d)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(got))
	// Output: {"a":1,"b":[1,2],"c":null}
}

func Example_mergeAdvanced() {
	a := []byte(`
	{
      "log": {"level": "debug"},
      "inbounds": [{"tag": "in-1"}],
      "outbounds": [{"_order": 100, "tag": "out-1"}],
      "route": {"rules": [
        {"_tag":"rule1","inbound":["in-1"],"outbound":"out-1"}
	  ]}
	}`)
	b := []byte(`
	{
      "log": {"level": "error"},
      "outbounds": [{"_order": -100, "tag": "out-2"}],
      "route": {"rules": [
        {"_tag":"rule1","inbound":["in-1.1"],"outbound":"out-1.1"}
      ]}
    }`)

	var m = jsons.NewMerger(
		jsons.WithMergeBy("tag"),
		jsons.WithMergeByAndRemove("_tag"),
		jsons.WithOrderByAndRemove("_order"),
		jsons.WithIndent("", "  "),
	)
	got, err := m.Merge(a, b)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(got))
	// Output: {
	//   "log": {
	//     "level": "error"
	//   },
	//   "inbounds": [
	//     {
	//       "tag": "in-1"
	//     }
	//   ],
	//   "outbounds": [
	//     {
	//       "tag": "out-2"
	//     },
	//     {
	//       "tag": "out-1"
	//     }
	//   ],
	//   "route": {
	//     "rules": [
	//       {
	//         "inbound": [
	//           "in-1",
	//           "in-1.1"
	//         ],
	//         "outbound": "out-1.1"
	//       }
	//     ]
	//   }
	// }
}

func ExampleMerger_RegisterOrderedLoader() {
	const FormatYAML jsons.Format = "yaml"
	m := jsons.NewMerger()
	m.RegisterOrderedLoader(
		FormatYAML,
		[]string{".yaml", ".yml"},
		func(b []byte) (*jsons.OrderedMap, error) {
			m := jsons.NewOrderedMap()
			// "github.com/goccy/go-yaml" is recommended since it's able to use json.Unmarshaler
			// err := yaml.UnmarshalWithOptions(
			// 	b, m,
			// 	yaml.UseJSONUnmarshaler(), // important
			// )
			// if err != nil {
			// 	return nil, err
			// }
			return m, nil
		},
	)
}

func ExampleMerger_RegisterLoader() {
	const FormatTOML jsons.Format = "toml"
	m := jsons.NewMerger()
	m.RegisterLoader(
		FormatTOML,
		[]string{".toml"},
		func(b []byte) (map[string]interface{}, error) {
			m := make(map[string]interface{})
			// err := toml.Unmarshal(b, &m)
			// if err != nil {
			// 	return nil, err
			// }
			return m, nil
		},
	)
}
func ExampleWithTypeOverride() {
	a := []byte(`{"a":1}`)
	b := []byte(`{"a":false}`)

	m := jsons.NewMerger(
		jsons.WithTypeOverride(true),
	)
	got, err := m.Merge(a, b)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(got))
	// Output: {"a":false}
}

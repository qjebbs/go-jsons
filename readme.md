# go-jsons

A universal `JSON` merge library for `Go`. 

(test coverage: 100.0%)

## Installation

```bash
go get github.com/qjebbs/go-jsons
```

## Usage

```go
a := []byte(`{"a":1}`)
b := []byte(`{"b":[1]}`)
c := []byte(`{"b":[2]}`)
got, err := jsons.Merge(a, b, c) // got = []byte(`{"a":1,"b":[1,2]}`)
```

### Accepted input

- `string`: path to a local file
- `[]string`: paths of local files
- `[]byte`: content of a file
- `[][]byte`: content list of files
- `io.Reader`: content reader
- `[]io.Reader`: content readers

## Merge rules

The strandard merger is intuitive and easy to understand:

- Simple values (`string`, `number`, `boolean`) are overwritten by later ones.
- Container values (`object`, `array`) are merged recursively.

To work with complex contents, you may create a custom merger to applies more options:

```go
var myMerger = jsons.NewMerger(
	jsons.WithMergeBy("tag"),
	jsons.WithMergeByAndRemove("_tag"),
	jsons.WithOrderByAndRemove("_order"),
)
myMerger.Merge("a.json", "b.json")
```

which means:

- Elements with same `tag` or `_tag` in an array will be merged.
- Elements in an array will be sorted by the value of `_order` field, the smaller ones are in front.

> `_tag` and `_order` fields will be removed after merge, according to the codes above.

Suppose we have...

`a.json`:

```json
{
  "log": {"level": "debug"},
  "inbounds": [{"tag": "in-1"}],
  "outbounds": [{"_order": 100, "tag": "out-1"}],
  "route": {"rules": [
    {"_tag":"rule1","inbound":["in-1"],"outbound":"out-1"}
  ]}
}
```

`b.json`:

```json
{
  "log": {"level": "error"},
  "outbounds": [{"_order": -100, "tag": "out-2"}],
  "route": {"rules": [
    {"_tag":"rule1","inbound":["in-1.1"],"outbound":"out-1.1"}
  ]}
}
```

Output:

```jsonc
{
  // level field is overwritten by the latter value
  "log": {"level": "error"},
  "inbounds": [{"tag": "in-1"}],
  "outbounds": [
    // Although out-2 is a latecomer, but it's in 
    // the front due to the smaller "_order"
    {"tag": "out-2"},
    {"tag": "out-1"}
  ],
  "route": {"rules": [
    // 2 rules are merged into one due to the same "_tag",
    // outbound field is overwritten during the merging
    {"inbound":["in-1","in-1.1"],"outbound":"out-1.1"}
  ]}
}
```

## Load from other formats

`go-jsons` allows you to extend it to load other formats easily.

For example, to load from `YAML` files and merge to `JSON`:

```go
package main

import (
	"fmt"

	"github.com/qjebbs/go-jsons"
	// goccy/go-yaml is able to use json.Unmarshaler
	"github.com/goccy/go-yaml"
)

func ExampleMerger_RegisterLoader() {
	const FormatYAML jsons.Format = "yaml"
	m := jsons.NewMerger()
	m.RegisterOrderedLoader(
		FormatYAML,
		[]string{".yaml", ".yml"},
		func(b []byte) (*jsons.OrderedMap, error) {
			// YAML fields order will be kept
			m := jsons.NewOrderedMap()
			err := yaml.UnmarshalWithOptions(
				b, m,
				// important
				yaml.UseJSONUnmarshaler(),
			)
			if err != nil {
				return nil, err
			}
			return m, nil
		},
	)
	a := []byte(`{"a":1,"z":1}`)    // json
	b := []byte("b: 1\nc: 1\nd: 1") // yaml
	got, err := m.Merge(a, b)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(got))
	// Output: {"a":1,"z":1,"b":1,"c":1,"d":1}
}
```

## Why not support remote files?

Here are some considerations:

- It makes the your program support remote file unexpectedly, which may be a security risk.
- Users need to choose their own strategy for loading remote files, not hard-coded logic in the library
- You can still merge downloaded content by `[]byte` or `io.Reader`
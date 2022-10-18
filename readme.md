# go-jsons

A universal `JSON` merge library for `Go`.

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

### Accepted Input

- `[]byte`: content of a file
- `string`: path to a file, either local or remote
- `[]string`: a list of files, either local or remote
- `io.Reader`: a file content reader

## Merge Rules

The merge logic is intuitive and easy to understand:

- Simple values (`string`, `number`, `boolean`) are overwritten, others (`array`, `object`) are merged.

To work with complex files and contents, especially when the order of fields of output matters, we can add `_tag` and `_priority` fields to apply more merge rules:

- Elements with same `_tag` in an array will be merged.
- Elements in an array will be sorted by `_priority` field, the smaller the higher priority.

Suppose we have 2 `JSON` files:

`a.json`:

```json
{
  "log": {"level": "debug"},
  "inbounds": [{"tag": "in-1"}],
  "outbounds": [{"_priority": 100, "tag": "out-1"}],
  "route": {"rules": [
    {"_tag":"rule1","inbound":["in-1"],"outbound":"out-1"}
  ]}
}
```

`b.json`:

```json
{
  "log": {"level": "error"},
  "outbounds": [{"_priority": -100, "tag": "out-2"}],
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
    // the front due to the smaller "_priority"
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

## Other Formats Support

`go-jsons` supports only `JSON`, but it allows you to extend it to support other formats easily.

For example, to support merge from `YAML` files to `JSON`:

```go
package main

import (
	"fmt"

	"github.com/qjebbs/go-jsons"
	"github.com/qjebbs/go-jsons/merge"
	"gopkg.in/yaml.v2"
)

func main() {
	const FormatYAML jsons.Format = "yaml"
	jsons.Register(
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
	a := []byte(`a: 1`)
	b := []byte(`b: 1`)
	got, err := jsons.Merge(a, b)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(got)) // {"a":1,"b":1}
}
```
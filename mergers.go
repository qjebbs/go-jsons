package jsons

import (
	"encoding/json"
	"fmt"
	"strings"
)

func init() {
	must(Register(
		FormatJSON,
		[]string{".json"},
		func(v []byte) (map[string]interface{}, error) {
			m := make(map[string]interface{})
			if err := json.Unmarshal(v, &m); err != nil {
				return nil, err
			}
			return m, nil
		},
	))
}

var (
	mergersByName = make(map[Format]*Merger)
	mergersByExt  = make(map[string]*Merger)
)

// Register register a new format.
func Register(name Format, extensions []string, converter ConvertFunc) error {
	return registerMerger(makeMerger(name, extensions, converter))
}

// Unregister unregister a format.
func Unregister(name Format) {
	format, found := mergersByName[name]
	if !found {
		return
	}
	delete(mergersByName, name)
	for _, ext := range format.Extensions {
		lext := strings.ToLower(ext)
		if _, found := mergersByExt[lext]; found {
			delete(mergersByExt, lext)
		}
	}
}

// registerMerger add a new Merger.
func registerMerger(format *Merger) error {
	if _, found := mergersByName[format.Name]; found {
		return fmt.Errorf("%s already registered", format.Name)
	}
	mergersByName[format.Name] = format
	for _, ext := range format.Extensions {
		lext := strings.ToLower(ext)
		if f, found := mergersByExt[lext]; found {
			return fmt.Errorf("%s already registered to %s", ext, f.Name)
		}
		mergersByExt[lext] = format
	}
	return nil
}

func must(errs ...error) {
	for _, err := range errs {
		if err != nil {
			panic(err)
		}
	}
}

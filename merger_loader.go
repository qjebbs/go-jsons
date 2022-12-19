package jsons

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RegisterLoader register a new format loader.
func (m *Merger) RegisterLoader(name Format, extensions []string, fn LoadFunc) error {
	if _, found := m.loadersByName[name]; found {
		return fmt.Errorf("%s already registered", name)
	}
	loader := newLoader(name, extensions, fn)
	m.loadersByName[name] = loader
	for _, ext := range extensions {
		lext := strings.ToLower(ext)
		if f, found := m.loadersByExt[lext]; found {
			return fmt.Errorf("%s already registered to %s", ext, f.Name)
		}
		m.loadersByExt[lext] = loader
	}
	return nil
}

// RegisterDefaultLoader register the default json loader.
func (m *Merger) RegisterDefaultLoader() error {
	return m.RegisterLoader(
		FormatJSON,
		[]string{".json"},
		func(v []byte) (map[string]interface{}, error) {
			m := make(map[string]interface{})
			if err := json.Unmarshal(v, &m); err != nil {
				return nil, err
			}
			return m, nil
		},
	)
}

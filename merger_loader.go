package jsons

import (
	"fmt"
	"strings"

	"github.com/qjebbs/go-jsons/internal/ordered"
)

// RegisterOrderedLoader register a new format loader that loads into an ordered map,
// who keeps the fields order between merges.
func (m *Merger) RegisterOrderedLoader(name Format, extensions []string, fn LoadOrderedFunc) error {
	if name == FormatAuto {
		return fmt.Errorf("cannot register with reserved name: '%s'", FormatAuto)
	}
	if old, found := m.loadersByName[name]; found {
		for _, format := range old.Extensions {
			delete(m.loadersByExt, format)
		}
	}
	loader := newLoader(name, extensions, fn)
	m.loadersByName[name] = loader
	for _, ext := range extensions {
		lext := strings.ToLower(ext)
		if f, found := m.loadersByExt[lext]; found {
			return fmt.Errorf("file extension '%s' is already registered to '%s'", ext, f.Name)
		}
		m.loadersByExt[lext] = loader
	}
	return nil
}

// RegisterLoader register a new format loader.
// The fields order is not guaranteed between merges.
func (m *Merger) RegisterLoader(name Format, extensions []string, fn LoadFunc) error {
	fn2 := func(b []byte) (*ordered.Map, error) {
		m, err := fn(b)
		if err != nil {
			return nil, err
		}
		return ordered.FromMap(m), nil
	}
	return m.RegisterOrderedLoader(name, extensions, fn2)
}

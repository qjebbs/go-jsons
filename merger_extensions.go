package jsons

import (
	"fmt"
)

// GetExtensions get extensions of given format
func (m *Merger) GetExtensions(formatName Format) ([]string, error) {
	if formatName == FormatAuto {
		return m.GetAllExtensions(), nil
	}
	f, found := m.loadersByName[formatName]
	if !found {
		return nil, fmt.Errorf("%s not found", formatName)
	}
	return f.Extensions, nil
}

// GetAllExtensions get all extensions supported
func (m *Merger) GetAllExtensions() []string {
	extensions := make([]string, 0)
	for _, f := range m.loadersByName {
		extensions = append(extensions, f.Extensions...)
	}
	return extensions
}

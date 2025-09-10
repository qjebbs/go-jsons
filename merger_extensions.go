package jsons

import "fmt"

// Extensions get supported extensions of given formats.
// If formatNames is empty or contains FormatAuto, it returns all extensions.
func (m *Merger) Extensions(formatNames ...Format) ([]string, error) {
	if len(formatNames) == 0 || contains(formatNames, FormatAuto) {
		return m.getAllExtensions(), nil
	}
	var extensions []string
	seen := make(map[string]struct{})
	for _, formatName := range formatNames {
		f, found := m.loadersByName[formatName]
		if !found {
			return nil, fmt.Errorf("%s not found", formatName)
		}
		for _, ext := range f.Extensions {
			if _, found := seen[ext]; !found {
				extensions = append(extensions, ext)
				seen[ext] = struct{}{}
			}
		}
	}
	return extensions, nil
}

// getAllExtensions get all extensions supported
func (m *Merger) getAllExtensions() []string {
	extensions := make([]string, 0)
	for ext := range m.loadersByExt {
		extensions = append(extensions, ext)
	}
	return extensions
}

func contains[T comparable](slice []T, str T) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

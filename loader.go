package jsons

import (
	"fmt"
	"io"
	"os"

	"github.com/qjebbs/go-jsons/internal/ordered"
)

// OrderedMap is an alias of ordered.Map
type OrderedMap = ordered.Map

// NewOrderedMap is an alias of ordered.New
var NewOrderedMap = ordered.New

// LoadFunc load the input bytes to map[string]interface{}
type LoadFunc func([]byte) (map[string]interface{}, error)

// LoadOrderedFunc load the input bytes to *OrderedMap, which keeps the fields order
type LoadOrderedFunc func([]byte) (*OrderedMap, error)

// loader is a configurable loader for specific format files.
type loader struct {
	Name       Format
	Extensions []string
	LoadFunc   LoadOrderedFunc
}

// makeLoader makes a merger who merge the format by converting it to JSON
func newLoader(name Format, extensions []string, fn LoadOrderedFunc) *loader {
	return &loader{
		Name:       name,
		Extensions: extensions,
		LoadFunc:   fn,
	}
}

// makeLoadFunc makes a merge func who merge the input to
func (l *loader) Load(input interface{}) ([]*OrderedMap, error) {
	if input == nil {
		return nil, nil
	}
	switch v := input.(type) {
	case string:
		return l.loadFiles([]string{v})
	case []string:
		return l.loadFiles(v)
	case []byte:
		return l.loadSlices([][]byte{v})
	case [][]byte:
		return l.loadSlices(v)
	case io.Reader:
		return l.loadReaders([]io.Reader{v})
	case []io.Reader:
		return l.loadReaders(v)
	default:
		return nil, fmt.Errorf("unsupported input type: %T", input)
	}
}

func (l *loader) loadFiles(files []string) ([]*OrderedMap, error) {
	maps := make([]*OrderedMap, 0, len(files))
	for _, file := range files {
		m, err := l.loadFile(file)
		if err != nil {
			return nil, err
		}
		maps = append(maps, m)
	}
	return maps, nil
}

func (l *loader) loadReaders(readers []io.Reader) ([]*OrderedMap, error) {
	maps := make([]*OrderedMap, 0, len(readers))
	for _, r := range readers {
		m, err := l.loadReader(r)
		if err != nil {
			return nil, err
		}
		maps = append(maps, m)
	}
	return maps, nil
}

func (l *loader) loadSlices(slices [][]byte) ([]*OrderedMap, error) {
	maps := make([]*OrderedMap, 0, len(slices))
	for _, slice := range slices {
		m, err := l.LoadFunc(slice)
		if err != nil {
			return nil, err
		}
		maps = append(maps, m)
	}
	return maps, nil
}

func (l *loader) loadFile(file string) (*OrderedMap, error) {
	bs, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return l.LoadFunc(bs)
}

func (l *loader) loadReader(reader io.Reader) (*OrderedMap, error) {
	bs, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return l.LoadFunc(bs)
}

package jsons

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/qjebbs/go-jsons/internal/merge"
	"github.com/qjebbs/go-jsons/internal/ordered"
)

// Merger is the json merger
type Merger struct {
	loadersByName map[Format]*loader
	loadersByExt  map[string]*loader
	options       *Options
}

// NewMerger returns a new Merger
func NewMerger(options ...Option) *Merger {
	o := &Options{}
	for _, opt := range options {
		opt(o)
	}
	m := &Merger{
		loadersByName: make(map[Format]*loader),
		loadersByExt:  make(map[string]*loader),
		options:       o,
	}
	// never return error
	_ = m.RegisterOrderedLoader(
		FormatJSON,
		[]string{".json"},
		func(v []byte) (*ordered.Map, error) {
			m := ordered.New()
			if err := json.Unmarshal(v, m); err != nil {
				return nil, err
			}
			return m, nil
		},
	)
	return m
}

// Merge merges inputs into a single json.
//
// It detects the format by file extension, or try all mergers
// if no extension found
//
// Accepted Input:
//
//   - string: path to a local file
//   - []string: paths of local files
//   - []byte: content of a file
//   - [][]byte: content list of files
//   - io.Reader: content reader
//   - []io.Reader: content readers
func (m *Merger) Merge(inputs ...interface{}) ([]byte, error) {
	target := ordered.New()
	for _, input := range inputs {
		err := m.mergeToMap(input, target)
		if err != nil {
			return nil, err
		}
	}
	err := m.options.apply(target)
	if err != nil {
		return nil, err
	}
	if m.options.MarshalIndent != "" {
		return json.MarshalIndent(target, m.options.MarshalPrefix, m.options.MarshalIndent)
	}
	return json.Marshal(target)
}

// MergeAs loads inputs of the specific format and merges into a single json.
//
// Accepted Input:
//
//   - string: path to a local file
//   - []string: paths of local files
//   - []byte: content of a file
//   - [][]byte: content list of files
//   - io.Reader: content reader
//   - []io.Reader: content readers
func (m *Merger) MergeAs(format Format, inputs ...interface{}) ([]byte, error) {
	target := ordered.New()
	for _, input := range inputs {
		err := m.mergeToMapAs(format, input, target)
		if err != nil {
			return nil, err
		}
	}
	err := m.options.apply(target)
	if err != nil {
		return nil, err
	}
	if m.options.MarshalIndent != "" {
		return json.MarshalIndent(target, m.options.MarshalPrefix, m.options.MarshalIndent)
	}
	return json.Marshal(target)
}

func (m *Merger) mergeToMapAs(formatName Format, input interface{}, target *ordered.Map) error {
	if formatName == FormatAuto {
		return m.mergeToMap(input, target)
	}
	f, found := m.loadersByName[formatName]
	if !found {
		return fmt.Errorf("unknown format: %s", formatName)
	}
	maps, err := f.Load(input)
	if err != nil {
		return err
	}
	return merge.OrderedMaps(m.options.TypeOverride, target, maps...)
}

func (m *Merger) mergeToMap(input interface{}, target *ordered.Map) error {
	if input == nil {
		return nil
	}
	switch v := input.(type) {
	case string:
		// load by file extension
		if ext := getExtension(v); ext != "" {
			lext := strings.ToLower(ext)
			if f, found := m.loadersByExt[lext]; found {
				mp, err := f.Load(v)
				if err != nil {
					return err
				}
				return merge.OrderedMaps(m.options.TypeOverride, target, mp...)
			}
		}
		err := m.tryLoaders(v, target)
		if err != nil {
			return err
		}
	case io.Reader:
		// read into []byte in case it's drained when try different load
		bs, err := io.ReadAll(v)
		if err != nil {
			return err
		}
		err = m.tryLoaders(bs, target)
		if err != nil {
			return err
		}
	case []string:
		for _, v := range v {
			err := m.mergeToMap(v, target)
			if err != nil {
				return err
			}
		}
	case []io.Reader:
		for _, v := range v {
			err := m.mergeToMap(v, target)
			if err != nil {
				return err
			}
		}
	default:
		return m.tryLoaders(v, target)
	}
	return nil
}

func (m *Merger) tryLoaders(input interface{}, target *ordered.Map) error {
	var errs []string
	for _, f := range m.loadersByName {
		mp, err := f.Load(input)
		if err == nil {
			return merge.OrderedMaps(m.options.TypeOverride, target, mp...)
		}
		errs = append(errs, fmt.Sprintf("[%s] %s", f.Name, err))
	}
	return fmt.Errorf("tried all formats but failed: %s", strings.Join(errs, "; "))
}

func getExtension(filename string) string {
	ext := filepath.Ext(filename)
	return strings.ToLower(ext)
}

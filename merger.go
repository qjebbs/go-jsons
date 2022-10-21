package jsons

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/qjebbs/go-jsons/merge"
)

// Merger is the json merger
type Merger struct {
	loadersByName map[Format]*loader
	loadersByExt  map[string]*loader
}

// NewMerger returns a new Merger
func NewMerger() *Merger {
	return &Merger{
		loadersByName: make(map[Format]*loader),
		loadersByExt:  make(map[string]*loader),
	}
}

// Merge merges inputs into a single json.
//
// It detects the format by file extension, or try all mergers
// if no extension found
//
// Accepted Input:
//
//  - `string`: path to a local file
//  - `[]string`: paths of local files
//  - `[]byte`: content of a file
//  - `[][]byte`: content list of files
//  - `io.Reader`: content reader
//  - `[]io.Reader`: content readers
func (m *Merger) Merge(inputs ...interface{}) ([]byte, error) {
	tmp := make(map[string]interface{})
	for _, input := range inputs {
		err := m.mergeToMap(input, tmp)
		if err != nil {
			return nil, err
		}
	}
	err := merge.ApplyRules(tmp)
	if err != nil {
		return nil, err
	}
	merge.RemoveHelperFields(tmp)
	return json.Marshal(tmp)
}

// MergeAs loads inputs of the specific format and merges into a single json.
//
// Accepted Input:
//
//  - `string`: path to a local file
//  - `[]string`: paths of local files
//  - `[]byte`: content of a file
//  - `[][]byte`: content list of files
//  - `io.Reader`: content reader
//  - `[]io.Reader`: content readers
func (m *Merger) MergeAs(format Format, inputs ...interface{}) ([]byte, error) {
	tmp := make(map[string]interface{})
	for _, input := range inputs {
		err := m.mergeToMapAs(format, input, tmp)
		if err != nil {
			return nil, err
		}
	}
	err := merge.ApplyRules(tmp)
	if err != nil {
		return nil, err
	}
	merge.RemoveHelperFields(tmp)
	return json.Marshal(tmp)
}

// mergeToMapAs load inputs of the specific format into target
//
// Accepted Input:
//
//  - `string`: path to a local file
//  - `[]string`: paths of local files
//  - `[]byte`: content of a file
//  - `[][]byte`: content list of files
//  - `io.Reader`: content reader
//  - `[]io.Reader`: content readers
//
// it will neither apply "_priority" sort or "_tag" merge rules
// nor remove helper fields. You may want to call them manually:
//
//  err := merge.ApplyRules(target)
//  if err != nil {
//  	return nil, err
//  }
//  merge.RemoveHelperFields(target)
func (m *Merger) mergeToMapAs(formatName Format, input interface{}, target map[string]interface{}) error {
	if formatName == FormatAuto {
		return m.mergeToMap(input, target)
	}
	f, found := m.loadersByName[formatName]
	if !found {
		return fmt.Errorf("unknown format: %s", formatName)
	}
	return f.Merge(input, target)
}

// mergeToMap loads inputs and merges them into target.
// It detects the format by file extension, or try all mergers
// if no extension found
//
// Accepted Input:
//
//  - `string`: path to a local file
//  - `[]string`: paths of local files
//  - `[]byte`: content of a file
//  - `[][]byte`: content list of files
//  - `io.Reader`: content reader
//  - `[]io.Reader`: content readers
//
// it will neither apply "_priority" sort or "_tag" merge rules
// nor remove helper fields. You may want to call them manually:
//
//  err := merge.ApplyRules(target)
//  if err != nil {
//  	return nil, err
//  }
//  merge.RemoveHelperFields(target)
func (m *Merger) mergeToMap(input interface{}, target map[string]interface{}) error {
	if input == nil {
		return nil
	}
	switch v := input.(type) {
	case string:
		// try to load by extension
		if ext := getExtension(v); ext != "" {
			lext := strings.ToLower(ext)
			if f, found := m.loadersByExt[lext]; found {
				return f.Merge(v, target)
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

func (m *Merger) tryLoaders(input interface{}, target map[string]interface{}) error {
	var errs []string
	for _, f := range m.loadersByName {
		err := f.Merge(input, target)
		if err == nil {
			return nil
		}
		errs = append(errs, fmt.Sprintf("[%s] %s", f.Name, err))
	}
	return fmt.Errorf("tried all formats but failed: %s", strings.Join(errs, "; "))
}

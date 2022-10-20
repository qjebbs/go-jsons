package jsons

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
//  - `[]byte`: content of a file
//  - `string`: path to a file, either local or remote
//  - `[]string`: a list of files, either local or remote
//  - `io.Reader`: a file content reader
func (m *Merger) Merge(inputs ...interface{}) ([]byte, error) {
	tmp := make(map[string]interface{})
	for _, input := range inputs {
		err := m.MergeToMap(input, tmp)
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
//  - `[]byte`: content of a file
//  - `string`: path to a file, either local or remote
//  - `[]string`: a list of files, either local or remote
//  - `io.Reader`: a file content reader
func (m *Merger) MergeAs(format Format, inputs ...interface{}) ([]byte, error) {
	tmp := make(map[string]interface{})
	for _, input := range inputs {
		err := m.MergeToMapAs(format, input, tmp)
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

// MergeToMapAs load inputs of the specific format into target
//
// Accepted Input:
//
//  - `[]byte`: content of a file
//  - `string`: path to a file, either local or remote
//  - `[]string`: a list of files, either local or remote
//  - `io.Reader`: a file content reader
//
// it will neither apply "_priority" sort or "_tag" merge rules
// nor remove helper fields. You may want to call them manually:
//
//  err := merge.ApplyRules(target)
//  if err != nil {
//  	return nil, err
//  }
//  merge.RemoveHelperFields(target)
func (m *Merger) MergeToMapAs(formatName Format, input interface{}, target map[string]interface{}) error {
	if formatName == FormatAuto {
		return m.MergeToMap(input, target)
	}
	f, found := m.loadersByName[formatName]
	if !found {
		return fmt.Errorf("unknown format: %s", formatName)
	}
	return f.Merge(input, target)
}

// MergeToMap loads inputs and merges them into target.
// It detects the format by file extension, or try all mergers
// if no extension found
//
// Accepted Input:
//
//  - `[]byte`: content of a file
//  - `string`: path to a file, either local or remote
//  - `[]string`: a list of files, either local or remote
//  - `io.Reader`: a file content reader
//
// it will neither apply "_priority" sort or "_tag" merge rules
// nor remove helper fields. You may want to call them manually:
//
//  err := merge.ApplyRules(target)
//  if err != nil {
//  	return nil, err
//  }
//  merge.RemoveHelperFields(target)
func (m *Merger) MergeToMap(input interface{}, target map[string]interface{}) error {
	if input == nil {
		return nil
	}
	switch v := input.(type) {
	case string:
		err := m.mergeContent(v, target)
		if err != nil {
			return err
		}
	case []string:
		for _, file := range v {
			err := m.mergeContent(file, target)
			if err != nil {
				return err
			}
		}
	case []byte:
		err := m.mergeContent(v, target)
		if err != nil {
			return err
		}
	case io.Reader:
		// read to []byte incase it tries different mergers
		bs, err := ioutil.ReadAll(v)
		if err != nil {
			return err
		}
		err = m.mergeContent(bs, target)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported input type: %T", input)
	}
	return nil
}

func (m *Merger) mergeContent(input interface{}, target map[string]interface{}) error {
	if file, ok := input.(string); ok {
		ext := getExtension(file)
		if ext != "" {
			lext := strings.ToLower(ext)
			f, found := m.loadersByExt[lext]
			if !found {
				return fmt.Errorf("unsupported file extension: %s", ext)
			}
			return f.Merge(file, target)
		}
	}
	var errs []string
	// no extension, try all mergers
	for _, f := range m.loadersByName {
		err := f.Merge(input, target)
		if err == nil {
			return nil
		}
		errs = append(errs, fmt.Sprintf("[%s] %s", f.Name, err))
	}
	return fmt.Errorf("tried all formats but failed for: \n\n%s\n\nerrors:\n\n  %s", input, strings.Join(errs, "\n  "))
}

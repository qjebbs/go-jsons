package jsons

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/qjebbs/go-jsons/merge"
)

// Merge loads inputs and merges them into json in []byte.
//
// It detects the file extension for format merger selecting, or try all mergers
// if no extension found
//
// Accepted Input:
//
//  - `[]byte`: content of a file
//  - `string`: path to a file, either local or remote
//  - `[]string`: a list of files, either local or remote
//  - `io.Reader`: a file content reader
//
// It will apply "_priority" sort and "_tag" merge
// rules, and remove helper fields before the final output
func Merge(inputs ...interface{}) ([]byte, error) {
	m := make(map[string]interface{})
	for _, input := range inputs {
		err := MergeToMap(input, m)
		if err != nil {
			return nil, err
		}
	}
	err := merge.ApplyRules(m)
	if err != nil {
		return nil, err
	}
	merge.RemoveHelperFields(m)
	return json.Marshal(m)
}

// MergeAs loads inputs as 'format' and merges them into json in []byte
//
// Accepted Input:
//
//  - `[]byte`: content of a file
//  - `string`: path to a file, either local or remote
//  - `[]string`: a list of files, either local or remote
//  - `io.Reader`: a file content reader
//
//
// It will apply "_priority" sort and "_tag" merge
// rules, and remove helper fields before the final output
func MergeAs(format Format, inputs ...interface{}) ([]byte, error) {
	m := make(map[string]interface{})
	for _, input := range inputs {
		err := MergeToMapAs(format, input, m)
		if err != nil {
			return nil, err
		}
	}
	err := merge.ApplyRules(m)
	if err != nil {
		return nil, err
	}
	merge.RemoveHelperFields(m)
	return json.Marshal(m)
}

// MergeToMapAs load input and merge as specified format into m
//
// Accepted Input:
//
//  - `[]byte`: content of a file
//  - `string`: path to a file, either local or remote
//  - `[]string`: a list of files, either local or remote
//  - `io.Reader`: a file content reader
//
//
// Note: it will neither apply "_priority" sort or "_tag" merge
// rules nor remove helper fields
func MergeToMapAs(formatName Format, input interface{}, target map[string]interface{}) error {
	f, found := mergersByName[formatName]
	if !found {
		return fmt.Errorf("unknown format: %s", formatName)
	}
	return f.Merge(input, target)
}

// MergeToMap loads inputs and merges them into target.
// It detects the file extension for format merger selecting, or try all mergers
// if no extension found
//
// Accepted Input:
//
//  - `[]byte`: content of a file
//  - `string`: path to a file, either local or remote
//  - `[]string`: a list of files, either local or remote
//  - `io.Reader`: a file content reader
//
//
// Note: it will neither apply "_priority" sort or "_tag" merge
// rules nor remove helper fields
func MergeToMap(input interface{}, target map[string]interface{}) error {
	switch v := input.(type) {
	case string:
		err := mergeSingleFile(v, target)
		if err != nil {
			return err
		}
	case []string:
		for _, file := range v {
			err := mergeSingleFile(file, target)
			if err != nil {
				return err
			}
		}
	case []byte:
		err := mergeSingleFile(v, target)
		if err != nil {
			return err
		}
	case io.Reader:
		// read to []byte incase it tries different mergers
		bs, err := ioutil.ReadAll(v)
		if err != nil {
			return err
		}
		err = mergeSingleFile(bs, target)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknow supported input type: %T", input)
	}
	return nil
}

func mergeSingleFile(input interface{}, m map[string]interface{}) error {
	if file, ok := input.(string); ok {
		ext := getExtension(file)
		if ext != "" {
			lext := strings.ToLower(ext)
			f, found := mergersByExt[lext]
			if !found {
				return fmt.Errorf("unsupported file extension: %s", ext)
			}
			return f.Merge(file, m)
		}
	}
	var errs []string
	// no extension, try all mergers
	for _, f := range mergersByName {
		if f.Name == FormatAuto {
			continue
		}
		err := f.Merge(input, m)
		if err == nil {
			return nil
		}
		errs = append(errs, fmt.Sprintf("[%s] %s", f.Name, err))
	}
	return fmt.Errorf("tried all formats but failed for: \n\n%s\n\nerrors:\n\n  %s", input, strings.Join(errs, "\n  "))
}

func getExtension(filename string) string {
	ext := filepath.Ext(filename)
	return strings.ToLower(ext)
}

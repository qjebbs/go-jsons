package jsons

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/qjebbs/go-jsons/merge"
)

// LoadFunc load the input bytes to map[string]interface{}
type LoadFunc func([]byte) (map[string]interface{}, error)

// loader is a configurable loader for specific format files.
type loader struct {
	Name       Format
	Extensions []string
	Merge      loadFunc
}

// loadFunc is a function to load the input into map[string]interface{}
type loadFunc func(input interface{}, target map[string]interface{}) error

// makeLoader makes a merger who merge the format by converting it to JSON
func makeLoader(name Format, extensions []string, converter LoadFunc) *loader {
	return &loader{
		Name:       name,
		Extensions: extensions,
		Merge:      makeLoadFunc(converter),
	}
}

// makeLoadFunc makes a merge func who merge the input to
func makeLoadFunc(converter LoadFunc) loadFunc {
	return func(input interface{}, target map[string]interface{}) error {
		if target == nil {
			return errors.New("target is nil")
		}
		if input == nil {
			return nil
		}
		switch v := input.(type) {
		case string:
			err := loadFile(v, target, converter)
			if err != nil {
				return err
			}
		case []string:
			err := loadFiles(v, target, converter)
			if err != nil {
				return err
			}
		case []byte:
			err := loadBytes(v, target, converter)
			if err != nil {
				return err
			}
		case io.Reader:
			err := loadReader(v, target, converter)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported input type: %T", input)
		}
		return nil
	}
}

func loadFiles(files []string, target map[string]interface{}, converter LoadFunc) error {
	for _, file := range files {
		err := loadFile(file, target, converter)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadFile(file string, target map[string]interface{}, converter LoadFunc) error {
	bs, err := loadToBytes(file)
	if err != nil {
		return fmt.Errorf("fail to load %s: %s", file, err)
	}
	return loadBytes(bs, target, converter)
}

func loadReader(reader io.Reader, target map[string]interface{}, converter LoadFunc) error {
	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return loadBytes(bs, target, converter)
}

func loadBytes(bs []byte, target map[string]interface{}, converter LoadFunc) error {
	m, err := converter(bs)
	if err != nil {
		return err
	}
	return merge.Maps(target, m)
}

package jsons_test

import (
	"errors"
	"testing"

	"github.com/qjebbs/go-jsons"
)

func TestLoadToNilTarget(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterDefaultLoader()
	err := m.MergeToMap([]byte("{}"), nil)
	if err == nil {
		t.Error("want error, got nil")
	}
}
func TestLoadToNilTarget2(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterDefaultLoader()
	err := m.MergeToMapAs(jsons.FormatJSON, []byte("{}"), nil)
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestLoadUnsupportInput(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterDefaultLoader()
	target := make(map[string]interface{})
	err := m.MergeToMap(true, target)
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestLoadUnsupportInput2(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterDefaultLoader()
	target := make(map[string]interface{})
	err := m.MergeToMapAs(jsons.FormatJSON, true, target)
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestLoadNilInput(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterDefaultLoader()
	target := make(map[string]interface{})
	err := m.MergeToMap(nil, target)
	if err != nil {
		t.Errorf("want nil, got: %s", err)
	}
}

func TestLoadNilInput2(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterDefaultLoader()
	target := make(map[string]interface{})
	err := m.MergeToMapAs(jsons.FormatJSON, nil, target)
	if err != nil {
		t.Errorf("want nil, got: %s", err)
	}
}

func TestLoadReaderError(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterDefaultLoader()

	target := make(map[string]interface{})
	err := m.MergeToMapAs(jsons.FormatJSON, &errReader{}, target)
	if err == nil {
		t.Error("want error, got nil")
	}
}

type errReader struct{}

func (r *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error")
}

package jsons_test

import (
	"errors"
	"testing"

	"github.com/qjebbs/go-jsons"
)

func TestLoadUnsupportInput(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterDefaultLoader()
	_, err := m.Merge(true)
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestLoadUnsupportInput2(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterDefaultLoader()
	_, err := m.MergeAs(jsons.FormatAuto, true)
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestLoadNilInput(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterDefaultLoader()
	_, err := m.Merge(nil)
	if err != nil {
		t.Errorf("want nil, got: %s", err)
	}
}

func TestLoadNilInput2(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterDefaultLoader()
	_, err := m.MergeAs(jsons.FormatJSON, nil)
	if err != nil {
		t.Errorf("want nil, got: %s", err)
	}
}

func TestLoadReaderError(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterDefaultLoader()
	_, err := m.MergeAs(jsons.FormatJSON, &errReader{})
	if err == nil {
		t.Error("want error, got nil")
	}
}

type errReader struct{}

func (r *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error")
}

package jsons_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/qjebbs/go-jsons"
)

func TestMergeAsUnknownFormat(t *testing.T) {
	_, err := jsons.MergeAs("unknown", []byte(`{}`))
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestMergeFileError(t *testing.T) {
	_, err := jsons.Merge("file_not_exist.json")
	if err == nil {
		t.Error("want error, got nil")
	}
	_, err = jsons.Merge([]string{"file_not_exist.unknown"})
	if err == nil {
		t.Error("want error, got nil")
	}
	_, err = jsons.Merge([]string{"file_not_exist.json"})
	if err == nil {
		t.Error("want error, got nil")
	}
	_, err = jsons.MergeAs(jsons.FormatJSON, "file_not_exist.json")
	if err == nil {
		t.Error("want error, got nil")
	}
	_, err = jsons.MergeAs(jsons.FormatJSON, []string{"file_not_exist.unknown"})
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestLoadBadBytes(t *testing.T) {
	a := []byte("{")
	_, err := jsons.Merge(a)
	if err == nil {
		t.Error("want error, got nil")
	}
	_, err = jsons.MergeAs(jsons.FormatJSON, a)
	if err == nil {
		t.Error("want error, got nil")
	}
}
func TestLoadBadBytesSlice(t *testing.T) {
	a := [][]byte{
		[]byte("{"),
	}
	_, err := jsons.Merge(a)
	if err == nil {
		t.Error("want error, got nil")
	}
	_, err = jsons.MergeAs(jsons.FormatJSON, a)
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestLoadBadReader(t *testing.T) {
	a := strings.NewReader(`{}`)
	b := strings.NewReader(`{`)
	_, err := jsons.Merge(a, b)
	if err == nil {
		t.Error("want error, got nil")
	}
	_, err = jsons.MergeAs(jsons.FormatJSON, a, b)
	if err == nil {
		t.Error("want error, got nil")
	}
}
func TestLoadBadReaders(t *testing.T) {
	a := []io.Reader{
		strings.NewReader(`{}`),
		strings.NewReader(`{`),
	}
	_, err := jsons.Merge(a)
	if err == nil {
		t.Error("want error, got nil")
	}
	_, err = jsons.MergeAs(jsons.FormatJSON, a)
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestLoadReaderError(t *testing.T) {
	a := &errReader{}
	_, err := jsons.Merge(a)
	if err == nil {
		t.Error("want error, got nil")
	}
	_, err = jsons.MergeAs(jsons.FormatJSON, a)
	if err == nil {
		t.Error("want error, got nil")
	}
}
func TestLoadReadersError(t *testing.T) {
	a := []io.Reader{
		&errReader{},
	}
	_, err := jsons.Merge(a)
	if err == nil {
		t.Error("want error, got nil")
	}
	_, err = jsons.MergeAs(jsons.FormatJSON, a)
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestMergeApplyRulesError(t *testing.T) {
	a := []byte(`
	  {
		"a": [
		  {
			"_tag": "a",
			"value": 1
		  }
		]
	  }
	`)
	b := []byte(`
	  {
		"a": [
		  {
			"_tag": "a",
			"value": false
		  }
		]
	  }
	`)
	_, err := jsons.Merge(a, b)
	if err == nil {
		t.Error("want error, got nil")
	}
	_, err = jsons.MergeAs(jsons.FormatJSON, a, b)
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestLoadUnsupportInput(t *testing.T) {
	_, err := jsons.Merge(true)
	if err == nil {
		t.Error("want error, got nil")
	}
	_, err = jsons.MergeAs(jsons.FormatJSON, true)
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestLoadNilInput(t *testing.T) {
	_, err := jsons.Merge(nil)
	if err != nil {
		t.Errorf("want nil, got: %s", err)
	}
	_, err = jsons.MergeAs(jsons.FormatJSON, nil)
	if err != nil {
		t.Errorf("want nil, got: %s", err)
	}
}

func TestGetExtensionsError(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterLoader("a", []string{".a1", ".a2"}, nil)
	_, err := m.Extensions("b")
	if err == nil {
		t.Error("want error, got nil")
	}
}
func TestRegisterLoaderError(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterLoader("a", []string{".a1", ".a2"}, nil)
	err := m.RegisterLoader("a", []string{".a1", ".a2"}, nil)
	if err == nil {
		t.Error("want error, got nil")
	}
	err = m.RegisterLoader("b", []string{".a1"}, nil)
	if err == nil {
		t.Error("want error, got nil")
	}
}

type errReader struct{}

func (r *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error")
}

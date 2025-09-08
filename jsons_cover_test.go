package jsons_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/qjebbs/go-jsons"
)

func TestMergeAs(t *testing.T) {
	a := []byte(`{"a":1}`)
	b := [][]byte{
		[]byte(`{"b":1}`),
	}
	c := strings.NewReader(`{"c":1}`)
	d := []io.Reader{
		strings.NewReader(`{"d":1}`),
	}
	want := []byte(`{"a":1,"b":1,"c":1,"d":1}`)
	m := jsons.NewMerger()
	got, err := m.MergeAs(jsons.FormatJSON, a, b, c, d)
	if err != nil {
		t.Fatal(err)
	}
	assertJSONEqual(t, want, got)
}

func TestMergeAsAuto(t *testing.T) {
	a := []byte(`{"a":1}`)
	b := [][]byte{
		[]byte(`{"b":1}`),
	}
	c := strings.NewReader(`{"c":1}`)
	d := []io.Reader{
		strings.NewReader(`{"d":1}`),
	}
	want := []byte(`{"a":1,"b":1,"c":1,"d":1}`)
	m := jsons.NewMerger()
	got, err := m.MergeAs(jsons.FormatAuto, a, b, c, d)
	if err != nil {
		t.Fatal(err)
	}
	assertJSONEqual(t, want, got)
}

func TestMergeAsUnknownFormat(t *testing.T) {
	m := jsons.NewMerger()
	_, err := m.MergeAs("unknown", []byte(`{}`))
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestMergeFileError(t *testing.T) {
	m := jsons.NewMerger()
	_, err := m.Merge("file_not_exist.json")
	if err == nil {
		t.Error("want error, got nil")
	}
	_, err = m.Merge([]string{"file_not_exist.unknown"})
	if err == nil {
		t.Error("want error, got nil")
	}
	_, err = m.Merge([]string{"file_not_exist.json"})
	if err == nil {
		t.Error("want error, got nil")
	}
	_, err = m.MergeAs(jsons.FormatJSON, "file_not_exist.json")
	if err == nil {
		t.Error("want error, got nil")
	}
	_, err = m.MergeAs(jsons.FormatJSON, []string{"file_not_exist.unknown"})
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
	m := jsons.NewMerger()
	_, err = m.MergeAs(jsons.FormatJSON, a)
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
	m := jsons.NewMerger()
	_, err = m.MergeAs(jsons.FormatJSON, a)
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
	m := jsons.NewMerger()
	_, err = m.MergeAs(jsons.FormatJSON, a, b)
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
	m := jsons.NewMerger()
	_, err = m.MergeAs(jsons.FormatJSON, a)
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
	m := jsons.NewMerger()
	_, err = m.MergeAs(jsons.FormatJSON, a)
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
	m := jsons.NewMerger()
	_, err = m.MergeAs(jsons.FormatJSON, a)
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestMergeApplyRulesError(t *testing.T) {
	a := []byte(`
	  {
		"a": [
		  {
			"tag": "a",
			"value": 1
		  }
		]
	  }
	`)
	b := []byte(`
	  {
		"a": [
		  {
			"tag": "a",
			"value": false
		  }
		]
	  }
	`)
	m := jsons.NewMerger(
		jsons.WithMergeBy("tag"),
	)
	_, err := m.Merge(a, b)
	if err == nil {
		t.Error("want error, got nil")
	}
	_, err = m.MergeAs(jsons.FormatJSON, a, b)
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestLoadUnsupportInput(t *testing.T) {
	_, err := jsons.Merge(true)
	if err == nil {
		t.Error("want error, got nil")
	}
	m := jsons.NewMerger()
	_, err = m.MergeAs(jsons.FormatJSON, true)
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestLoadNilInput(t *testing.T) {
	_, err := jsons.Merge(nil)
	if err != nil {
		t.Errorf("want nil, got: %s", err)
	}
	m := jsons.NewMerger()
	_, err = m.MergeAs(jsons.FormatJSON, nil)
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
	err := m.RegisterLoader("a", []string{".a1", ".a2"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = m.RegisterLoader("a", []string{".a1", ".a3"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = m.RegisterLoader("b", []string{".a2"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = m.RegisterLoader("b", []string{".a1"}, nil)
	if err == nil {
		t.Error("want error, got nil")
	}
	err = m.RegisterLoader(jsons.FormatAuto, []string{".a1"}, nil)
	if err == nil {
		t.Error("want error, got nil")
	}
}

type errReader struct{}

func (r *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error")
}

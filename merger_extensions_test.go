package jsons_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/qjebbs/go-jsons"
)

func TestGetExtensions(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterLoader("a", []string{".a1", ".a2"}, nil)
	want := []string{".a1", ".a2"}
	got, err := m.Extensions("a")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestGetAllExtensions(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterLoader("a", []string{".a1", ".a2"}, nil)
	m.RegisterLoader("b", []string{".b1", ".b2"}, nil)
	want := []string{".a1", ".a2", ".b1", ".b2"}
	got, err := m.Extensions(jsons.FormatAuto)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(got)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

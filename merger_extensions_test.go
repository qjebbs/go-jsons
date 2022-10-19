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
	got, err := m.GetExtensions("a")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestGetExtensionsError(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterLoader("a", []string{".a1", ".a2"}, nil)
	_, err := m.GetExtensions("b")
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestGetAllExtensions(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterLoader("a", []string{".a1", ".a2"}, nil)
	m.RegisterLoader("b", []string{".b1", ".b2"}, nil)
	want := []string{".a1", ".a2", ".b1", ".b2"}
	got, err := m.GetExtensions(jsons.FormatAuto)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(got)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

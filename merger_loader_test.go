package jsons_test

import (
	"testing"

	"github.com/qjebbs/go-jsons"
)

func TestRegisterLoaderError(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterLoader("a", []string{".a1", ".a2"}, nil)
	err := m.RegisterLoader("a", []string{".a1", ".a2"}, nil)
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestRegisterLoaderError2(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterLoader("a", []string{".a1", ".a2"}, nil)
	err := m.RegisterLoader("b", []string{".a1"}, nil)
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestLoadBadJSON(t *testing.T) {
	m := jsons.NewMerger()
	m.RegisterDefaultLoader()
	_, err := m.Merge([]byte("{"))
	if err == nil {
		t.Error("want error, got nil")
	}
}

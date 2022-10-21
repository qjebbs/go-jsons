package jsons

import (
	"errors"
	"testing"
)

func TestMust(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("want panic error, got nil")
		}
	}()
	must(nil, errors.New("test"))
}

func TestMust2(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Error("want ok, got panic")
		}
	}()
	must(nil, nil)
}

package preprocess_test

import (
	"os"
	"testing"

	"github.com/qjebbs/go-jsons/preprocess"
)

func TestExpandTypedEnv(t *testing.T) {
	// Setup test environment variables.
	os.Setenv("TT_NUM", "42")
	os.Setenv("TT_FLOAT", "3.14")
	os.Setenv("TT_NEG", "-7")
	os.Setenv("TT_BOOL_TRUE", "true")
	os.Setenv("TT_BOOL_FALSE", "false")
	os.Setenv("TT_BOOL_1", "1")
	os.Setenv("TT_BOOL_0", "0")
	os.Setenv("TT_STR", "hello")
	os.Setenv("TT_INVALID_NUM", "not_a_number")
	os.Setenv("TT_INVALID_BOOL", "not_a_boolean")
	defer func() {
		os.Unsetenv("TT_NUM")
		os.Unsetenv("TT_FLOAT")
		os.Unsetenv("TT_NEG")
		os.Unsetenv("TT_BOOL_TRUE")
		os.Unsetenv("TT_BOOL_FALSE")
		os.Unsetenv("TT_BOOL_1")
		os.Unsetenv("TT_BOOL_0")
		os.Unsetenv("TT_STR")
		os.Unsetenv("TT_INVALID_NUM")
		os.Unsetenv("TT_INVALID_BOOL")
	}()

	tests := []struct {
		input string
		want  any
	}{
		// String env var
		{"${TT_STR}", "hello"},
		// Number env vars
		{"${TT_NUM:number}", int64(42)},
		{"${TT_FLOAT:number}", 3.14},
		{"${TT_NEG:number}", int64(-7)},
		// Boolean env vars
		{"${TT_BOOL_TRUE:boolean}", true},
		{"${TT_BOOL_FALSE:boolean}", false},
		{"${TT_BOOL_1:boolean}", true},
		{"${TT_BOOL_0:boolean}", false},
		// Fallback to string on invalid value
		{"${TT_INVALID_NUM:number}", "not_a_number"},
		{"${TT_INVALID_BOOL:boolean}", "not_a_boolean"},
		// With extra path — type annotation stripped
		{"${TT_NUM:number}/some/path", "42/some/path"},
		{"${TT_BOOL_TRUE:boolean}/suffix", "true/suffix"},
		{"${TT_STR}/more", "hello/more"},
		// Plain text — returned as-is
		{"plain text", "plain text"},
		{"", ""},
		// Unset variable
		{"${TT_UNSET}", ""},
		{"${TT_UNSET:number}", int64(0)},
		{"${TT_UNSET:boolean}", false},
		// Unknown type name treated as string
		{"${TT_STR:unknown}", "hello"},
		// Multiple typed refs — all annotations stripped, expanded as string
		{"${TT_STR:string}/${TT_STR:string}/path", "hello/hello/path"},
		// No braces
		{"$TT_BOOL_1:boolean", "1:boolean"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := preprocess.ExpandTypedEnv("test", tt.input)
			if got != tt.want {
				t.Errorf("expandTypedEnv(%q) = %#v (%T), want %#v (%T)",
					tt.input, got, got, tt.want, tt.want)
			}
		})
	}
}

package preprocess

import (
	"os"
	"strconv"
	"strings"
	"unicode"
)

// ExpandTypedEnv is a preprocessor function that expands environment variables in string values using the TypedEnv syntax.
//
// TypedEnv syntax allows specifying the expected JSON type of an environment variable reference.
//
// Format: ${VAR:type}
// Currently supported types are: number, boolean, and string (default).
//
//	"${VAR}"              -> "1"
//	"${VAR:number}"       -> 1
//	"${VAR:boolean}"      -> true
//	"${VAR:number}/rest"  -> "1/reset"
//	"${VAR:boolean}/rest" -> "true/reset"
func ExpandTypedEnv(key string, value any) any {
	if str, ok := value.(string); ok {
		return expandTypedEnv(str)
	}
	return value
}

// expandTypedEnv expands a string using the TypedEnv syntax.
//
//   - If the whole string is a single typed env reference (e.g. "${VAR:number}"),
//     it resolves the environment variable and converts it to the target type.
//   - If the string contains extra content or multiple references, all are
//     expanded as strings and concatenated.
//   - If type conversion fails, the raw string value is returned (never fails).
func expandTypedEnv(raw string) any {
	exprs := parse(raw)

	// Single typed env reference (with type != String) — try type conversion.
	if len(exprs) == 1 {
		return exprs[0].buildTyped()
	}
	// Concatenate all segments as strings.
	var buf strings.Builder
	for _, e := range exprs {
		buf.WriteString(e.buildString())
	}
	return buf.String()
}

// expr represents a parsed segment of a TypedEnv expression.
type expr interface {
	buildTyped() any
	buildString() string
}

// plainExpr is literal text outside any ${...} reference.
type plainExpr string

func (v plainExpr) buildTyped() any     { return string(v) }
func (v plainExpr) buildString() string { return string(v) }

type stringExpr struct {
	varName string
}

func (v stringExpr) buildTyped() any     { return os.Getenv(v.varName) }
func (v stringExpr) buildString() string { return os.Getenv(v.varName) }

type numberExpr struct {
	varName string
}

func (v numberExpr) buildTyped() any {
	val := os.Getenv(v.varName)
	if val == "" {
		return int64(0)
	}
	if n, err := strconv.ParseInt(val, 10, 64); err == nil {
		return n
	}
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		return f
	}
	// fallback to string
	return val
}

func (v numberExpr) buildString() string { return os.Getenv(v.varName) }

type booleanExpr struct {
	varName string
}

func (v booleanExpr) buildTyped() any {
	val := os.Getenv(v.varName)
	if val == "" {
		return false
	}
	switch strings.ToLower(val) {
	case "true", "1", "yes", "on", "ok", "enabled":
		return true
	case "false", "0", "no", "off", "disabled":
		return false
	default:
		// fallback to string
		return val
	}
}
func (v booleanExpr) buildString() string {
	switch val := v.buildTyped().(type) {
	case bool:
		return strconv.FormatBool(val)
	case string:
		return val
	}
	return ""
}

func parse(raw string) []expr {
	runes := []rune(raw)
	var exprs []expr
	var plainBuf strings.Builder

	flushPlain := func() {
		if plainBuf.Len() > 0 {
			exprs = append(exprs, plainExpr(plainBuf.String()))
			plainBuf.Reset()
		}
	}

	for i := 0; i < len(runes); i++ {
		if runes[i] != '$' || i+1 >= len(runes) {
			plainBuf.WriteRune(runes[i])
			continue
		}

		// Bare $VAR syntax (no braces) — expand as string only
		if runes[i+1] != '{' {
			if !unicode.IsLetter(runes[i+1]) && runes[i+1] != '_' {
				plainBuf.WriteRune(runes[i])
				continue
			}
			j := i + 1
			for j < len(runes) && (unicode.IsLetter(runes[j]) || unicode.IsDigit(runes[j]) || runes[j] == '_') {
				j++
			}
			varName := string(runes[i+1 : j])
			flushPlain()
			exprs = append(exprs, stringExpr{varName})
			i = j - 1
			continue
		}

		// Potential ref starting at i
		j := i + 2 // skip ${
		varStart := j

		// Read var name until '}', ':', or end
		for j < len(runes) && runes[j] != '}' && runes[j] != ':' {
			j++
		}

		if j >= len(runes) {
			// No closing brace — treat rest as plain
			plainBuf.WriteString(string(runes[i:]))
			break
		}

		varName := string(runes[varStart:j])
		if varName == "" {
			// Empty var name — not a valid ref, emit as plain
			plainBuf.WriteRune(runes[i])
			continue
		}

		if runes[j] == ':' {
			// Typed ref — read type name
			j++ // skip ':'
			typeStart := j
			for j < len(runes) && unicode.IsLetter(runes[j]) {
				j++
			}

			if j >= len(runes) || runes[j] != '}' {
				// Malformed type — treat '$' as plain
				plainBuf.WriteRune(runes[i])
				continue
			}

			flushPlain()

			switch string(runes[typeStart:j]) {
			case "number":
				exprs = append(exprs, numberExpr{varName})
			case "boolean":
				exprs = append(exprs, booleanExpr{varName})
			default:
				exprs = append(exprs, stringExpr{varName})
			}
			i = j // loop will increment past '}'
		} else {
			// Plain env ref: runes[j] == '}'
			flushPlain()
			exprs = append(exprs, stringExpr{varName})
			i = j // loop will increment past '}'
		}
	}

	flushPlain()
	return exprs
}

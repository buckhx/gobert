package tokenization_test

import (
	"testing"

	"github.com/buckhx/gobert/tokenization"
)

func TestIsWhitespace(t *testing.T) {
	for _, test := range []struct {
		char  rune
		valid bool
	}{
		{' ', true},
		{'\t', true},
		{'\r', true},
		{'\n', true},
		{'\u00A0', true},
		{'A', false},
		{'-', false},
	} {
		if tokenization.IsWhitespace(test.char) != test.valid {
			t.Errorf("Invalid Whitespace Validation - %U, %t", test.char, test.valid)
		}
	}
}

func TestIsControl(t *testing.T) {
	for _, test := range []struct {
		char  rune
		valid bool
	}{
		{'\u0005', true},
		{'A', false},
		{' ', false},
		{'\t', false},
		{'\r', false},
		{'\U0001F4A9', false},
	} {
		if tokenization.IsControl(test.char) != test.valid {
			t.Errorf("Invalid Control Validation - %U, %t", test.char, test.valid)
		}
	}
}

func TestIsPunctuation(t *testing.T) {
	for _, test := range []struct {
		char  rune
		valid bool
	}{
		{'-', true},
		{'$', true},
		{'`', true},
		{'.', true},
		{'A', false},
		{' ', false},
	} {
		if tokenization.IsPunctuation(test.char) != test.valid {
			t.Errorf("Invalid Punctuation Validation - %U, %t", test.char, test.valid)
		}
	}
}

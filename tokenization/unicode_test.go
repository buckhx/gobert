package tokenization

import (
	"testing"
)

/***** Ported Tests *****/

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
		if isWhitespace(test.char) != test.valid {
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
		if isControl(test.char) != test.valid {
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
		if isPunctuation(test.char) != test.valid {
			t.Errorf("Invalid Punctuation Validation - %U, %t", test.char, test.valid)
		}
	}
}

/***** New Tests *****/

func TestIsChinese(t *testing.T) {
	for _, test := range []struct {
		char  rune
		valid bool
	}{
		{'\u535A', true},
		{'\u63A8', true},
		{'A', false},
		{' ', false},
	} {
		if isChinese(test.char) != test.valid {
			t.Errorf("Invalid Chinese (CJK) Validation - %U, %t", test.char, test.valid)
		}
	}
}

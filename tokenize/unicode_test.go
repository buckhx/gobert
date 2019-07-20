package tokenize

import (
	"testing"
)

/***** Ported Tests *****/

func Test_isWhitespace(t *testing.T) {
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

func Test_isControl(t *testing.T) {
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

func Test_isPunctuation(t *testing.T) {
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

func Test_isChinese(t *testing.T) {
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

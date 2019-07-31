package tokenize

import (
	"reflect"
	"testing"
)

func TestTestBasicLower(t *testing.T) {
	for i, test := range []struct {
		text   string
		lower  bool
		tokens []string
	}{
		{"HELLó THERE", true, []string{"hello", "there"}},
		{"HELLó THERE", false, []string{"HELLó", "THERE"}},
	} {
		toks := Basic{Lower: test.lower}.Tokenize(test.text)
		if !reflect.DeepEqual(toks, test.tokens) {
			t.Errorf("Test %d - Invalid Lowering %t - Want: %s, Got: %s", i, test.lower, test.tokens, toks)

		}
	}
}

func Test_tokenizeWhitespace(t *testing.T) {
	for i, test := range []struct {
		text   string
		tokens []string
	}{
		{"", nil},
		{"a", []string{"a"}},
		{" a", []string{"a"}},
		{"a ", []string{"a"}},
		{" a ", []string{"a"}},
		{"abc", []string{"abc"}},
		{"a  b", []string{"a", "b"}},
		{" abc ", []string{"abc"}},
		{"abc  def ", []string{"abc", "def"}},
		{"abc \u535A\u535A \u535A ", []string{"abc", "\u535A\u535A", "\u535A"}},
	} {
		toks := tokenizeWhitespace(test.text)
		if !reflect.DeepEqual(toks, test.tokens) {
			t.Errorf("Test %d - Invalid Whitespace Tokenization - Want: %s, Got: %s", i, test.tokens, toks)
		}
	}
}

func Test_padChinese(t *testing.T) {
	for i, test := range []struct {
		text   string
		padded string
	}{
		{"", ""},
		{"a", "a"},
		{" a", " a"},
		{"\u535A", " \u535A "},
		{"\u535A\u535A", " \u535A  \u535A "},
		{"abc \u535A\u535A \u535A ", "abc  \u535A  \u535A   \u535A  "},
		{"ah\u535A\u63A8zz", "ah \u535A  \u63A8 zz"},
	} {
		padded := padChinese(test.text)
		if padded != test.padded {
			t.Errorf("Test %d - Invalid Chinese Padding - Want: %q, Got: %q", i, test.padded, padded)
		}
	}
}

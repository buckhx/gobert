package tokenization_test

import (
	"testing"

	"github.com/buckhx/gobert/tokenization"
)

func TestBasic(t *testing.T) {
	for _, test := range []struct {
		name   string
		lower  bool
		text   string
		tokens []string
	}{
		{"chinese", false, "ah\u535A\u63A8zz", []string{"ah", "\u535A", "\u63A8", "zz"}},
		{"lower multi", true, " \tHeLLo!how  \n Are yoU?  ", []string{"hello", "!", "how", "are", "you", "?"}},
		{"lower single", true, "H\u00E9llo", []string{"hello"}},
		{"no lower multi", false, " \tHeLLo!how  \n Are yoU?  ", []string{"HeLLo", "!", "how", "Are", "yoU", "?"}},
		{"no lower single", false, "H\u00E9llo", []string{"H\u00E9llo"}},
	} {
		tkz := tokenization.Basic{Lower: test.lower}
		toks := tkz.Tokenize(test.text)
		if !equalTokens(toks, test.tokens) {
			t.Errorf("Test %s - Invalid Tokenization - Want: %v, Got: %v", test.name, test.tokens, toks)
		}
	}
}

func TestWordPiece(t *testing.T) {
	vocab := tokenization.NewVocab([]string{"[UNK]", "[CLS]", "[SEP]", "want", "##want", "##ed", "wa", "un", "runn", "##ing"})
	for i, test := range []struct {
		text   string
		tokens []string
	}{
		{"", []string{}},
		{"unwanted running", []string{"un", "##want", "##ed", "runn", "##ing"}},
		{"unwantedX running", []string{"[UNK]", "runn", "##ing"}},
	} {
		tkz := tokenization.WordPiece{Vocab: vocab}
		toks := tkz.Tokenize(test.text)
		if !equalTokens(toks, test.tokens) {
			t.Errorf("Test %d - Invalid Tokenization - Want: %v, Got: %v", i, test.tokens, toks)
		}
	}
}

// https://arxiv.org/pdf/1609.08144.pdf

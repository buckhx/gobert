package tokenization_test

import (
	"testing"

	"github.com/buckhx/gobert/tokenization"
)

func TestNewVocab(t *testing.T) {
	toks := []string{"abc", "def", "\u535A"}
	vocab := tokenization.NewVocab(toks)
	// TODO better testing and cover set semantics
	for i, tok := range toks {
		if vocab.Get(tok) != tokenization.ID(i) {
			t.Error("New Vocab Error")
		}
	}
}

func TestVocabLongestSubstring(t *testing.T) {
	toks := []string{"a", "aa", "aaa", "\u535A"}
	vocab := tokenization.NewVocab(toks)
	for i, test := range []struct {
		term string
		sub  string
	}{
		{"", ""},
		{"bbb", ""},
		{"aabb", "aa"},
		{"\u535Aaabb", "\u535A"},
	} {
		sub := vocab.LongestSubstring(test.term)
		if sub != test.sub {
			t.Errorf("Test %d - Invalid Longest Substring - Want: %v, Got: %v", i, test.sub, sub)
		}
	}
}

func TestVocabConvertTokens(t *testing.T) {
	vocab := tokenization.NewVocab([]string{"[UNK]", "[CLS]", "[SEP]", "want", "##want", "##ed", "wa", "un", "runn", "##ing"})
	for i, test := range []struct {
		tokens []string
		ids    []tokenization.ID
	}{
		{[]string{"un", "##want", "##ed", "runn", "##ing"}, []tokenization.ID{7, 4, 5, 8, 9}},
	} {
		ids := vocab.ConvertTokens(test.tokens)
		if !equalIDs(ids, test.ids) {
			t.Errorf("Test %d - Invalid Vocab IDs - Want: %v, Got: %v", i, test.ids, ids)
		}
	}
}

func equalTokens(x []string, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func equalIDs(x []tokenization.ID, y []tokenization.ID) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

package vocab_test

import (
	"reflect"
	"testing"

	"github.com/buckhx/gobert/vocab"
)

func TestNewVocab(t *testing.T) {
	toks := []string{"abc", "def", "\u535A"}
	voc := vocab.New(toks)
	// TODO better testing and cover set semantics
	for i, tok := range toks {
		if voc.Get(tok) != vocab.ID(i) {
			t.Error("New Vocab Error")
		}
	}
}

func TestVocabLongestSubstring(t *testing.T) {
	toks := []string{"a", "aa", "aaa", "\u535A"}
	voc := vocab.New(toks)
	for i, test := range []struct {
		term string
		sub  string
	}{
		{"", ""},
		{"bbb", ""},
		{"aabb", "aa"},
		{"\u535Aaabb", "\u535A"},
	} {
		sub := voc.LongestSubstring(test.term)
		if sub != test.sub {
			t.Errorf("Test %d - Invalid Longest Substring - Want: %v, Got: %v", i, test.sub, sub)
		}
	}
}

func TestVocabConvertTokens(t *testing.T) {
	voc := vocab.New([]string{"[UNK]", "[CLS]", "[SEP]", "want", "##want", "##ed", "wa", "un", "runn", "##ing"})
	for i, test := range []struct {
		tokens []string
		ids    []vocab.ID
	}{
		{[]string{"un", "##want", "##ed", "runn", "##ing"}, []vocab.ID{7, 4, 5, 8, 9}},
	} {
		ids := voc.ConvertTokens(test.tokens)
		if !reflect.DeepEqual(ids, test.ids) {
			t.Errorf("Test %d - Invalid Vocab IDs - Want: %v, Got: %v", i, test.ids, ids)
		}
	}
}

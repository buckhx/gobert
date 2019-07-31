package tokenize

import (
	"reflect"
	"testing"

	"github.com/buckhx/gobert/tokenize/vocab"
)

func TestFeatureCount(t *testing.T) {
	voc := vocab.New([]string{"[CLS]", "[SEP]", "the", "dog", "is", "hairy", "."})
	ff := FeatureFactory{Tokenizer: NewTokenizer(voc), SeqLen: 7}
	for _, test := range []struct {
		text  string
		count int
	}{
		{"", 2},
		{"the", 3},
		{"hello", 3},
		{"there we go", 6},
		{"mama mia, there we go again", 7},
	} {
		f := ff.Features(test.text)[0]
		if f.Count() != test.count {
			t.Errorf("Invalid Feature Count - Want: %d, Got %d", test.count, f.Count())
		}
	}
}

func Test_sequenceFeature(t *testing.T) {
	voc := vocab.New([]string{"[CLS]", "[SEP]", "the", "dog", "is", "hairy", "."})
	tkz := NewTokenizer(voc)
	for _, test := range []struct {
		text    string
		feature Feature
	}{
		// TODO more tests, but this one covers some good edge cases
		{"the dog is hairy. ||| the ||| a dog is hairy", Feature{
			ID:       0,
			Text:     "the dog is hairy. ||| the ||| a dog is hairy",
			Tokens:   []string{"[CLS]", "the", "dog", "[SEP]", "the", "[SEP]", "[UNK]", "[SEP]"},
			TokenIDs: []int32{0, 2, 3, 1, 2, 1, -1, 1},
			Mask:     []int32{1, 1, 1, 1, 1, 1, 1, 1},
			TypeIDs:  []int32{0, 0, 0, 0, 1, 1, 2, 2},
		}},
	} {
		f := sequenceFeature(tkz, 8, test.text)
		if !reflect.DeepEqual(f, test.feature) {
			t.Errorf("Invalid Sequence Feature - Want: %+v, Got: %+v", test.feature, f)
		}
	}
}

func Test_sequenceTruncate(t *testing.T) {
	for _, test := range []struct {
		seqs   [][]string
		len    int32
		tokens [][]string
	}{
		{nil, 1, nil},
		{[][]string{}, 1, [][]string{}},
		{[][]string{{"a1"}, {"b1"}, {"c1", "c2"}}, -1, [][]string{{}, {}, {}}},
		{[][]string{{"a1"}, {"b1"}, {"c1", "c2"}}, 0, [][]string{{}, {}, {}}},
		{[][]string{{"a1"}, {"b1"}, {"c1", "c2"}}, 1, [][]string{{"a1"}, {}, {}}},
		{[][]string{{"a1"}, {"b1"}, {"c1", "c2"}}, 3, [][]string{{"a1"}, {"b1"}, {"c1"}}},
		{[][]string{{"a1"}, {"b1"}, {"c1", "c2"}}, 10, [][]string{{"a1"}, {"b1"}, {"c1", "c2"}}},
	} {
		toks := truncate(test.seqs, test.len)
		if !reflect.DeepEqual(toks, test.tokens) {
			t.Errorf("Invalid Sequence Truncate - Want: %+v, Got: %+v", test.tokens, toks)
		}
	}
}

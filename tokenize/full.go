package tokenize

import "github.com/buckhx/gobert/tokenize/vocab"

// Full is a FullTokenizer which comprises of a Basic & Wordpiece tokenizer
type Full struct {
	Basic     Basic
	Wordpiece Wordpiece
}

// Tokenize will tokenize the input text
// First basic is applited, then wordpiece on the tokens froms basic
func (f Full) Tokenize(text string) []string {
	var toks []string
	for _, tok := range f.Basic.Tokenize(text) {
		toks = append(toks, f.Wordpiece.Tokenize(tok)...)
	}
	return toks
}

// Vocab returns the vocub used for this tokenizer
func (f Full) Vocab() vocab.Dict {
	return f.Wordpiece.vocab
}

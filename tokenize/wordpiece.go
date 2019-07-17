package tokenize

import (
	"fmt"

	"github.com/buckhx/gobert/vocab"
)

// DefaultMaxWordChars is the max length of a token for it to be tokenized, otherwise marked as unknown
const DefaultMaxWordChars = 200

// DefaultUnknownToken is the token used to signify an unkown token
const DefaultUnknownToken = "[UNK]"

type Wordpiece struct {
	vocab        vocab.Dict
	maxWordChars int
	unknownToken string
}

// NewWordpiece returns a WordpieceTokenizer with the default settings.
// Generally should be used in a FullTokenizer
func NewWordpiece(voc vocab.Dict) Wordpiece {
	return Wordpiece{
		vocab:        voc,
		maxWordChars: DefaultMaxWordChars,
		unknownToken: DefaultUnknownToken,
	}
}

func (wp Wordpiece) Tokenize(text string) []string {
	// TODO: determine if utf8 conversion is necessary, per python impl
	// text = convert_to_unicode(text)
	var toks []string
	for _, tok := range tokenizeWhitespace(text) {
		if len(tok) > wp.maxWordChars {
			toks = append(toks, wp.unknownToken)
			continue
		}
		for len(tok) > 0 && tok != "##" {
			sub := wp.vocab.LongestSubstring(tok)
			if sub == "" {
				toks = append(toks, wp.unknownToken)
				break
			}
			toks = append(toks, sub)
			tok = fmt.Sprintf("##%s", tok[len(sub):])
		}
	}
	return toks
}

// SetUnkownToken will set the max chars for a word to be tokenized,
// generally this should be congfigured through the FullTokenizer
func (wp Wordpiece) SetMaxWordChars(c int) {
	wp.maxWordChars = c
}

// SetUnkownToken will set the , generally this should be congfigured through the FullTokenizer
func (wp Wordpiece) SetUnknownToken(tok string) {
	wp.unknownToken = tok
}

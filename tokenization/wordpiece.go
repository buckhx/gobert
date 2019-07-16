package tokenization

import "fmt"

const maxWordChars = 200
const unknownToken = "[UNK]"

type WordPiece struct {
	Vocab Vocab
}

func (wt WordPiece) Tokenize(text string) []string {
	// TODO: determine if utf8 conversion is necessary, per python impl
	// text = convert_to_unicode(text)
	var toks []string
	for _, tok := range tokenizeWhitespace(text) {
		if len(tok) > maxWordChars {
			toks = append(toks, unknownToken)
			continue
		}
		for len(tok) > 0 && tok != "##" {
			sub := wt.Vocab.LongestSubstring(tok)
			if sub == "" {
				toks = append(toks, unknownToken)
				break
			}
			toks = append(toks, sub)
			tok = fmt.Sprintf("##%s", tok[len(sub):])
		}
	}
	return toks
}

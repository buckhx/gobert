package tokenization

import "strings"

// Tokenizer is an interface for chunking a string into it's tokens as per the BERT implematation
type Tokenizer interface {
	Tokenize(text string) (tokens []string, err error)
}

//tokenizeWhitespace splits text into tokeens by whitespace, per python semantics empty strings are not included
func tokenizeWhitespace(text string) []string {
	split := strings.Split(text, " ")
	var toks []string
	for _, tok := range split {
		if tok != "" {
			toks = append(toks, tok)
		}
	}
	return toks

}

//padChinese will add space padding around all CJK chars
// This implementation matches BasicTokenizer._tokenize_chinese_chars
func padChinese(text string) string {
	var b strings.Builder
	for _, c := range text {
		if isChinese(c) {
			b.WriteRune(' ')
			b.WriteRune(c)
			b.WriteRune(' ')
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}

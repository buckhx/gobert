package vocab

import (
	"bufio"
	"os"
)

// ID is used to identify vocab items
type ID int32

/*
func (id ID) String() string {
	return fmt.Sprint(id)
}
*/

// Dict is a container for tokens
// NOTE: python uses an OrderedDict, unsure of implications
type Dict struct {
	tokens map[string]ID
}

// FromFile will read a newline delimited file into a Dict
func FromFile(path string) (Dict, error) {
	// TODO test
	f, err := os.Open(path)
	if err != nil {
		// TODO wrap w/ stdlib
		return Dict{}, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	voc := Dict{}
	for scanner.Scan() {
		voc.Add(scanner.Text())
	}
	return voc, nil
}

func New(tokens []string) Dict {
	v := make(map[string]ID, len(tokens))
	for i, t := range tokens {
		v[t] = ID(i)
	}
	return Dict{tokens: v}
}

// Add will add an item to the vocabulary, is not thread-safe
func (v Dict) Add(token string) {
	v.tokens[token] = ID(v.Size())
}

// Get will return the ID of the token in the vocab. Will be negative if it doesn't exists
func (v Dict) Get(token string) ID {
	id, ok := v.tokens[token]
	if !ok {
		return ID(-1)
	}
	return ID(id)
}

/*
// GetToken will get a token by the ID, returns the mepty string if ID does not exist
func (v Dict) GetToken(id ID) token {
	for k, v := range v.tokens {
		if v =

	}
}

// HasID returns true if the vocab contains the token
func (v Dict) HasID(id ID) bool {
	for k, v := range v.tokens {
		if v =
	}
}

// HasToken returns true if the
func (v Dict) HasToken(token string) bool {

}
*/

// SIze returns the size of the vocabulary
func (v Dict) Size() int {
	return len(v.tokens)
}

// LongestSubstring returns the longest token that is a substring of the token
func (v Dict) LongestSubstring(token string) string {
	// Greedt, optimize to trie if needed
	for i := len(token); i > 0; i-- {
		sub := token[:i]
		if _, ok := v.tokens[sub]; ok {
			return sub
		}
	}
	return ""
}

func (v Dict) ConvertItems(items []string) []ID {
	ids := make([]ID, len(items))
	for i, m := range items {
		ids[i] = v.tokens[m]
	}
	return ids
}

func (v Dict) ConvertTokens(tokens []string) []ID {
	return v.ConvertItems(tokens)
}

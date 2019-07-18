package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config is a BERT MOdel Configuration
type Config struct {
	AttentionProbsDropoutProb float64 `json:"attention_probs_dropout_prob"`
	HiddenAct                 string  `json:"hidden_act"`
	HiddenDropoutProb         float64 `json:"hidden_dropout_prob"`
	HiddenSize                int     `json:"hidden_size"`
	InitializerRange          float64 `json:"initializer_range"`
	IntermediateSize          int     `json:"intermediate_size"`
	MaxPositionEmbeddings     int     `json:"max_position_embeddings"`
	NumAttentionHeads         int     `json:"num_attention_heads"`
	NumHiddenLayers           int     `json:"num_hidden_layers"`
	TypeVocabSize             int     `json:"type_vocab_size"`
	VocabSize                 int     `json:"vocab_size"`
}

// LoadConfig reads a config from filepath
func LoadConfig(path string) (Config, error) {
	c := Config{}
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(d, &c)
	if err != nil {
		err = fmt.Errorf("Config Unmarshal Error: %s", err)
	}
	return c, err
}

package tokenizer

import (
	"fmt"

	"g1.wpp2.hsnr/inr/boolret/index"
)

type Tokenizer struct {
	index *index.Index
}

func InitTokenizer(ind *index.Index) *Tokenizer {
	return &Tokenizer{index: ind}
}

func (t *Tokenizer) ParseFile(path string) error {
	//TODO
	return fmt.Errorf("Not implemented")
}

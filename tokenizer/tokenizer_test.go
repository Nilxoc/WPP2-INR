package tokenizer

import (
	"testing"

	"g1.wpp2.hsnr/inr/boolret/index"
)

const text1 string = "FOO\tBar\t\"Title\"\tHallo Welt, dies ist ein Test!"
const text1TokenCount int = 7

var text1Tokens = []string{"title", "hallo", "welt", "dies", "ist", "ein", "test"}

func TestTokenizer(t *testing.T) {

	//TODO
	idx := index.NewIndexEmpty(1, 1, 0.5)
	tk := InitTokenizer(idx)

	err := tk.ParseString(text1)
	if err != nil {
		t.Errorf("Parse String failed! %v", err)
	}

	if idx.Len() != text1TokenCount {
		t.Errorf("Expected Token count of %d got %d instead", text1TokenCount, idx.Len())
	}

	if idx.GetTerm("title") == nil || idx.GetTerm("title").Docs[0].Pos[0] != 1 {
		t.Errorf("Expected Term %s to exist and to be at position %d", "title", 1)
	}

	for _, x := range text1Tokens {
		if idx.GetTerm(x) == nil {
			t.Errorf("Expected Term %s to exist", x)
		}
	}

}

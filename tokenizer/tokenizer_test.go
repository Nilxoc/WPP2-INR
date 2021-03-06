package tokenizer

import (
	"testing"

	"g1.wpp2.hsnr/inr/boolret/config"

	"g1.wpp2.hsnr/inr/boolret/index"
)

const text1 string = "FOO\tBar\tnewTitle\t\"Title\" Hallo Welt, dies i'm % ist ein Test!"
const text1TokenCount int = 8

var text1Tokens = map[string]int{"title": 1, "hallo": 2, "welt": 3, "dies": 4, "i'm": 5, "ist": 6, "ein": 7, "test": 8}

func TestTokenizer(t *testing.T) {
	//TODO
	cfg := config.Config{
		PDoc:         "",
		PDict:        "",
		KGram:        1,
		JThresh:      0.5,
		CSpell:       false,
		CSpellThresh: 1,
	}

	idx := index.NewIndexEmpty(&cfg)
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

	for k := range text1Tokens {
		if idx.GetTerm(k) == nil {
			t.Errorf("Expected Term %s to exist", k)
		}
	}

	if idx.GetTerm("@") != nil {
		t.Errorf("Expected Term %s to not exist", "@")
	}

	for k, pos := range text1Tokens {
		idxPosT := idx.GetTerm(k)
		if idxPosT == nil {
			t.Errorf("Could not find %s!", k)
			continue
		}
		idxPos := idxPosT.Docs[0].Pos[0]
		if int(idxPos) != pos {
			t.Errorf("Invalid Position for token %s. Expected %d got %d", k, pos, idxPos)
		}
	}

}

package tokenizer

import (
	"fmt"
	"strings"

	"g1.wpp2.hsnr/inr/boolret/file"
	"g1.wpp2.hsnr/inr/boolret/index"
)

type Tokenizer struct {
	index *index.Index
}

func InitTokenizer(ind *index.Index) *Tokenizer {
	return &Tokenizer{index: ind}
}

func (t *Tokenizer) ParseFile(path string) error {
	fileString, err := file.ReadAsString(path)
	if err != nil {
		return err
	}
	return t.ParseString(fileString)

}

func (t *Tokenizer) ParseString(fileString string) error {

	//Creating a Array of lines for each doc
	docs := strings.Split(fileString, "\n")
	for docCounter, doc := range docs {
		//Loop over lines (documents)
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}
		//lines: 1: Doc-Name, 2: Doc
		parts := strings.Split(doc, "\t")
		if len(parts) != 4 {
			return fmt.Errorf("invalid Line")
		}
		text := strings.TrimSpace(parts[2]) + " " + strings.TrimSpace(parts[3])

		if err := t.evaluateText(text, docCounter+1); err != nil {
			return err
		}
	}
	return nil
}

func (t *Tokenizer) evaluateText(text string, docID int) error {
	tokensRaw := strings.Split(text, " ")
	tokenSub := 0

	hm := make(map[string]*index.Posting)

	for i, token := range tokensRaw {
		token = strings.ToLower(strings.Trim(token, "()\".;=-:,|][{}%/'!?&$§<>-_#+@*"))

		if token == "" {
			tokenSub += 1
			continue
		}

		if post, f := hm[token]; f {
			post.Pos = append(post.Pos, int64((i+1)-tokenSub))
		} else {
			tarr := make([]int64, 1)
			tarr[0] = int64((i + 1) - tokenSub)
			hm[token] = &index.Posting{
				DocID: int64(docID),
				Pos:   tarr,
			}
		}

	}

	for token, posting := range hm {
		t.index.AddTerm(token, posting)
	}

	return nil
}

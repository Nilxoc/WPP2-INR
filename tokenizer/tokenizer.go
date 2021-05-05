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

	//Creating a Array of lines for each doc
	docs := strings.Split(fileString, "\n")
	for _, doc := range docs {
		//Loop over lines (documents)
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}

		parts := strings.Split(doc, "\t")
		if len(parts) != 3 {
			return fmt.Errorf("Invalid Line")
		}
		text := strings.TrimSpace(parts[2])

		if err := t.evaluateText(text); err != nil {
			return err
		}
	}
	return nil
}

func (t *Tokenizer) evaluateText(text string) error {
	fmt.Println(text)
	return nil
}

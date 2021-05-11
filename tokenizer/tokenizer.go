package tokenizer

import (
	"fmt"
	"path"
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

func (t *Tokenizer) ParseSingleFile(path string) error {
	fileString, err := file.ReadAsString(path)
	if err != nil {
		return err
	}
	return t.ParseString(fileString)
}

/* UNUSED
func getIDFromFilename(name string) int {
	parts := strings.Split(name, "-")
	if len(parts) == 2 {
		id, err := strconv.ParseInt(parts[1], 10, 32)
		if err != nil {
			return -1
		}
		return int(id)
	}
	return -1
}
*/

func (t *Tokenizer) ParseMultiFile(pathD string) error {
	docs, err := file.ListFiles(pathD)
	if err != nil {
		return err
	}
	for dID, docPath := range docs {
		content, err := file.ReadAsString(docPath)
		if err != nil {
			return err
		}
		baseFileName := strings.TrimSuffix(path.Base(docPath), path.Ext(docPath))
		//docID := getIDFromFilename(baseFileName)
		/*if docID == -1 {
			return fmt.Errorf("Invalid Document ID found: %s", baseFileName)
		}*/
		if err = t.evaluateText(content, dID+1, baseFileName); err != nil {
			return err
		}
	}
	return nil
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
		//strings.TrimSpace(parts[2]) + " " +
		text := strings.TrimSpace(parts[3])
		idName := strings.TrimSpace(parts[0])

		if err := t.evaluateText(text, docCounter+1, idName); err != nil {
			return err
		}
	}
	return nil
}

func (t *Tokenizer) evaluateText(text string, docID int, docName string) error {
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
		t.index.AddTerm(token, posting, docName)
	}

	return nil
}

package index

import (
	"fmt"
	"strings"

	word2Vec "code.sajari.com/word2vec"
	"g1.wpp2.hsnr/inr/boolret/file"
)

type w2VIndex struct {
	//This Index Maps DOCUMENTS,not Terms, to their respective Vectors
	//A Document Vector is the normalized Vector of all its terms
	Index map[string]word2Vec.Vector
	Model *word2Vec.Model
}

type resultEntry struct {
	doc    string
	cosSim float32
}

func calcTextVec(text string, model *word2Vec.Model) word2Vec.Vector {
	parts := strings.Split(text, " ")
	wordMap := model.Map(parts)

	var docVec word2Vec.Vector
	wordCount := float32(0)

	for _, vec := range wordMap {
		if docVec == nil {
			docVec = vec
		} else {
			docVec.Add(1, vec)
		}
		wordCount += 1
	}

	for i, elem := range docVec {
		docVec[i] = elem / wordCount
	}

	return docVec
}

func buildIndex() (*w2VIndex, error) {
	filename := "/home/colin/Documents/WPP2/WPP2-INR/docs.txt"
	modelPath := "/home/colin/Documents/WPP2/WPP2-INR/binGlove.bin"

	modelReader, err := file.GetFileReader(modelPath)
	if err != nil {
		return nil, fmt.Errorf("error reading model: %v", err)
	}

	model, err := word2Vec.FromReader(modelReader)
	index := w2VIndex{Index: make(map[string]word2Vec.Vector), Model: model}

	if err != nil {
		return nil, fmt.Errorf("error building model: %v", err)
	}

	content, err := file.ReadAsString(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading docs file: %v", err)
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		sections := strings.Split(line, "\t")
		if len(sections) < 2 {
			continue
		}
		vec := calcTextVec(sections[1], model)
		index.Index[sections[0]] = vec
	}

	return &index, nil
}

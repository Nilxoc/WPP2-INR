package index

import (
	"fmt"
	"sort"
	"strings"

	word2Vec "code.sajari.com/word2vec"
	"g1.wpp2.hsnr/inr/boolret/file"
)

type w2VIndex struct {
	//This Index Maps DOCUMENTS,not Terms, to their respective Vectors
	//A Document Vector is the average Vector of all its terms
	Index map[string]word2Vec.Vector
	Texts map[string]string
	Model *word2Vec.Model
}

type resultEntry struct {
	Doc    string
	cosSim float32
}

//Calculates the average vector of a given Bag of Words(a query or a document)
func calcTextVec(text string, model *word2Vec.Model) word2Vec.Vector {
	parts := strings.Split(text, " ")
	wordMap := model.Map(parts)

	wordWeight := make(map[string]float64)

	for _, word := range parts {
		if val, ok := wordWeight[word]; ok {
			wordWeight[word] = val + 1
		} else {
			wordWeight[word] = 1
		}
	}

	var docVec []float64

	wordCount := float64(0)

	for word, vec := range wordMap {
		if docVec == nil {
			docVec = make([]float64, len(vec))
		}
		for i := range vec {
			docVec[i] += wordWeight[word] * float64(vec[i])
		}
		wordCount += wordWeight[word]
	}

	for i, elem := range docVec {
		docVec[i] = elem / wordCount
	}

	var resVec word2Vec.Vector = make(word2Vec.Vector, len(docVec))
	for i := range resVec {
		resVec[i] = float32(docVec[i])
	}

	return resVec
}

//Constructs an Index from the given file using the given model file
func BuildIndex(filename string, modelPath string) (*w2VIndex, error) {
	modelReader, err := file.GetFileReader(modelPath)
	if err != nil {
		return nil, fmt.Errorf("error reading model: %v", err)
	}

	model, err := word2Vec.FromReader(modelReader)
	index := w2VIndex{Index: make(map[string]word2Vec.Vector), Model: model, Texts: make(map[string]string)}

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
		vec := calcTextVec(sections[1], index.Model)
		index.Index[sections[0]] = vec
		index.Texts[sections[0]] = sections[1]
	}

	return &index, nil
}

//Evaluates a query on the index and returns a ranked list of results
func (index *w2VIndex) EvaluateQuery(query string) []resultEntry {
	queryVec := calcTextVec(query, index.Model)

	results := make([]resultEntry, len(index.Index))

	for doc, docVec := range index.Index {
		results = append(results, resultEntry{Doc: doc, cosSim: queryVec.Dot(docVec) / (queryVec.Norm() * docVec.Norm())})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].cosSim > results[j].cosSim
	})

	return results
}

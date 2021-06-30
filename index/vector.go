package index

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type BagTerm struct {
	//	term     string
	DocCount uint
	Count    uint
	Docs     DocumentRefs
}

type TermBag map[string]BagTerm

func makeTermBag(query string, idx *Index) (TermBag, map[int64]float64) {
	termList := strings.Split(query, " ")
	bag := make(TermBag, 0)
	documentsWeightMap := make(map[int64]float64)
	for _, t := range termList { //Range over all Terms
		if v, f := bag[t]; f {
			v.Count++ //Term found multiple times
		} else {
			termEntry := idx.GetTerm(t) //Term found first time. get all available info
			if termEntry == nil {
				fmt.Printf("Term %s not found in Index!\n", t)
				continue
			} else {
				bag[t] = BagTerm{DocCount: uint(termEntry.DocCount), Count: 1, Docs: termEntry.Docs}
				for _, d := range termEntry.Docs {
					documentsWeightMap[int64(d.Document.DocID)] = 0
				}
			}
		}
	}
	return bag, documentsWeightMap
}

func (idx *Index) Weighting(query TermBag, doc DocumentRef, k float64) float64 {
	var sum float64 = 0
	for _, term := range query {
		sum += float64(term.Count) * ((float64(doc.TermCount)) / (float64(doc.TermCount) + k*((float64(doc.Document.TotalLength))/(float64(idx.AvgDocLength))))) * math.Log10((float64(idx.DocCount))/(float64(term.DocCount)))
	}
	return sum
}

func (idx *Index) FastCosine(query string, n int) ([]string, error) {
	bag, scores := makeTermBag(query, idx)

	for _, term := range bag {
		for _, doc := range term.Docs {
			scores[int64(doc.Document.DocID)] += idx.Weighting(bag, doc, float64(idx.K)) // Increase Weighting with formular
		}
	}

	type DocumentListEntry struct {
		ID    int64
		Score float64
	}
	docList := make([]DocumentListEntry, 0)

	for k, _ := range scores {
		docList = append(docList, DocumentListEntry{ID: k, Score: scores[k] / float64(idx.Documents[k].TotalLength)})
	}
	sort.Slice(docList, func(i, j int) bool {
		return docList[i].Score > docList[j].Score
	})

	if len(docList) < n {
		n = len(docList)
	}

	fmt.Printf("Found total of %d\n", len(docList))
	res := make([]string, n)
	for i, d := range docList {
		if i >= n {
			break
		}
		res = append(res, idx.GetDocDisplay(d.ID))
	}
	return res, nil

	//TODO add length for docs
	//TODO have document lsit with length and terms in Index?
	// Divide by document length
	//Sort (heap?)
}

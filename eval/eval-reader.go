package eval

import (
	"fmt"
	"g1.wpp2.hsnr/inr/boolret/file"
	"strconv"
	"strings"
)

// RelevanceMap Maps Query IDs to slice of relevant documents
type RelevanceMap map[int64][]int64

func parseId(input string) int64 {
	parts := strings.Split(input, "-")
	res, _ := strconv.Atoi(parts[1])

	return int64(res)
}

func readRelevance(path string) (*RelevanceMap, error) {
	content, err := file.ReadAsString(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	relMap := RelevanceMap{}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		sections := strings.Split(line, "\t")
		if len(sections) < 4 {
			continue
		}
		qid := parseId(sections[0])
		docId := parseId(sections[2])

		if _, ok := relMap[qid]; !ok {
			relMap[qid] = make([]int64, 0, 1)
		}

		relMap[qid] = append(relMap[qid], docId)
	}

	return &relMap, nil
}

type ConfM struct {
	tp int32
	fp int32
	fn int32
}

func CalculateConfusion(found []int64, relevant []int64) ConfM {
	tp := int32(0)
	fp := int32(0)
	fn := int32(0)

	for _, docID := range found {
		contains := false

		for j, relID := range relevant {
			if docID == relID {
				contains = true
				relevant = append(relevant[:j], relevant[j+1:]...)
				break
			}
		}

		if contains {
			tp += 1
		} else {
			fp += 1
		}
	}

	fn = int32(len(relevant))

	return ConfM{tp, fp, fn}
}

func (m *ConfM) Precision() float64 {
	return float64(m.tp) / float64(m.tp+m.fp)
}

func (m *ConfM) Recall() float64 {
	return float64(m.tp) / float64(m.tp+m.fn)
}

func (m *ConfM) FMeasure(a float64) float64 {
	b2 := (1 - a) / a
	return ((b2 + 1) * m.Precision() * m.Recall()) / (b2*m.RPrecision() + m.Recall())
}

func (m *ConfM) F1Measure() float64 {
	return m.FMeasure(0.5)
}

func (m *ConfM) RPrecision() float64 {
	return float64(m.tp) / float64(m.tp+m.fn)
}

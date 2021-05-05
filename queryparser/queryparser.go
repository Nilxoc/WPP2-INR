package queryparser

import (
	"fmt"

	"g1.wpp2.hsnr/inr/boolret/index"
)

type QueryParser struct {
	index *index.Index
}

func InitParser(ind *index.Index) *QueryParser {
	return &QueryParser{index: ind}
}

func (qp *QueryParser) Evaluate(query string) (index.PostingList, error) {
	return nil, fmt.Errorf("Not implemented")
}

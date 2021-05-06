package query

import (
	"g1.wpp2.hsnr/inr/boolret/index"
)

type Query interface {
	Evaluate() (index.PostingList, error)
	Print()
}

type Parser interface {
	parse(query string) Query
}

func InitParser(ind *index.Index) *AstQueryParser {
	return &AstQueryParser{index: ind}
}

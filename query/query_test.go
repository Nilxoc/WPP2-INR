package query

import (
	"testing"

	"g1.wpp2.hsnr/inr/boolret/config"
	"g1.wpp2.hsnr/inr/boolret/index"
	"github.com/stretchr/testify/assert"
)

func TestWhenIdxEmptyNoResult(t *testing.T) {
	res := expectResult(t, "q")
	assert.True(t, res.Entry.Empty())
}

func TestWhenQueryEqDocResultFound(t *testing.T) {
	p := createParser()
	idx := p.Context.Index
	pl := []int64{1}
	idx.AddTerm("term", &index.Posting{DocID: 1, Pos: pl})

	res := expectResultP(t, "term", p)
	assert.Equal(t, 1, len(*res.Entry))
}

// FIXME: infinite loop in PostingList#Union?
func TestSubexpression(t *testing.T) {
	p := createParser()
	idx := p.Context.Index
	idx.AddTerm("java", &index.Posting{DocID: 1, Pos: []int64{1}})
	idx.AddTerm("is", &index.Posting{DocID: 1, Pos: []int64{2}})
	idx.AddTerm("what", &index.Posting{DocID: 2, Pos: []int64{1}})
	idx.AddTerm("java", &index.Posting{DocID: 2, Pos: []int64{2}})
	idx.AddTerm("does", &index.Posting{DocID: 3, Pos: []int64{1}})

	res := expectResultP(t, "(what OR does) AND java", p)
	assert.Equal(t, 1, len(*res.Entry))
}

// utility

func expectResult(t *testing.T, q string) *Result {
	parser := createParser()
	return expectResultP(t, q, parser)
}

func expectResultP(t *testing.T, q string, p *AstQueryParser) *Result {
	parsed, err := p.Parse(q)
	if err != nil {
		t.Errorf("%s", err)
	}
	eval, err := parsed.Evaluate()
	if err != nil {
		t.Errorf("%s", err)
	}
	return eval
}

func createParser() *AstQueryParser {
	cfg, idx := emptyIdx()
	ctx := Context{
		Index:  idx,
		Config: cfg,
	}
	parser := AstQueryParser{Context: &ctx}
	return &parser
}

func emptyIdx() (*config.Config, *index.Index) {
	cfg := config.DefaultConfig()
	idx := index.NewIndexEmpty(cfg)
	return cfg, idx
}

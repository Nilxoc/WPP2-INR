package query

import (
	"g1.wpp2.hsnr/inr/boolret/config"
	"g1.wpp2.hsnr/inr/boolret/index"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWhenIdxEmptyNoResult(t *testing.T) {
	res := expectResult(t, "q")
	assert.True(t, res.Entry.Empty())
}

// FIXME: currently fails, wrong result count
func WhenQueryEqDocResultFound(t *testing.T) {
	p := createParser()
	idx := p.Context.Index
	pl := []int64{1}
	idx.AddTerm("term", &index.Posting{DocID: 1, Pos: pl})

	res := expectResult(t, "term")
	assert.Equal(t, 1, len(*res.Entry))
}

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

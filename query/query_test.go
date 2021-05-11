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
	idx.AddTerm("term", &index.Posting{DocID: 1, Pos: []int64{1}}, "1")

	res := expectResultP(t, "term", p)
	assert.Equal(t, 1, len(*res.Entry))
}

func TestPhraseQuery(t *testing.T) {
	p := createParser()
	idx := p.Context.Index
	idx.AddTerm("first", &index.Posting{DocID: 1, Pos: []int64{1}}, "1")
	idx.AddTerm("second", &index.Posting{DocID: 1, Pos: []int64{2}}, "1")
	idx.AddTerm("third", &index.Posting{DocID: 1, Pos: []int64{3}}, "1")

	res := expectResultP(t, `"first second third"`, p)

	e := *res.Entry
	assert.Equal(t, 1, len(e))
	assert.Equal(t, int64(1), e[0].DocID)
}

func TestAndQuery(t *testing.T) {
	p := createParser()
	idx := p.Context.Index
	idx.AddTerm("deleted", &index.Posting{DocID: 1, Pos: []int64{1}}, "1")
	idx.AddTerm("code", &index.Posting{DocID: 1, Pos: []int64{2}}, "1")
	idx.AddTerm("is", &index.Posting{DocID: 1, Pos: []int64{3}}, "1")
	idx.AddTerm("debugged", &index.Posting{DocID: 1, Pos: []int64{4}}, "1")
	idx.AddTerm("code", &index.Posting{DocID: 1, Pos: []int64{5}}, "1")

	idx.AddTerm("debugged", &index.Posting{DocID: 2, Pos: []int64{1}}, "1")
	idx.AddTerm("code", &index.Posting{DocID: 2, Pos: []int64{2}}, "1")

	idx.AddTerm("code", &index.Posting{DocID: 3, Pos: []int64{1}}, "1")

	res := expectResultP(t, `code AND debugged`, p)

	e := *res.Entry
	assert.Equal(t, 2, len(e))
	assert.Equal(t, int64(1), e[0].DocID)
	assert.Equal(t, int64(2), e[1].DocID)
}

func TestSubexpression(t *testing.T) {
	p := createParser()
	idx := p.Context.Index
	idx.AddTerm("java", &index.Posting{DocID: 1, Pos: []int64{1}}, "1")
	idx.AddTerm("is", &index.Posting{DocID: 1, Pos: []int64{2}}, "1")
	idx.AddTerm("what", &index.Posting{DocID: 2, Pos: []int64{1}}, "1")
	idx.AddTerm("java", &index.Posting{DocID: 2, Pos: []int64{2}}, "1")
	idx.AddTerm("does", &index.Posting{DocID: 3, Pos: []int64{1}}, "1")

	res := expectResultP(t, "(what OR does) AND java", p)
	assert.Equal(t, 1, len(*res.Entry))
}

func TestProximity(t *testing.T) {
	p := createParser()
	idx := p.Context.Index
	idx.AddTerm("term1", &index.Posting{DocID: 1, Pos: []int64{1}}, "1")
	idx.AddTerm("term2", &index.Posting{DocID: 1, Pos: []int64{2}}, "1")
	idx.AddTerm("term3", &index.Posting{DocID: 1, Pos: []int64{3}}, "1")

	res := expectResultP(t, "term1 /2 term3", p)
	first := (*res.Entry)[0]
	assert.Equal(t, int64(1), first.DocID)
	pos := first.Pos
	assert.Equal(t, 2, len(pos))
	assert.Equal(t, int64(1), pos[0])
	assert.Equal(t, int64(3), pos[1])
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

package query

import (
	"g1.wpp2.hsnr/inr/boolret/config"
	"g1.wpp2.hsnr/inr/boolret/index"
	"testing"
)

// validates if sample queries can be compiled

func TestParseSingle(t *testing.T) {
	parseNoError(t, "test1")
}

func TestParseAnd(t *testing.T) {
	parseNoError(t, "test1 AND test2")
}

func TestParseProximity(t *testing.T) {
	parseNoError(t, "term4 /2 term5")
}

func TestParsePhrase(t *testing.T) {
	parseNoError(t, `term1 term2 term3`)
}

func TestParsePhraseNested(t *testing.T) {
	parseNoError(t, `term1 AND "term2 term3"`)
}

func TestParseBrackets(t *testing.T) {
	parseNoError(t, "( term1 AND term2 ) AND term3")
}

func TestParseFull(t *testing.T) {
	parseNoError(t, `("term1 term2" OR term3) AND NOT term4 /3 term5`)
}

func parseNoError(t *testing.T, q string) {
	cfg := &config.Config{}
	idx := index.NewIndexEmpty(cfg)
	ctx := Context{
		Index:  idx,
		Config: cfg,
	}
	parser := AstQueryParser{Context: &ctx}

	parsed, err := parser.Parse(q)
	if err != nil {
		t.Errorf("%s", err)
	} else {
		parsed.Print()
	}
}

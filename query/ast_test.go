package query

import (
	"g1.wpp2.hsnr/inr/boolret/config"
	"g1.wpp2.hsnr/inr/boolret/index"
	"testing"
)

// validates if sample queries can be compiled

func TestParseSingle(t *testing.T) {
	testCompile(t, "test1")
}

func TestParseAnd(t *testing.T) {
	testCompile(t, "test1 AND test2")
}

func TestParseProximity(t *testing.T) {
	testCompile(t, "term4 /2 term5")
}

func TestParsePhrase(t *testing.T) {
	testCompile(t, `term1 term2 term3`)
}

func TestParsePhraseNested(t *testing.T) {
	testCompile(t, `term1 AND "term2 term3"`)
}

func TestParseBrackets(t *testing.T) {
	testCompile(t, "( term1 AND term2 ) AND term3")
}

func TestParseFull(t *testing.T) {
	testCompile(t, `("term1 term2" OR term3) AND NOT term4 /3 term5`)
}

// samples from the official task

func TestCompileSample1(t *testing.T) {
	testCompile(t, "blood AND pressure")
}

func TestCompileSample2(t *testing.T) {
	testCompile(t, `blood AND NOT pressure`)
}

func TestCompileSample3(t *testing.T) {
	testCompile(t, `(blood OR pressure) AND cardiovascular`)
}

func TestCompileSample4(t *testing.T) {
	testCompile(t, `"blood pressure"`)
}

func TestCompileSample5(t *testing.T) {
	testCompile(t, `diet /10 health`)
}

func TestCompileSample6(t *testing.T) {
	testCompile(t, `diet /10 health AND "red wine"`)
}

func TestInvalidProximityNotCompile(t *testing.T) {
	testNotCompile(t, `term1 10 term2`)
}

func TestConsProximityNotCompile(t *testing.T) {
	testNotCompile(t, `term1 /10 /10 term2`)
}

func TestWrongBracketsNotCompile(t *testing.T) {
	testNotCompile(t, `[term1 AND term2]`)
}

func TestBracketNotClosedNotCompile(t *testing.T) {
	testNotCompile(t, `(term1 AND term2`)
}

func TestProximityNotClosedNotCompile(t *testing.T) {
	testNotCompile(t, `"term1 term2`)
}

// testCompile validates if the target query compiles
func testCompile(t *testing.T, q string) {
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

func testNotCompile(t *testing.T, q string) {
	cfg := &config.Config{}
	idx := index.NewIndexEmpty(cfg)
	ctx := Context{
		Index:  idx,
		Config: cfg,
	}
	parser := AstQueryParser{Context: &ctx}

	_, err := parser.Parse(q)
	if err != nil {
		// success
	} else {
		t.Errorf("expected compilation error: %s", q)
	}
}

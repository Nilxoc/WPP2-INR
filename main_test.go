package main

import (
	"path"
	"runtime"
	"testing"
	"time"

	"g1.wpp2.hsnr/inr/boolret/config"
	"g1.wpp2.hsnr/inr/boolret/query"

	"g1.wpp2.hsnr/inr/boolret/file"
	"g1.wpp2.hsnr/inr/boolret/index"
	"g1.wpp2.hsnr/inr/boolret/tokenizer"
)

func constructIndexForTest(t *testing.T) *index.Index {
	//TODO Sollte wahrscheinlich error verwenden wenn doc_dump nicht
	//gefunden
	var docSource string

	_, testFnPath, _, _ := runtime.Caller(0)
	workDir := path.Dir(testFnPath)
	docSource, err := file.AbsPath(path.Join(workDir, "doc_dump.txt"))
	if err != nil {
		t.Errorf("Could not get Abs Path to doc_dump.txt: %v", err)
		return nil
	}

	cfg := config.Config{
		PDoc:         "",
		PDict:        "",
		KGram:        2,
		JThresh:      0.2,
		CSpell:       true,
		CSpellThresh: 5,
	}
	var indexInstance *index.Index
	indexInstance = index.NewIndexEmpty(&cfg)

	tokenizer := tokenizer.InitTokenizer(indexInstance)

	tokStart := time.Now()

	if err := tokenizer.ParseSingleFile(docSource); err != nil {
		t.Logf("No doc_dump.txt file found. Skipping Performance-Testing")
		return nil
	}
	t.Logf("Loaded and parsed document within %s", time.Since(tokStart).String())

	return indexInstance
}

func TestMain(t *testing.T) {
	indexInstance := constructIndexForTest(t)
	if indexInstance == nil {
		return
	}
	token1 := "fresher"
	tStart := time.Now()
	res := indexInstance.GetTerm(token1)
	if res == nil {
		t.Error("Not found single token")
		return
	}
	t.Logf("Found single Term within %s", time.Since(tStart).String())

	token2 := "hrtzreher"

	rStart := time.Now()
	pl := indexInstance.GetTermListCorrected(token2)
	if pl == nil {
		t.Error("Not found corrected tokens")
		return
	}
	if len(pl) == 0 {
		t.Error("No results found. Try to lower j")
		return
	}
	t.Logf("Found corrected Terms within %s", time.Since(rStart).String())

}

func tryQuery(query string, t *testing.T, parser *query.AstQueryParser) {
	q, _ := parser.Parse(query)
	_, _ = q.Evaluate()
}

// FIXME: examples currently do not work
func TestExamples(t *testing.T) {
	indexInstance := constructIndexForTest(t)

	cfg, _ := config.Parse()

	ctx := query.Context{
		Index:  indexInstance,
		Config: cfg,
	}

	parser := query.AstQueryParser{Context: &ctx}

	//TODO performance
	tryQuery("blood AND pressure", t, &parser)
	tryQuery("blood AND NOT pressure", t, &parser)
	tryQuery("(blood OR pressure) AND cardiovascular", t, &parser)
	tryQuery("\"blood pressure\"", t, &parser)
	tryQuery("diet \\10 health", t, &parser)
	tryQuery("diet \\10 health AND \"red wine\"", t, &parser)
}

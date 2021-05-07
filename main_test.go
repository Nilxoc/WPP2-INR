package main

import (
	"path"
	"runtime"
	"testing"
	"time"

	"g1.wpp2.hsnr/inr/boolret/file"
	"g1.wpp2.hsnr/inr/boolret/index"
	"g1.wpp2.hsnr/inr/boolret/tokenizer"
)

func TestMain(t *testing.T) {

	var docSource string

	var k int = 2
	var r int = 5
	var j float64 = 0.2
	_, testFnPath, _, _ := runtime.Caller(0)
	workDir := path.Dir(testFnPath)
	docSource, err := file.AbsPath(path.Join(workDir, "doc_dump.txt"))
	if err != nil {
		t.Errorf("Could not get Abs Path to doc_dump.txt: %v", err)
		return
	}

	var indexInstance *index.Index
	indexInstance = index.NewIndexEmpty(k, r, float32(j))

	tokenizer := tokenizer.InitTokenizer(indexInstance)

	tokStart := time.Now()

	if err := tokenizer.ParseFile(docSource); err != nil {
		t.Logf("No doc_dump.txt file found. Skipping Performance-Testing")
		return
	}
	t.Logf("Loaded and parsed document within %s", time.Since(tokStart).String())

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
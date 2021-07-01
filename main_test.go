package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"g1.wpp2.hsnr/inr/boolret/eval"
	"g1.wpp2.hsnr/inr/boolret/file"

	"g1.wpp2.hsnr/inr/boolret/config"
	"g1.wpp2.hsnr/inr/boolret/query"

	"g1.wpp2.hsnr/inr/boolret/index"
	"g1.wpp2.hsnr/inr/boolret/tokenizer"
)

func constructIndexForTest(t *testing.T, b *testing.B) *index.Index {
	//TODO Sollte wahrscheinlich error verwenden wenn doc_dump nicht
	//gefunden
	var docSource string

	_, testFnPath, _, _ := runtime.Caller(0)
	workDir := filepath.Dir(testFnPath)
	docSource, err := filepath.Abs(filepath.Join(workDir, "doc_dump.txt"))
	if err != nil {
		if t != nil {
			t.Errorf("Could not get Abs Path to doc_dump.txt: %v", err)
		} else {
			b.Errorf("Could not get Abs Path to doc_dump.txt: %v", err)
		}
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
		if t != nil {
			t.Logf("No doc_dump.txt file found. Skipping Performance-Testing")
		} else {
			b.Logf("No doc_dump.txt file found. Skipping Performance-Testing")
		}
		return nil
	}
	if t != nil {
		t.Logf("Loaded and parsed document within %s", time.Since(tokStart).String())
	} else {
		b.Logf("Loaded and parsed document within %s", time.Since(tokStart).String())
	}

	return indexInstance
}

func TestMain(t *testing.T) {
	indexInstance := constructIndexForTest(t, nil)
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
	if len(pl.Docs) == 0 {
		t.Error("No results found. Try to lower j")
		return
	}
	t.Logf("Found corrected Terms within %s", time.Since(rStart).String())

}

func tryQuery(query string, t *testing.T, parser *query.AstQueryParser) {
	q, _ := parser.Parse(query)
	_, _ = q.Evaluate()
}

func TestExamples(t *testing.T) {
	indexInstance := constructIndexForTest(t, nil)

	if indexInstance == nil {
		return
	}

	cfg := config.DefaultConfig()

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
	tryQuery("diet /10 health", t, &parser)
	tryQuery("diet /10 health AND \"red wine\"", t, &parser)
}

func BenchmarkINR(b *testing.B) {
	idx := constructIndexForTest(nil, b)
	if idx == nil {
		return
	}
	cfg := config.DefaultConfig()
	ctx := query.Context{
		Index:  idx,
		Config: cfg,
	}
	parser := query.AstQueryParser{Context: &ctx}

	benchmarks := []string{
		"blood AND pressure",
		"blood AND NOT pressure",
		"(blood OR pressure) AND cardiovascular",
		"\"blood pressure\"",
		"diet /10 health",
		"diet /10 health AND \"red wine\"",
		"blod",
		"presure",
		"analysi",
	}

	for _, bm := range benchmarks {
		b.Run(bm, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				q, _ := parser.Parse(bm)
				_, _ = q.Evaluate()
			}
		})
	}
}

func TestConstruction(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.PDoc = workDirPath("doc_dump.txt")

	idx := index.NewIndexEmpty(cfg)
	tokenz := tokenizer.InitTokenizer(idx)

	absDocPath, err := file.AbsPath(cfg.PDoc)
	if err != nil {
		t.Errorf("cannot create path: %e", err)
	}

	// file is only present in local environments
	if _, err := os.Stat(absDocPath); os.IsNotExist(err) {
		t.Skipf("test file not found: %s", absDocPath)
	}

	err = tokenz.ParseSingleFile(absDocPath)
	if err != nil {
		t.Skipf("cannot parse file: %e", err)
	}
}

func BenchmarkConstruction(b *testing.B) {
	cfg := config.DefaultConfig()
	cfg.PDoc = workDirPath("doc_dump.txt")

	idx := index.NewIndexEmpty(cfg)
	tokenz := tokenizer.InitTokenizer(idx)

	absDocPath, err := file.AbsPath(cfg.PDoc)
	if err != nil {
		b.Errorf("cannot create path: %e", err)
	}

	b.Run("index-construction", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err = tokenz.ParseSingleFile(absDocPath)
			if err != nil {
				b.Errorf("cannot parse file: %e", err)
			}
		}
	})
}

func BenchmarkConstructionMultifile(b *testing.B) {
	cfg := config.DefaultConfig()
	cfg.PDoc = workDirPath(filepath.Join("prep", "out"))

	idx := index.NewIndexEmpty(cfg)
	tokenz := tokenizer.InitTokenizer(idx)

	absDocPath, err := file.AbsPath(cfg.PDoc)
	if err != nil {
		b.Errorf("cannot create path: %e", err)
	}

	b.Run("index-construction", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err = tokenz.ParseMultiFile(absDocPath)
			if err != nil {
				b.Errorf("cannot parse file: %e", err)
			}
		}
	})
}

func workDirPath(path string) string {
	_, testFnPath, _, _ := runtime.Caller(0)
	pwd := filepath.Dir(testFnPath)
	docSource, err := filepath.Abs(filepath.Join(pwd, path))
	if err != nil {
		panic(err)
	}
	return docSource
}

func TestW2VAccuracy(t *testing.T) {
	w2vInd, err := index.BuildIndex(workDirPath("docs.txt"), workDirPath("gnews.bin"))
	if err != nil {
		t.Errorf("Couldn't construct index, %v", err)
	}

	steps := [4]int32{5, 10, 20, 50}

	totMap := float64(0)

	//queries, err := loadQueries()
	allRels, _ := eval.ReadRelevance(workDirPath("qrel.txt"))
	queries := [5]Query{{QueryID: "PLAIN-121", Query: "berries to prevent muscle soreness", RelevantDocuments: (*allRels)[121]},
		{QueryID: "PLAIN-1021", Query: "diabetes", RelevantDocuments: (*allRels)[1021]},
		{QueryID: "PLAIN-15", Query: "why do heart doctors favor surgery and drugs over diet", RelevantDocuments: (*allRels)[15]},
		{QueryID: "PLAIN-145", Query: "fukushima radiation and seafood", RelevantDocuments: (*allRels)[145]},
		{QueryID: "PLAIN-1336", Query: "heart rate variability", RelevantDocuments: (*allRels)[1336]}}

	for _, query := range queries {
		recalls := make([]float64, 0, 4)
		precisions := make([]float64, 0, 4)
		f1s := make([]float64, 0, 4)

		results := w2vInd.EvaluateQuery(query.Query)

		i64Results := make([]int64, len(results))

		for i, res := range results {
			if res.Doc == "" {
				continue
			}
			i64Results[i] = parseId(res.Doc)
		}

		for _, step := range steps {
			confMat := eval.CalculateConfusion(i64Results[:step], query.RelevantDocuments)

			recalls = append(recalls, confMat.Recall())
			precisions = append(precisions, confMat.Precision())
			f1s = append(f1s, confMat.F1Measure())
		}

		confMat := eval.CalculateConfusion(i64Results, query.RelevantDocuments)

		rPrec := confMat.RPrecision()

		t.Logf("%.2f", rPrec)
		totMap += eval.MAPScore(i64Results, query.RelevantDocuments)
	}

	mapScore := totMap / float64(len(queries))

	t.Logf("%.2f", mapScore)
}

type Query struct {
	QueryID           string
	Query             string
	RelevantDocuments []int64
}

func parseId(input string) int64 {
	parts := strings.Split(input, "-")
	res, _ := strconv.Atoi(parts[1])

	return int64(res)
}

func loadQueries() ([]Query, error) {
	allRels, err := eval.ReadRelevance(workDirPath("qrel.txt"))
	if err != nil {
		return nil, err
	}
	file, err := file.ReadAsString(workDirPath("queries.txt"))
	if err != nil {
		return nil, err
	}
	lines := strings.Split(file, "\n")
	res := make([]Query, len(lines))
	for i, line := range lines {
		parts := strings.Split(strings.TrimSpace(line), "\t")
		if len(parts) < 2 {
			continue
		}
		res[i] = Query{QueryID: parts[0], Query: parts[1], RelevantDocuments: (*allRels)[parseId(parts[0])]}
	}
	return res, nil
}

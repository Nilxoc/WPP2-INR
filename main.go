package main

import (
	"flag"
	"fmt"

	"g1.wpp2.hsnr/inr/boolret/cli"
	"g1.wpp2.hsnr/inr/boolret/file"
	"g1.wpp2.hsnr/inr/boolret/index"
	"g1.wpp2.hsnr/inr/boolret/tokenizer"
)

func main() {

	var docSource string
	var dictSource string

	var k int = 0
	var r int = 0
	var j float64 = 0

	flag.StringVar(&docSource, "doc", "", "path to input file (txt)")
	flag.StringVar(&dictSource, "dict", "", "path to dict dump (skip index build) - if docSource also provided saves dict")
	flag.IntVar(&k, "k", 2, "k for kgram index")
	flag.IntVar(&r, "r", 5, "min count of terms to return")
	flag.Float64Var(&j, "j", 0.2, "jaccard threshold")
	flag.Parse()

	var indexInstance *index.Index

	if docSource != "" {
		//CREATE INDEX
		indexInstance = index.NewIndexEmpty(k, r, float32(j))

		tokenizer := tokenizer.InitTokenizer(indexInstance)

		absDocPath, err := file.AbsPath(docSource)
		if err != nil {
			panic(err)
		}

		if err := tokenizer.ParseFile(absDocPath); err != nil {
			panic(err)
		}
		fmt.Printf("File Parsed. Found %d terms\n", indexInstance.Len())

		if dictSource != "" {
			//SAVE INDEX
			absDictPath, err := file.AbsPath(dictSource)
			if err != nil {
				panic(err)
			}
			if err = indexInstance.SaveIndex(absDictPath); err != nil {
				panic(err)
			}
		}
	} else if dictSource != "" {
		//READ INDEX
		absDictPath, err := file.AbsPath(dictSource)
		if err != nil {
			panic(err)
		}
		indexInstance, err = index.NewIndexFromFile(absDictPath, k, r, float32(j))
		if err != nil {
			panic(err)
		}

		fmt.Printf("Dump loaded. Containing %d terms\n", indexInstance.Len())
	}

	if dictSource == "" && docSource == "" {
		panic(fmt.Errorf("No Input specified"))
	}

	//START CLI
	cl := cli.Init()

	//qp := queryparser.InitParser(indexInstance)

	cl.Print("Welcome to INR-System please insert your first Query..")
	for {

		//Temporary Token Finder (No complex Queries, just one single token)
		token := cl.GetInput()
		pl := indexInstance.GetTermListCorrected(token)
		if pl == nil {
			fmt.Println("Not found")
			continue
		}
		if len(pl) == 0 {
			fmt.Println("No results found. Try to lower j")
			continue
		}
		fmt.Println(pl.String())

		/* TODO
		q := cl.GetInput()
		pl, err := qp.Evaluate(q)
		if err != nil {
			cl.Print("Error executing Query: " + err.Error())
			continue
		}
		//TODO: PRINT pl
		cl.Print(pl.String())*/
	}

}

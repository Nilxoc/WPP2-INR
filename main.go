package main

import (
	"flag"
	"fmt"

	"g1.wpp2.hsnr/inr/boolret/cli"
	"g1.wpp2.hsnr/inr/boolret/file"
	"g1.wpp2.hsnr/inr/boolret/index"
	"g1.wpp2.hsnr/inr/boolret/queryparser"
	"g1.wpp2.hsnr/inr/boolret/tokenizer"
)

func main() {

	var docSource string
	var dictSource string

	flag.StringVar(&docSource, "doc", "", "path to input file (txt)")
	flag.StringVar(&dictSource, "dict", "", "path to dict dump (skip index build) - if docSource also provided saves dict")

	var indexInstance *index.Index

	if docSource != "" {
		//CREATE INDEX
		indexInstance = index.NewIndexEmpty()

		tokenizer := tokenizer.InitTokenizer(indexInstance)

		absDocPath, err := file.AbsPath(docSource)
		if err != nil {
			panic(err)
		}

		if err := tokenizer.ParseFile(absDocPath); err != nil {
			panic(err)
		}

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
	}
	if dictSource != "" {
		//READ INDEX
		absDictPath, err := file.AbsPath(dictSource)
		if err != nil {
			panic(err)
		}
		indexInstance, err = index.NewIndexFromFile(absDictPath)
		if err != nil {
			panic(err)
		}
	}

	if dictSource == "" && docSource == "" {
		panic(fmt.Errorf("No Input specified"))
	}

	//START CLI
	cl := cli.Init()

	qp := queryparser.InitParser(indexInstance)

	cl.Print("Welcome to INR-System please insert your first Query..")
	for {
		q := cl.GetInput()
		pl, err := qp.Evaluate(q)
		if err != nil {
			cl.Print("Error executing Query: " + err.Error())
			continue
		}
		//TODO: PRINT pl
		cl.Print(pl.String())
	}

}

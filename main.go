package main

import (
	"fmt"
	"log"

	"g1.wpp2.hsnr/inr/boolret/config"

	"g1.wpp2.hsnr/inr/boolret/cli"
	"g1.wpp2.hsnr/inr/boolret/file"
	"g1.wpp2.hsnr/inr/boolret/index"
	"g1.wpp2.hsnr/inr/boolret/tokenizer"
)

func main() {

	cfg, err := config.Parse()
	if err != nil {
		log.Panic(err)
	}

	var indexInstance *index.Index

	docPath := cfg.PDoc
	dictPath := cfg.PDict

	if docPath != "" {
		//CREATE INDEX
		indexInstance = index.NewIndexEmpty(cfg)

		tokenizer := tokenizer.InitTokenizer(indexInstance)

		absDocPath, err := file.AbsPath(docPath)
		if err != nil {
			panic(err)
		}

		if err := tokenizer.ParseMultiFile(absDocPath); err != nil {
			panic(err)
		}
		fmt.Printf("File Parsed. Found %d terms\n", indexInstance.Len())

		if dictPath != "" {
			//SAVE INDEX
			absDictPath, err := file.AbsPath(dictPath)
			if err != nil {
				panic(err)
			}
			if err = indexInstance.SaveIndex(absDictPath); err != nil {
				panic(err)
			}
		}
	} else if dictPath != "" {
		//READ INDEX
		indexInstance, err = index.NewIndexFromFile(cfg)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Dump loaded. Containing %d terms\n", indexInstance.Len())
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

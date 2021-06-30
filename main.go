package main

import (
	"fmt"
	"log"
	"strings"

	"g1.wpp2.hsnr/inr/vecret/config"

	"g1.wpp2.hsnr/inr/vecret/cli"
	"g1.wpp2.hsnr/inr/vecret/file"
	"g1.wpp2.hsnr/inr/vecret/index"
	"g1.wpp2.hsnr/inr/vecret/tokenizer"
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

		if cfg.UseSingleFileInput {
			if err := tokenizer.ParseSingleFile(absDocPath); err != nil {
				panic(err)
			}
		} else {
			if err := tokenizer.ParseMultiFile(absDocPath); err != nil {
				panic(err)
			}
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

	/*ctx := query.Context{
		Index:  indexInstance,
		Config: cfg,
	}
	parser := query.AstQueryParser{Context: &ctx}*/

	cl.Print("\r\n\r\n                                       .-') _  _  .-')   \r\n                                      ( OO ) )( \\( -O )  \r\n  ,----.     .-'),-----.   ,-.-') ,--./ ,--,'  ,------.  \r\n '  .-./-') ( OO'  .-.  '  |  |OO)|   \\ |  |\\  |   /`. ' \r\n |  |_( O- )/   |  | |  |  |  |  \\|    \\|  | ) |  /  | | \r\n |  | .--, \\\\_) |  |\\|  |  |  |(_/|  .     |/  |  |_.' | \r\n(|  | '. (_/  \\ |  | |  | ,|  |_.'|  |\\    |   |  .  '.' \r\n |  '--'  |    `'  '-'  '(_|  |   |  | \\   |   |  |\\  \\  \r\n  `------'       `-----'   `--'   `--'  `--'   `--' '--' \r\n\r\n")
	for {

		//Temporary Token Finder (No complex Queries, just one single token)
		/*token := cl.GetInput()
		pl := indexInstance.GetTermListCorrected(token)
		if pl == nil {
			fmt.Println("Not found")
			continue
		}
		if len(pl) == 0 {
			fmt.Println("No results found. Try to lower j")
			continue
		}
		fmt.Println(pl.String())*/

		in := cl.GetInput()

		res, err := indexInstance.FastCosine(in, 5)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(strings.Join(res, " "))

		/*q, err := parser.Parse(in)
		if err != nil {
			cl.Print(fmt.Sprintf("cannot parse query '%s': %s", in, err.Error()))
			continue
		}

		res, err := q.Evaluate()
		if err != nil {
			cl.Print(fmt.Sprintf("cannot execute query: %s", err.Error()))
			continue
		}
		fmt.Println(res.Entry.String(indexInstance))*/

	}

}

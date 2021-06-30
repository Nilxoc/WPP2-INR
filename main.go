package main

import (
	"fmt"
	"log"

	"g1.wpp2.hsnr/inr/boolret/query"

	"g1.wpp2.hsnr/inr/boolret/config"

	"g1.wpp2.hsnr/inr/boolret/cli"
	"g1.wpp2.hsnr/inr/boolret/file"
	"g1.wpp2.hsnr/inr/boolret/index"
	"g1.wpp2.hsnr/inr/boolret/tokenizer"
)

func w2vMain(cfg *config.Config) {
	fmt.Println("Running in Word2Vec Mode")

	w2vInd, _ := index.BuildIndex(cfg.PDoc, cfg.ModelPath)
	cl := cli.Init()
	cl.Print("\r\n\r\n                                       .-') _  _  .-')   \r\n                                      ( OO ) )( \\( -O )  \r\n  ,----.     .-'),-----.   ,-.-') ,--./ ,--,'  ,------.  \r\n '  .-./-') ( OO'  .-.  '  |  |OO)|   \\ |  |\\  |   /`. ' \r\n |  |_( O- )/   |  | |  |  |  |  \\|    \\|  | ) |  /  | | \r\n |  | .--, \\\\_) |  |\\|  |  |  |(_/|  .     |/  |  |_.' | \r\n(|  | '. (_/  \\ |  | |  | ,|  |_.'|  |\\    |   |  .  '.' \r\n |  '--'  |    `'  '-'  '(_|  |   |  | \\   |   |  |\\  \\  \r\n  `------'       `-----'   `--'   `--'  `--'   `--' '--' \r\n\r\n")

	cl.Print("Ready, please enter first query")
	for {
		res := w2vInd.EvaluateQuery(cl.GetInput())
		for i := 0; i < 10; i++ {
			fmt.Print(res[i].Doc)
			fmt.Print(", ")
		}
		fmt.Print("\n")
	}
}

func main() {
	cfg, err := config.Parse()
	if err != nil {
		log.Panic(err)
	}
	if cfg.Word2Vec {
		w2vMain(cfg)
		return
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

	ctx := query.Context{
		Index:  indexInstance,
		Config: cfg,
	}
	parser := query.AstQueryParser{Context: &ctx}

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

		q, err := parser.Parse(in)
		if err != nil {
			cl.Print(fmt.Sprintf("cannot parse query '%s': %s", in, err.Error()))
			continue
		}

		res, err := q.Evaluate()
		if err != nil {
			cl.Print(fmt.Sprintf("cannot execute query: %s", err.Error()))
			continue
		}
		fmt.Println(res.Entry.String(indexInstance))

	}

}

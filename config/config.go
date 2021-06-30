package config

import (
	"flag"
	"fmt"
)

// Config global application configuration
type Config struct {
	// PDoc path to the input file to parse
	PDoc string
	// PDict path to the prebuild dictionary file
	// if supplied the index will not be rebuild
	PDict string

	// KGram size for k-gram index terms
	KGram   uint8
	JThresh float32
	CSpell  bool
	// CSpellThresh threshold for spelling correction
	// active if less results than the specified value are found
	CSpellThresh int
	//Wether to use only doc_dump.txt instead of folder of textfiles
	UseSingleFileInput bool

	Word2Vec bool

	ModelPath string
}

func DefaultConfig() *Config {
	defaultCfg := Config{
		PDoc:               "",
		PDict:              "",
		KGram:              2,
		JThresh:            0.2,
		CSpell:             true,
		CSpellThresh:       5,
		UseSingleFileInput: false,
		Word2Vec:           false,
		ModelPath:          "",
	}
	return &defaultCfg
}

func Parse() (*Config, error) {
	defaultCfg := DefaultConfig()
	return parseConfig(defaultCfg)
}

func parseConfig(c *Config) (*Config, error) {
	flag.StringVar(&c.PDoc, "docs", "", "path to input folder")
	flag.StringVar(&c.PDict, "dict", "", "path to dict dump (skip index build) - if docSource also provided saves dict")
	k := 0
	flag.IntVar(&k, "k", 2, "k for kgram index")
	flag.IntVar(&c.CSpellThresh, "r", 5, "min count of terms to return")
	j := 0.0
	flag.Float64Var(&j, "j", 0.2, "jaccard threshold")
	flag.BoolVar(&c.CSpell, "correction", true, "enable error correction")
	var docPath string
	flag.StringVar(&docPath, "doc", "", "path to input file (txt)")
	flag.StringVar(&c.ModelPath, "model", "", "path to word2vec model")
	flag.BoolVar(&c.Word2Vec, "word2vec", false, "Switch to Word 2 Vec Retrieval")
	flag.Parse()

	c.KGram = uint8(k)
	c.JThresh = float32(j)

	if c.PDoc == "" && docPath != "" {
		c.PDoc = docPath
		c.UseSingleFileInput = true
	}

	if c.PDoc == "" && c.PDict == "" {
		return nil, fmt.Errorf("input file path or dictionary path required ")
	}

	return c, nil
}

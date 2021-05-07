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
}

func Parse() (*Config, error) {
	defaultCfg := Config{
		PDoc:         "",
		PDict:        "",
		KGram:        2,
		JThresh:      0.2,
		CSpell:       true,
		CSpellThresh: 5,
	}

	return parseConfig(&defaultCfg)
}

func parseConfig(c *Config) (*Config, error) {
	flag.StringVar(&c.PDoc, "doc", "", "path to input file (txt)")
	flag.StringVar(&c.PDict, "dict", "", "path to dict dump (skip index build) - if docSource also provided saves dict")
	k := 0
	flag.IntVar(&k, "k", 2, "k for kgram index")
	flag.IntVar(&c.CSpellThresh, "r", 5, "min count of terms to return")
	j := 0.0
	flag.Float64Var(&j, "j", 0.2, "jaccard threshold")
	flag.BoolVar(&c.CSpell, "correction", true, "enable error correction")
	flag.Parse()

	c.KGram = uint8(k)
	c.JThresh = float32(j)

	if c.PDoc == "" && c.PDict == "" {
		return nil, fmt.Errorf("input file path or dictionary path required ")
	}

	return c, nil
}

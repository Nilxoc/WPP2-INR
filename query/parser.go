package query

import (
	"g1.wpp2.hsnr/inr/boolret/config"
	"g1.wpp2.hsnr/inr/boolret/index"
)

// Query root interface to process user queries
type Query interface {
	// Evaluate evaluates the
	Evaluate() (*Result, error)
	Print()
}

type Result struct {
	Entry *index.PostingList
}

// Context common evaluation context
type Context struct {
	Index  *index.Index
	Config *config.Config
}

type Parser interface {
	Parse(s string) *Query
}

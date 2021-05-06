package query

import (
	"fmt"
	"g1.wpp2.hsnr/inr/boolret/index"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer/stateful"
	"github.com/alecthomas/repr"
)

type AstQuery struct {
	index *index.Index
	root  BinExpr
}

type AstQueryParser struct {
	index *index.Index
}

type astNode interface {
	nodeType() string
	terminal() bool
	evaluate(index *index.Index) (*index.PostingList, error)
}

func Evaluate(q *AstQuery) (*index.PostingList, error) {
	return nil, nil
}

// AST-Nodes

type Term struct {
	Group  *Group   `parser:"@@ |"`
	Expr   *BinExpr `parser:"@@ |"`
	Phrase *Phrase  `parser:"@@ |"`
	Value  *Value   `parser:"@@"`
}

//goland:noinspection GoVetStructTag
type Group struct {
	Expr *Term `"\"" @@* "\""`
}

type BinExpr struct {
	Left  *Value        `parser:"@@"`
	Right []*CondOpTerm `parser:"@@*"`
}

type CondOpTerm struct {
	Op    *string `parser:"@('AND'|'OR'|'AND NOT'|Proxim)"`
	Value *Value  `parser:"@@"`
}

type Value struct {
	Value *string `parser:"@Ident"`
}

//goland:noinspection GoVetStructTag
type Phrase struct {
	Value []*string `"\"" @@* "\""`
}

var lexer = stateful.MustSimple([]stateful.Rule{
	{`Ident`, `[a-zA-Z][a-zA-Z_\d]*`, nil},
	{`Proxim`, `\/\d`, nil},
	{"String", `"(\\"|[^"])*"`, nil},
	{"whitespace", `\s+`, nil}})

var parser = participle.MustBuild(&Term{},
	participle.Lexer(lexer),
	//participle.UseLookahead(1),
)

func Parse(query string) (*Term, error) {
	gram := &Term{}
	err := parser.ParseString("", query, gram)
	if err != nil {
		return nil, fmt.Errorf("cannot parse string: %v", err)
	}
	return gram, nil
}

func Print(r *Term) {
	repr.Println(r, repr.Indent(" "), repr.OmitEmpty(true))
}

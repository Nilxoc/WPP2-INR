package query

import (
	"fmt"
	"g1.wpp2.hsnr/inr/boolret/index"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer/stateful"
	"github.com/alecthomas/repr"
	"strconv"
	"strings"
)

type AstQueryParser struct {
	Context *Context
}

func (p AstQueryParser) Parse(s string) (*AstQuery, error) {
	gram := &Expression{}
	err := parser.ParseString("", s, gram)
	if err != nil {
		return nil, fmt.Errorf("cannot parse string: %v", err)
	}
	q := &AstQuery{context: p.Context, root: gram}
	return q, nil
}

type AstQuery struct {
	context *Context
	root    *Expression
}

func (q AstQuery) Evaluate() (*Result, error) {
	pl, err := q.root.eval(q.context)
	if err != nil {
		return nil, err
	}
	res := Result{Entry: pl}

	return &res, nil
}

func (q AstQuery) Print() {
	repr.Println(q.root, repr.Indent(" "), repr.OmitEmpty(true))
}

// AST-Nodes

type Phrase struct {
	Value []*string `"\"" @Ident* "\""`
}

type Value struct {
	Text          *string     `parser:"( @Ident"`
	Phrase        *string     `parser:"| @String)"`
	Subexpression *Expression `parser:"| '(' @@ ')'"`
}

type Term struct {
	Left  *Value   `parser:"@@"`
	Right []*Value `@@*`
}

type OpTerm struct {
	OP   *string `parser:"( @BoolOp |"`
	K    *string `parser:"@Proxim )"`
	Term *Term   `parser:"@@"`
}

type Expression struct {
	Left  *Term     `parser:"@@"`
	Right []*OpTerm `parser:"@@*"`
}

var lexer = stateful.MustSimple([]stateful.Rule{
	{"BoolOp", `(?i)AND\sNOT|AND|OR`, nil},
	{"Ident", `[a-zA-Z_]\w*`, nil},
	{"String", `"(\\"|[^"])*"`, nil},
	{`Proxim`, `/\d+`, nil}, // MUST be over 'Punct'
	{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`, nil},
	{"Number", `[-+]?(\d*\.)?\d+`, nil},
	{"EOL", `[\n\r]+`, nil},
	{"whitespace", `[ \t]+`, nil},
})

var parser = participle.MustBuild(&Expression{},
	participle.Lexer(lexer),
	//participle.UseLookahead(3),
)

func (e *Expression) eval(ctx *Context) (*index.PostingList, error) {
	left, err := e.Left.Eval(ctx)
	if err != nil {
		return nil, err
	}

	for _, r := range e.Right {
		eval, err := r.Term.Eval(ctx)
		if err != nil {
			return nil, err
		}

		left, err = r.eval(left, eval)
		if err != nil {
			return nil, err
		}
	}
	return left, nil
}

func (t *OpTerm) eval(l *index.PostingList, r *index.PostingList) (*index.PostingList, error) {
	op := t.OP
	if op != nil {
		return evalOper(op, l, r)
	}
	if t.K != nil {
		return l.Proximity(r, parseK(t.K)), nil
	}
	return nil, fmt.Errorf("unsupported operand: '%v'", op)
}

func evalOper(op *string, l *index.PostingList, r *index.PostingList) (*index.PostingList, error) {
	switch *op {
	case "OR":
		return l.Union(r), nil
	case "AND":
		return l.Intersect(r), nil
	case "AND NOT":
		return l.Difference(r), nil
	default:
		return nil, fmt.Errorf("unknown operator: '%s'", *op)
	}
}

func parseK(k *string) int64 {
	p, err := strconv.ParseInt(*k, 0, 32)
	if err != nil {
		panic(fmt.Errorf("cannot parse proximity: '%s'", *k))
	}
	return p
}

func (t *Term) Eval(ctx *Context) (*index.PostingList, error) {
	left, err := t.Left.Eval(ctx)
	if err != nil {
		return nil, err
	}

	for _, r := range t.Right {
		eval, err := r.Eval(ctx)
		if err != nil {
			return nil, err
		}
		return eval, nil
	}
	return left, nil
}

func (t *Value) Eval(ctx *Context) (*index.PostingList, error) {
	text := t.Text
	if text != nil {
		return getTerm(ctx, *text), nil
	}

	phrase := t.Phrase
	if phrase != nil {
		fields := strings.Fields(*phrase)

		fieldLen := len(fields)
		if fieldLen == 0 {
			empty := make(index.PostingList, 0)
			return &empty, nil
		} else if fieldLen == 1 {
			return getTerm(ctx, fields[0]), nil
		} else {
			terms := getTerms(ctx, fields)
			first := terms[0]

			others := make([]*index.PostingList, 0, fieldLen-1)
			for i := 1; i <= fieldLen; i++ {
				others[i] = getTerm(ctx, fields[i])
			}

			return first.PhraseIntersect(others), nil
		}
	}

	// TODO: sub-query
	return nil, nil
}

func getTerm(ctx *Context, text string) *index.PostingList {
	term := ctx.Index.GetTermListCorrected(text)
	if term != nil {
		return &term
	}
	empty := make(index.PostingList, 0)
	return &empty
}

func getTerms(ctx *Context, text []string) []*index.PostingList {
	res := make([]*index.PostingList, 0, len(text))
	for _, t := range text {
		term := getTerm(ctx, t)
		res = append(res, term)
	}
	return res
}

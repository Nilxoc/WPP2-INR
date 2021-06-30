package index

import "fmt"

type Document struct {
	DocID       int64
	TotalLength int
}

type Documents map[int64]*Document

type DocumentRef struct {
	Document  *Document
	TermCount int
}

type DocumentRefs []DocumentRef

func (d *Document) String() string {
	return fmt.Sprintf("%d", d.DocID)
}

func (d *Documents) String() string {
	ret := ""
	for _, e := range *d {
		ret += e.String() + " "
	}
	return ret
}

/*func (d *DocumentRef) String() string {
	return fmt.Sprintf("%s(%d)", d.Document.String(), d.TermCount)
}

func (d *DocumentRefs) String() string {
	ret := ""
	for _, e := range *d {
		ret += e.String() + " "
	}
	return ret
}*/

func (d *DocumentRef) String(idx *Index) string {
	return fmt.Sprintf("%s(%d)", idx.GetDocDisplay(d.Document.DocID), d.TermCount)
}

func (d *DocumentRefs) String(idx *Index) string {
	ret := ""
	for _, e := range *d {
		ret += e.String(idx) + " "
	}
	return ret
}

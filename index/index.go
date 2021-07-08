package index

import (
	"fmt"
	"log"
	"sort"

	"g1.wpp2.hsnr/inr/vecret/config"

	"g1.wpp2.hsnr/inr/vecret/file"
	"g1.wpp2.hsnr/inr/vecret/index/spell"
)

type IndexEntry struct {
	Term     string
	Docs     DocumentRefs
	DocCount int
}

type Index struct {
	Index     map[string]*IndexEntry
	K         int
	r         int
	j         float32
	kgram     *KGramIndex
	DocIDMap  map[int64]string
	Documents Documents

	AvgDocLength int
	DocCount     int
}

type SpellCorrectionResult struct {
	Replacements []string
	Docs         DocumentRefs
}

func NewIndexEmpty(c *config.Config) *Index {
	idx := &Index{
		K:            int(c.KGram),
		r:            c.CSpellThresh,
		j:            c.JThresh,
		DocIDMap:     make(map[int64]string),
		Documents:    make(Documents, 0),
		DocCount:     0,
		AvgDocLength: 0,
	}
	idx.Index = make(map[string]*IndexEntry)
	//idx.kgram = InitKGramIndex(idx.K)
	return idx
}

func NewIndexFromFile(cfg *config.Config) (*Index, error) {
	/*absDictPath, err := file.AbsPath(cfg.PDict)
	if err != nil {
		panic(err)
	}

	idx, err := loadIndex(absDictPath)
	if err != nil {
		return nil, err
	}

	idx.r = cfg.CSpellThresh
	idx.j = cfg.JThresh
	idx.K = int(cfg.KGram)
	idx.kgram = InitKGramIndex(idx.K)

	//Regenerate KGram-Index on Load, because we cannot persist it..
	for k, v := range idx.Index {
		idx.kgram.AddKGram(k, v)
	}
	return idx, nil*/
	return nil, fmt.Errorf("Not Supported for Vector Index!")
}

func (i *Index) AddDocument(documentId int64, documentLength int, docIDString string) *Document {
	doc := &Document{DocID: documentId, TotalLength: documentLength, TermRefs: make([]*DocumentRef, 0)}
	i.Documents[documentId] = doc
	if _, f := i.DocIDMap[documentId]; !f {
		i.DocIDMap[documentId] = docIDString
	}
	return doc
}

///INDEX SPECIFIC METHODS
func (i *Index) AddTerm(term string, count int, documentRef *Document) {

	if documentRef == nil {
		log.Fatalln("Got empty document Reference!")
	}

	if entry, found := i.Index[term]; found {
		entry.DocCount++
		docRef := &DocumentRef{TermCount: count, Document: documentRef, TermRef: entry}
		entry.Docs = append(entry.Docs, docRef)
		documentRef.TermRefs = append(documentRef.TermRefs, docRef)
	} else {
		t := IndexEntry{
			Term: term,
			Docs: make(DocumentRefs, 1),
		}
		docRef := &DocumentRef{TermCount: count, Document: documentRef, TermRef: &t}
		t.Docs[0] = docRef
		i.Index[term] = &t
		documentRef.TermRefs = append(documentRef.TermRefs, docRef)
		//i.kgram.AddKGram(term, &t)
	}
}

func docsSum(inp []*IndexEntry) int {
	sum := 0
	for _, i := range inp {
		sum += len(i.Docs)
	}
	return sum
}

func (i *Index) GetTerm(term string) *IndexEntry {
	if e, f := i.Index[term]; f {
		return e
	}
	return nil
}

func (i *Index) GetTermSuggestions(term string) []*IndexEntry {
	if term == "" {
		return nil
	}

	res := make([]*IndexEntry, 1)
	if e, f := i.Index[term]; f {
		res[0] = e
		docs := docsSum(res)
		if docs < i.r {
			//Not enough documents found
			res = append(res, i.getCorrectedDocs(term, i.r-docs)...)
		}
		return res
	}
	return i.getCorrectedDocs(term, i.r)
}

func dedupe(entries []*IndexEntry) []*IndexEntry {
	res := make([]*IndexEntry, 0, len(entries))
	contains := func(x *IndexEntry) bool {
		for _, e := range res {
			if e.Term == x.Term {
				return true
			}
		}
		return false
	}

	for _, e := range entries {
		if !contains(e) {
			res = append(res, e)
		}
	}
	return res
}

func (i *Index) getCorrectedDocs(term string, altCount int) []*IndexEntry {

	type Candidate struct {
		ldist int
		entry *IndexEntry
	}

	candidates := make([]Candidate, 0)

	posTokens := i.kgram.FindTokens(term)

	for _, k := range posTokens {
		if k.Entry.Term != term { // SKIP ALREADY SELECED TERM
			if jv := spell.Jaccard(term, k.Entry.Term, k.Count, i.K); jv > i.j {
				candidates = append(candidates, Candidate{
					ldist: spell.LevenshteinDistance(term, k.Entry.Term),
					entry: k.Entry,
				})
			}
		}
	}

	//Sort by Distance & by string
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].ldist == candidates[j].ldist {
			il := len(candidates[i].entry.Docs)
			jl := len(candidates[j].entry.Docs)
			if il == jl {
				return candidates[i].entry.Term < candidates[j].entry.Term
			}
			return il > jl
		}
		return candidates[i].ldist < candidates[j].ldist
	})

	minFunc := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	min := minFunc(len(candidates), altCount)
	totalAdd := 0

	res := make([]*IndexEntry, 0, min)
	for i := 0; i < min; i++ {
		//Refactor: Move to CLI handler fmt.Printf("Did you mean %s ?\n", candidates[i].entry.Term)
		res = append(res, candidates[i].entry)
		totalAdd += len(candidates[i].entry.Docs)
		if totalAdd >= altCount {
			//Got enough Documents to return
			break
		}
	}

	return res
}

func (i *Index) GetTermListCorrected(term string) *SpellCorrectionResult {
	tmp := i.GetTermSuggestions(term)
	if len(tmp) == 1 {
		return &SpellCorrectionResult{Docs: tmp[0].Docs, Replacements: []string{tmp[0].Term}}
	}
	tms := make([]string, 0)
	res := make(DocumentRefs, 0)
	for _, e := range tmp {
		res = append(res, e.Docs...)
		tms = append(tms, e.Term)
	}
	return &SpellCorrectionResult{Docs: res, Replacements: tms}
}

func (i *Index) SaveIndex(path string) error {
	return file.SaveIndex(i, path)
}

func (i *Index) Len() int {
	return len(i.Index)
}

func (i *IndexEntry) String(idx *Index) string {
	return i.Docs.String(idx)
}

func (i *Index) GetDocDisplay(id int64) string {
	if s, f := i.DocIDMap[id]; f {
		return s
	}
	return fmt.Sprintf("{%d}", id)
}

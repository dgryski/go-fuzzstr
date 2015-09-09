// Package fuzzstr implements a fuzzy string search in the style of Sublime Text
package fuzzstr

import "unicode/utf8"

// DocID is a document ID
type DocID uint32

//  Posting is a document and character position
type Posting struct {
	Doc DocID
	Pos uint32
}

// Index is a character index
type Index struct {
	postings  map[rune][]Posting
	allDocIDs []DocID
}

// NewIndex returns an index for the strings in docs
func NewIndex(docs []string) *Index {

	idx := Index{
		postings: make(map[rune][]Posting),
	}

	for id, d := range docs {
		docid := DocID(id)
		idx.allDocIDs = append(idx.allDocIDs, docid)
		for i, r := range d {
			idxr := idx.postings[r]
			idx.postings[r] = append(idxr, Posting{Doc: docid, Pos: uint32(i)})
		}
	}

	return &idx
}

// Query returns all documents which contain the letters in s in order
func (idx *Index) Query(s string) []Posting {
	r, w := utf8.DecodeRuneInString(s) //Error caught on insertion
	p := idx.postings[r]

	result := make([]Posting, len(p))

	for _, r := range s[w:] {
		result = intersect(result[:0], p, idx.postings[r])
		p = result
	}

	return p

}

// Filter removes from p all strings that additionally contain characters in s
func (idx *Index) Filter(p []Posting, s string) []Posting {
	result := make([]Posting, len(p))

	for _, r := range s {
		result = intersect(result[:0], p, idx.postings[r])
		p = result
	}

	return p
}

// intersect returns the intersection of two posting lists with the characters
// in b occuring later in the string than the entries in a
func intersect(result, a, b []Posting) []Posting {

	var aidx, bidx int

scan:
	for aidx < len(a) && bidx < len(b) {
		for a[aidx].Doc == b[bidx].Doc {

			if a[aidx].Pos < b[bidx].Pos {
				result = append(result, b[bidx])
				d := a[aidx].Doc
				for aidx < len(a) && d == a[aidx].Doc {
					aidx++
				}
				if aidx >= len(a) {
					break scan
				}
			}
			bidx++
			if bidx >= len(b) {
				break scan
			}
		}

		for a[aidx].Doc < b[bidx].Doc {
			aidx++
			if aidx >= len(a) {
				break scan
			}
		}

		for bidx < len(b) && a[aidx].Doc > b[bidx].Doc {
			bidx++
			if bidx >= len(b) {
				break scan
			}
		}
	}

	return result
}

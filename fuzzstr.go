// Package fuzzstr implements a fuzzy string search in the style of Sublime Text
package fuzzstr

// DocID is a document ID
type DocID uint32

//  Posting is a document and character position
type Posting struct {
	Doc DocID
	Pos uint32
}

// Index is a character index
type Index struct {
	postings  map[byte][]Posting
	allDocIDs []DocID
}

// NewIndex returns an index for the strings in docs
func NewIndex(docs []string) *Index {

	idx := Index{
		postings: make(map[byte][]Posting),
	}

	for id, d := range docs {
		docid := DocID(id)
		idx.allDocIDs = append(idx.allDocIDs, docid)
		for i, r := range []byte(d) {
			idxr := idx.postings[r]
			idx.postings[r] = append(idxr, Posting{Doc: docid, Pos: uint32(i)})
		}
	}

	return &idx
}

// Query returns all documents which contain the letters in s in order
func (idx *Index) Query(s string) []Posting {

	p := idx.postings[s[0]]

	for _, r := range []byte(s[1:]) {
		p = intersect(p, idx.postings[r])
	}

	return p

}

// Filter removes from p all strings that additionally contain characters in s
func (idx *Index) Filter(p []Posting, s string) []Posting {

	for _, r := range []byte(s) {
		p = intersect(p, idx.postings[r])
	}

	return p
}

// intersect returns the intersection of two posting lists with the characters
// in b occuring later in the string than the entries in a
func intersect(a, b []Posting) []Posting {

	var aidx, bidx int

	var result []Posting

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

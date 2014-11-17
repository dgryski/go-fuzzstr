package fuzzstr

type DocID uint32

type Posting struct {
	Doc DocID
	pos uint32
}

type Index struct {
	postings  map[byte][]Posting
	allDocIDs []DocID
}

// NewIndex returns an index for the strings in docs
func NewIndex(docs []string) Index {

	idx := Index{
		postings: make(map[byte][]Posting),
	}

	for id, d := range docs {
		docid := DocID(id)
		idx.allDocIDs = append(idx.allDocIDs, docid)
		for i, r := range []byte(d) {
			idxr := idx.postings[r]
			idx.postings[r] = append(idxr, Posting{Doc: docid, pos: uint32(i)})
		}
	}

	return idx
}

func (idx *Index) Query(s string) []Posting {

	var p []Posting = idx.postings[s[0]]

	for _, r := range []byte(s[1:]) {
		p = intersect(p, idx.postings[r])
	}

	return p

}

func intersect(a, b []Posting) []Posting {

	var aidx, bidx int

	var result []Posting

scan:
	for aidx < len(a) && bidx < len(b) {
		for a[aidx].Doc == b[bidx].Doc {

			if a[aidx].pos < b[bidx].pos {
				result = append(result, b[bidx])
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

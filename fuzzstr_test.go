package fuzzstr

import (
	"testing"
)

func TestQuery(t *testing.T) {

	docs := []string{
		"reprecipitation",
		"grallic",
		"fir",
		"emigrate",
		"cataphrenia",
		"unpostponed",
		"prerogativity",
		"chiefly",
		"hup",
		"unzealously",
		"goldilocks",
		"especial",
		"exoticness",
		"polymorphean",
		"chalcosine",
		"tutworkman",
		"labrosaurid",
		"compactness",
		"superannuate",
		"uranist",
	}

	idx := NewIndex(docs)

	tests := []struct {
		q     string
		words []string
	}{
		{
			"ac",
			[]string{"grallic", "chalcosine", "compactness"},
		},
		{
			"cac",
			[]string{"chalcosine", "compactness"},
		},
		{
			"zz",
			nil,
		},
		{
			"epi",
			[]string{"reprecipitation", "especial"},
		},
	}

	for _, tt := range tests {
		postings := idx.Query(tt.q)

		if len(postings) != len(tt.words) {
			t.Errorf("Query(%q)=[%d]string, want [%d]string", tt.q, len(postings), len(tt.words))
			for _, d := range postings {
				t.Log(docs[d.Doc])
			}
			continue
		}

		for i, p := range postings {
			if docs[p.Doc] != tt.words[i] {
				t.Errorf("Query(%q)[%d]=%q, want %q", tt.q, i, docs[p.Doc], tt.words[i])
				continue
			}
		}
	}
}

func TestFilter(t *testing.T) {

	docs := []string{
		"foobar",
		"bazbar",
		"foobaz",
		"bazfoo",
		"qux",
		"zot",
	}

	idx := NewIndex(docs)

	p := idx.Query("ba")
	p = idx.Filter(p, "fo")

	if len(p) != 1 || docs[p[0].Doc] != "bazfoo" {
		t.Errorf("Filter(Query(ba), fo)=%v, want `bazfoo`", p)
	}
}

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

	postings := idx.Query("cac")

	for _, p := range postings {
		t.Log(docs[p.Doc])
	}

	postings = idx.Query("ac")

	t.Log()
	for _, p := range postings {
		t.Log(docs[p.Doc])
	}
}

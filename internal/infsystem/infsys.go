package infsystem

import (
	"evaluation/internal/infsysresults"
	"evaluation/internal/qrels"
)

type InfSystem struct {
	Relevances qrels.RelevantDocs
	ISResults  infsysresults.ResultMap
}

func (i InfSystem) Results(infNeed string) (docsID []string) {
	res := i.ISResults[infNeed]
	for _, r := range res {
		docsID = append(docsID, r.DocID)
	}
	return docsID
}

func (i InfSystem) TotalRelevant(infNeed string) int {
	rels := i.Relevances[infNeed]
	total := 0
	for _, r := range rels {
		if i.IsRelevant(r.DocID, infNeed) {
			total++
		}
	}
	return total
}

func (i InfSystem) IsRelevant(docID, infNeed string) bool {
	return i.Relevances.IsRelevant(docID, infNeed)
}

func (i InfSystem) GetInfNeeds() []string {
	var infNeeds []string
	for k := range i.Relevances {
		infNeeds = append(infNeeds, k)
	}
	return infNeeds
}

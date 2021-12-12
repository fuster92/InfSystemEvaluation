package metrics

func Recall(col ResultCollectioner, infNeed string) float64 {
	res := col.Results(infNeed)
	rels := 0.0
	for _, r := range res {
		if col.IsRelevant(r, infNeed) {
			rels++
		}

	}
	return rels / float64(col.TotalRelevant(infNeed))
}

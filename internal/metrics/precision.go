package metrics

func PrecisionAt(col ResultCollectioner, infNeed string, k int) float64 {
	res := col.Results(infNeed)
	rels := 0.0
	for i := 0; i < k; i++ {
		if col.IsRelevant(res[i], infNeed) {
			rels++
		}
	}

	if len(res) < 10 {
		return rels / 10.0
	}
	return rels / float64(k)
}

func Precision(col ResultCollectioner, infNeed string) float64 {
	res := col.Results(infNeed)
	rels := 0
	for _, r := range res {
		if col.IsRelevant(r, infNeed) {
			rels++
		}
	}

	return float64(rels) / float64(len(res))
}

func AveragePrecision(col ResultCollectioner, infNeed string) float64 {
	res := col.Results(infNeed)
	avrg := 0.0
	rels := 0.0
	for i := 0; i < len(res); i++ {
		if col.IsRelevant(res[i], infNeed) {
			avrg += PrecisionAt(col, infNeed, i+1)
			rels++
		}
	}
	if rels == 0 {
		return 0.0
	}
	return avrg / rels
}

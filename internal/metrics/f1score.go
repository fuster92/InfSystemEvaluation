package metrics

func F1Score(prec, rec, beta float64) float64 {
	betaS := beta * beta
	num := (betaS + 1) * prec * rec
	den := betaS*prec + rec

	return num / den
}

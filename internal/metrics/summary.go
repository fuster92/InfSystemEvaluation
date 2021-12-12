package metrics

import (
	"evaluation/internal/plot"
	"fmt"
	"gonum.org/v1/plot/plotter"
	"math"
	"strings"
)

type ResultCollectioner interface {
	Results(infNeed string) []string
	TotalRelevant(infNeed string) int
	IsRelevant(docID, infNeed string) bool
	GetInfNeeds() []string
}

type RecPrecPoint struct {
	Rec, Prec float64
}

func (r RecPrecPoint) String() string {
	return fmt.Sprintf("%.3f %.3f", r.Rec, r.Prec)
}

type RPPoints []RecPrecPoint

func (r RPPoints) Len() int {
	return len(r)
}

func (r RPPoints) XY(p int) (x, y float64) {
	return r[p].Rec, r[p].Prec
}

func (r RPPoints) GetInterpolatedPrec(rec float64) (pre float64) {
	pre = 0.0
	for _, p := range r {
		if p.Rec < rec {
			continue
		}

		pre = math.Max(p.Prec, pre)
	}
	return pre
}

func (r RPPoints) String() string {
	var b strings.Builder
	for _, point := range r {
		b.WriteString(fmt.Sprintf("%s\n", point.String()))
	}
	return b.String()
}

type InfNeedSummary struct {
	Precision   float64
	Recall      float64
	F1          float64
	Prec10      float64
	AveragePrec float64
	RecPre      RPPoints
	IntRecPre   RPPoints
}

type InfSysSummary struct {
	InfNeedSummary
}

func RecallPrecision(col ResultCollectioner, infNeed string) (points []RecPrecPoint) {
	res := col.Results(infNeed)
	rels := 0.0
	for i := 0; i < len(res); i++ {
		if col.IsRelevant(res[i], infNeed) {
			p := PrecisionAt(col, infNeed, i+1)
			rels += 1
			points = append(points, RecPrecPoint{
				Rec:  rels / float64(col.TotalRelevant(infNeed)),
				Prec: p,
			})
		}
	}
	return points
}

func InterpolatedRecallPrecision(col ResultCollectioner, infNeed string) (points []RecPrecPoint) {
	rppoints := RPPoints(RecallPrecision(col, infNeed))
	for i := 0.0; i < 1.0; i += 0.1 {
		points = append(points, RecPrecPoint{
			Rec:  i,
			Prec: rppoints.GetInterpolatedPrec(i),
		})
	}
	return points
}

func (i InfNeedSummary) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s %.3f\n", "precision", i.Precision))
	b.WriteString(fmt.Sprintf("%s %.3f\n", "recall", i.Recall))
	b.WriteString(fmt.Sprintf("%s %.3f\n", "F1", i.F1))
	b.WriteString(fmt.Sprintf("%s %.3f\n", "prec@10", i.Prec10))
	b.WriteString(fmt.Sprintf("%s %.3f\n", "average_precision", i.AveragePrec))
	b.WriteString(fmt.Sprintf("%s\n", "recall_precision"))
	b.WriteString(fmt.Sprintf("%s\n", i.RecPre.String()))
	b.WriteString(fmt.Sprintf("%s\n", "interpolated_recall_precision"))
	b.WriteString(fmt.Sprintf("%s\n", i.IntRecPre.String()))
	return b.String()
}

type Summary map[string]InfNeedSummary

func (s Summary) String() string {
	var b strings.Builder
	for k, v := range s {
		if k != "TOTAL" {
			b.WriteString("INFORMATION_NEED " + k + "\n")
			b.WriteString(v.String())
		}
	}
	t := s["TOTAL"]
	b.WriteString("TOTAL \n")
	b.WriteString(fmt.Sprintf("precision %.3f\n", t.Precision))
	b.WriteString(fmt.Sprintf("recall %.3f\n", t.Recall))
	b.WriteString(fmt.Sprintf("F1 %.3f\n", t.F1))
	b.WriteString(fmt.Sprintf("prec@10 %.3f\n", t.Prec10))
	b.WriteString(fmt.Sprintf("MAP %.3f\n", t.AveragePrec))
	b.WriteString(fmt.Sprintf("%s\n", "interpolated_recall_precision"))
	b.WriteString(fmt.Sprintf("%s\n", t.IntRecPre.String()))
	b.WriteString("\n")
	return b.String()
}

func CreateSummary(c ResultCollectioner) Summary {
	m := make(map[string]InfNeedSummary)
	needs := c.GetInfNeeds()
	for _, need := range needs {
		m[need] = createInfSummary(c, need)
	}
	Summary(m).createTotal()
	plotMap := make(map[string]plotter.XYer)
	for k, need := range m {
		plotMap[k] = need.IntRecPre
	}
	plot.RecallPrecision("rp.png", plotMap)
	return m
}

func createInfSummary(c ResultCollectioner, need string) InfNeedSummary {
	ins := InfNeedSummary{}
	ins.Precision = Precision(c, need)
	ins.Prec10 = PrecisionAt(c, need, 10)
	ins.AveragePrec = AveragePrecision(c, need)
	ins.Recall = Recall(c, need)
	ins.F1 = F1Score(ins.Precision, ins.Recall, 1)
	ins.RecPre = RecallPrecision(c, need)
	ins.IntRecPre = InterpolatedRecallPrecision(c, need)
	return ins
}

func (s Summary) createTotal() {
	ins := InfNeedSummary{}
	ins.Precision, ins.Recall, _, ins.Prec10, ins.AveragePrec = s.getMeanPrecRecF1()
	ins.F1 = F1Score(ins.Precision, ins.Recall, 1)
	ins.IntRecPre = s.getMeanInterpolatedRecPre()
	s["TOTAL"] = ins
}

func (s Summary) getMeanPrecRecF1() (p, r, f1, p10, MAP float64) {

	for _, infneedsum := range s {
		p += infneedsum.Precision
		r += infneedsum.Recall
		f1 += infneedsum.F1
		p10 += infneedsum.Prec10
		MAP += infneedsum.AveragePrec
	}
	p = p / float64(len(s))
	r = r / float64(len(s))
	f1 = f1 / float64(len(s))
	p10 = p10 / float64(len(s))
	MAP = MAP / float64(len(s))
	return p, r, f1, p10, MAP
}
func (s Summary) getMeanInterpolatedRecPre() (points []RecPrecPoint) {

	for i := 0; i < 10; i++ {
		points = append(points, RecPrecPoint{})
	}
	for _, ins := range s {
		for i := 0; i < 10; i++ {
			points[i].Rec += ins.IntRecPre[i].Rec
			points[i].Prec += ins.IntRecPre[i].Prec
		}
	}
	for i := 0; i < 10; i++ {
		points[i].Prec = points[i].Prec / float64(len(s))
		points[i].Rec = points[i].Rec / float64(len(s))
	}
	return points
}

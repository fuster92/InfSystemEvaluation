package plot

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func RecallPrecision(filename string, rpPoints map[string]plotter.XYer) {

	p := plot.New()

	p.Title.Text = "Precision Recall Plot"
	p.X.Label.Text = "Recall"
	p.Y.Label.Text = "Precision"
	var plots []interface{}
	for name, xyers := range rpPoints {
		plots = append(plots, name)
		plots = append(plots, xyers)
	}
	err := plotutil.AddLinePoints(p, plots...)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(20*vg.Centimeter, 20*vg.Centimeter, filename); err != nil {
		panic(err)
	}
}

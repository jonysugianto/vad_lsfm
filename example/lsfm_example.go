package main

import (
	"fmt"
	"math/rand"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"

	"github.com/jonysugianto/vad_lsfm/lsfm"
)

func main() {
	var wavdata, lsfmvalues = lsfm.Voice2Lsfm("atas1.wav")
	fmt.Println(lsfmvalues)
	rand.Seed(int64(0))

	pwave, err := plot.New()
	if err != nil {
		panic(err)
	}

	pwave.Title.Text = "Plot Wav"
	pwave.X.Label.Text = "Time"
	pwave.Y.Label.Text = "Wav"

	err = plotutil.AddLines(pwave,
		"Wav", plotWav(wavdata))
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := pwave.Save(8*vg.Inch, 4*vg.Inch, "wave.png"); err != nil {
		panic(err)
	}

	plsfm, err := plot.New()
	if err != nil {
		panic(err)
	}

	plsfm.Title.Text = "Plot LSFM"
	plsfm.X.Label.Text = "Time"
	plsfm.Y.Label.Text = "LSFM"

	err = plotutil.AddLines(plsfm,
		"LSFM", plotLSFM(lsfmvalues))
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := plsfm.Save(8*vg.Inch, 4*vg.Inch, "lsfm.png"); err != nil {
		panic(err)
	}
}

func plotWav(wavData []float64) plotter.XYs {
	pts := make(plotter.XYs, len(wavData))
	var size = len(wavData)
	for i := 0; i < size; i++ {
		pts[i].X = float64(i)
		pts[i].Y = wavData[i]
	}
	return pts
}

func plotLSFM(lsfmvalues []float64) plotter.XYs {
	pts := make(plotter.XYs, len(lsfmvalues))
	var size = len(lsfmvalues)
	for i := 0; i < size; i++ {
		pts[i].X = (float64(i) * 80) + 1
		if lsfmvalues[i] < -2 {
			pts[i].Y = -2
		} else {
			pts[i].Y = lsfmvalues[i]
		}
	}
	return pts
}

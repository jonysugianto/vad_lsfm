package lsfm

import (
	"math"

	"github.com/jonysugianto/vad_lsfm/voicefeatures"
)

type LSFM struct {
	M                 int
	R                 int
	S                 []float64
	AM                []float64
	GM                []float64
	MEL_FILTER_STREAM [][]float64
	S_STREAM          [][]float64
}

func CreateLSFM(m int, r int) *LSFM {
	var ret = new(LSFM)
	ret.M = m
	ret.R = r
	ret.S = make([]float64, voicefeatures.NUMBEROFTRIANGULARFILTERS)
	ret.AM = make([]float64, voicefeatures.NUMBEROFTRIANGULARFILTERS)
	ret.GM = make([]float64, voicefeatures.NUMBEROFTRIANGULARFILTERS)
	for i := 0; i < voicefeatures.NUMBEROFTRIANGULARFILTERS; i++ {
		ret.GM[i] = 1
	}
	return ret
}

func (this *LSFM) addMelFilter(melfilter []float64) {
	this.MEL_FILTER_STREAM = append(this.MEL_FILTER_STREAM, melfilter)
	for i := 0; i < voicefeatures.NUMBEROFTRIANGULARFILTERS; i++ {
		this.S[i] = this.S[i] + melfilter[i]
	}
	if len(this.MEL_FILTER_STREAM) > this.M {
		var firstmelfilter = this.MEL_FILTER_STREAM[0]
		for i := 0; i < voicefeatures.NUMBEROFTRIANGULARFILTERS; i++ {
			this.S[i] = this.S[i] - firstmelfilter[i]
		}
		this.MEL_FILTER_STREAM = this.MEL_FILTER_STREAM[1:]
	}
}

func (this *LSFM) compS() {
	var number_M = this.M
	if len(this.MEL_FILTER_STREAM) < number_M {
		number_M = len(this.MEL_FILTER_STREAM)
	}
	var newS []float64 = make([]float64, voicefeatures.NUMBEROFTRIANGULARFILTERS)
	for i := 0; i < voicefeatures.NUMBEROFTRIANGULARFILTERS; i++ {
		newS[i] = this.S[i] / float64(number_M)
	}

	this.S_STREAM = append(this.S_STREAM, newS)
	for i := 0; i < voicefeatures.NUMBEROFTRIANGULARFILTERS; i++ {
		this.AM[i] = this.AM[i] + newS[i]
		this.GM[i] = this.GM[i] * newS[i]
	}
	if len(this.S_STREAM) > this.R {
		var firstS = this.S_STREAM[0]
		for i := 0; i < voicefeatures.NUMBEROFTRIANGULARFILTERS; i++ {
			var pembagi = firstS[i]
			if pembagi < 0.00001 {
				pembagi = 0.00001
			}

			this.AM[i] = this.AM[i] - firstS[i]
			this.GM[i] = this.GM[i] / pembagi
		}
		this.S_STREAM = this.S_STREAM[1:]
	}
}

func (this *LSFM) compAM() []float64 {
	var number_R = this.R
	if len(this.S_STREAM) < number_R {
		number_R = len(this.S_STREAM)
	}
	var ret []float64 = make([]float64, voicefeatures.NUMBEROFTRIANGULARFILTERS)
	for i := 0; i < voicefeatures.NUMBEROFTRIANGULARFILTERS; i++ {
		ret[i] = this.AM[i] / float64(number_R)
	}
	return ret
}

func (this *LSFM) compGM() []float64 {
	var number_R = this.R
	if len(this.S_STREAM) < number_R {
		number_R = len(this.S_STREAM)
	}
	var ret []float64 = make([]float64, voicefeatures.NUMBEROFTRIANGULARFILTERS)
	for i := 0; i < voicefeatures.NUMBEROFTRIANGULARFILTERS; i++ {
		ret[i] = math.Pow(this.GM[i], 1.0/float64(number_R))
	}
	return ret
}

func (this *LSFM) compL() []float64 {
	var ret []float64 = make([]float64, voicefeatures.NUMBEROFTRIANGULARFILTERS)
	var am = this.compAM()
	var gm = this.compGM()
	for i := 0; i < voicefeatures.NUMBEROFTRIANGULARFILTERS; i++ {
		var pembagi = am[i]
		if pembagi < 0.00001 {
			pembagi = 0.00001
		}
		ret[i] = math.Log10(gm[i] / pembagi)
	}
	return ret
}

func (this *LSFM) compSingleL() float64 {
	var ret float64
	var am = this.compAM()
	var gm = this.compGM()
	for i := 0; i < voicefeatures.NUMBEROFTRIANGULARFILTERS; i++ {
		var pembagi = am[i]
		if pembagi < 0.00001 {
			pembagi = 0.00001
		}
		ret += math.Log10(gm[i] / pembagi)
	}
	return ret
}

func (this *LSFM) CompLsfm(melfilter []float64) float64 {
	this.addMelFilter(melfilter)
	this.compS()
	return this.compSingleL()
}

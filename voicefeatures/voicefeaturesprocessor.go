package voicefeatures

import (
	"fmt"
	"math"
	"math/cmplx"

	"github.com/jonysugianto/vad_lsfm/common"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/spectral"
	"github.com/mjibson/go-dsp/window"
)

var ALPHA_PREEMPHASIS float64 = 0.97
var MEL_MIN_FREQUENCY = float64(0)
var MEL_MAX_FREQUENCY = float64(8000)
var NUMBEROFTRIANGULARFILTERS int = 31
var CENTEROFFREQUENCY []float64
var FREQUENCY_STEP float64

var SAMPLERATE = float64(16000)
var FRAMELENGTH int = 400
var FRAMEPERIODE int = 80
var NOVERLAP = FRAMELENGTH - FRAMEPERIODE
var HANNWINDOW []float64

func init() {
	fmt.Println("init voicefeatures")
	CENTEROFFREQUENCY = calculateCentersOfFrequency()
	FREQUENCY_STEP = (SAMPLERATE / 2) / (float64(FRAMELENGTH) / 2)
	NOVERLAP = FRAMELENGTH - FRAMEPERIODE
	HANNWINDOW = window.Hann(FRAMELENGTH)
}

func CompPreEmphasis(signal []float64) []float64 {
	var ret []float64
	ret = append(ret, signal[0])
	var size = len(signal)
	for i := 1; i < size; i++ {
		ret = append(ret, signal[i]-ALPHA_PREEMPHASIS*signal[i-1])
	}
	return ret
}

func CompHannWindowSignal(signal []float64) []float64 {
	var ret = make([]float64, FRAMELENGTH)
	for i := 0; i < FRAMELENGTH; i++ {
		ret[i] = HANNWINDOW[i] * signal[i]
	}
	return ret
}

func SplittSignalIntoFrames(data []float64) [][]float64 {
	var segmented = spectral.Segment(data, FRAMELENGTH, NOVERLAP)
	return common.MapArray(segmented, CompHannWindowSignal)
}

func CompFft(data [][]float64) [][]float64 {
	var fftdata = fft.FFT2Real(data)
	var row = len(fftdata)
	var ret [][]float64
	for i := 0; i < row; i++ {
		var col = len(fftdata[i])
		var temp = make([]float64, col)
		for j := 0; j < col; j++ {
			temp[j] = cmplx.Abs(fftdata[i][j])
		}
		ret = append(ret, temp)
	}
	return ret
}

func fromMelToHz(mels float64) float64 {
	return 700 * (math.Pow(10, mels/2595) - 1)
}

func fromHzToMel(hzs float64) float64 {
	return 2595 * math.Log10(1+hzs/700)
}

func calculateCentersOfFrequency() []float64 {
	var delta = (fromHzToMel(MEL_MAX_FREQUENCY) - fromHzToMel(MEL_MIN_FREQUENCY)) / float64(NUMBEROFTRIANGULARFILTERS+1)
	var centerOfFrequency = make([]float64, NUMBEROFTRIANGULARFILTERS+2)
	centerOfFrequency[0] = MEL_MIN_FREQUENCY
	for i := 1; i <= NUMBEROFTRIANGULARFILTERS; i++ {
		centerOfFrequency[i] = fromMelToHz(delta*float64(i) + fromHzToMel(MEL_MIN_FREQUENCY))
		i += 1
	}
	centerOfFrequency[NUMBEROFTRIANGULARFILTERS+1] = MEL_MAX_FREQUENCY
	return centerOfFrequency
}

func getTriangleFilterCoefficient(triangleIndex int, frequencyIndex int) float64 {
	var frequency = float64(frequencyIndex) * FREQUENCY_STEP
	var previousCenterOfFrequency = CENTEROFFREQUENCY[triangleIndex-1]
	var currentCenterOfFrequency = CENTEROFFREQUENCY[triangleIndex]
	var nextCenterOfFrequency = CENTEROFFREQUENCY[triangleIndex+1]
	if previousCenterOfFrequency <= frequency && frequency < currentCenterOfFrequency {
		return (frequency - previousCenterOfFrequency) / (currentCenterOfFrequency - previousCenterOfFrequency)
	} else if currentCenterOfFrequency <= frequency && frequency < nextCenterOfFrequency {
		return (frequency - nextCenterOfFrequency) / (currentCenterOfFrequency - nextCenterOfFrequency)
	} else {
		return 0
	}
}

func CompMelfilter(input []float64) []float64 {
	var numberOfTriangularFilters = len(CENTEROFFREQUENCY) - 2
	var result = make([]float64, numberOfTriangularFilters)
	for i := 1; i <= numberOfTriangularFilters; i++ {
		var melCoefficient = float64(0)
		var size = int(len(input) / 2)
		for j := 0; j < size; j++ {
			melCoefficient += input[j] * getTriangleFilterCoefficient(i, j)
		}
		result[i-1] = melCoefficient
	}
	return result
}

func compLog(input []float64) []float64 {
	var size = len(input)
	var result = make([]float64, size)
	for i := 0; i < size; i++ {
		result[i] = math.Log1p(input[i] + 1)
	}
	return result
}

func CompLogMelfilter(input []float64) []float64 {
	var mf = CompMelfilter(input)
	return compLog(mf)
}

package lsfm

import (
	//	"github.com/jonysugianto/vad_lsfm/common"
	"github.com/jonysugianto/vad_lsfm/io"
	"github.com/jonysugianto/vad_lsfm/voicefeatures"
)

var M int = 5
var R int = 10

func Voice2Lsfm(filename string) ([]float64, []float64) {
	var lsfm = CreateLSFM(M, R)
	var wavdata = io.ReadWav(filename)
	var framesdata = voicefeatures.SplittSignalIntoFrames(wavdata)
	var fftdata = voicefeatures.CompFft(framesdata)
	//	var melfiltereddata = common.MapArray(fftdata, voicefeatures.CompMelfilter)
	var numberframe = len(framesdata)
	var lsfmdata []float64
	for i := 0; i < numberframe; i++ {
		var lsfmvalue = lsfm.CompLsfm(fftdata[i])
		lsfmdata = append(lsfmdata, lsfmvalue)
	}
	return wavdata, lsfmdata
}

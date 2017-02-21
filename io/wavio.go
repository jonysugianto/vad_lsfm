package io

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/mjibson/go-dsp/wav"
	"io"
	"os"
)

func ReadWav(filename string) []float64 {
	wavfile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer wavfile.Close()
	wavobj, err := wav.New(wavfile)
	if wavobj == nil {
		fmt.Println("Error reading wav from file '%s', nil returned", filename)
	}
	if err != nil {
		fmt.Println(err)
	}
	var retfloat32, _=wavobj.ReadFloats(wavobj.Samples)
	var size=len(retfloat32)
	var ret=make([]float64, size)
	for i:=0; i<size; i++{
		ret[i]=float64(retfloat32[i])
	}
	return ret
}

func write(w io.Writer, data interface{}) {
	if b, ok := data.([]byte); ok {
		for c := 0; c < len(b); {
			n, err := w.Write(b[c:])
			if err != nil {
				panic(err)
			}
			c += n
		}
		return
	}
	if err := binary.Write(w, binary.LittleEndian, data); err != nil {
		panic(err)
	}
}

func writeChunk(w io.Writer, id string, data []byte) (err error) {
	if len(id) != 4 {
		panic(errors.New("invalid chunk id"))
	}
	write(w, []byte(id))
	write(w, uint32(len(data)))
	write(w, data)
	return
}

func writeFmt(w io.Writer, SampleRate uint32, SignificantBits uint16, Channels uint16) (err error) {
	var b bytes.Buffer
	write(&b, uint16(1)) // uncompressed/PCM
	write(&b, Channels)
	write(&b, SampleRate)
	write(&b, uint32(Channels)*SampleRate*uint32(SignificantBits)/8) // bytes per second
	write(&b, SignificantBits/8*Channels)                            // block align
	write(&b, SignificantBits)
	return writeChunk(w, "fmt ", b.Bytes())
}

func WriteData(w io.Writer, data []byte, SampleRate uint32, SignificantBits uint16, Channels uint16) (err error) {
	defer func() {
		if e, ok := recover().(error); ok {
			err = e
		}
	}()
	var buf bytes.Buffer
	writeFmt(&buf, SampleRate, SignificantBits, Channels)
	writeChunk(&buf, "data", data)
	write(w, []byte("RIFF"))
	write(w, uint32(buf.Len()))
	write(w, []byte("WAVE"))
	write(w, buf.Bytes())
	return
}

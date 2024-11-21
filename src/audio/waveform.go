package audio

import (
	"fmt"

	waveform "github.com/cettoana/go-waveform"
)

func GetWaveForm(wavFileBytes []byte) {
	w := waveform.DecodeWav(wavFileBytes)
	fmt.Println(w.BitsPerSample)

	data, _ := w.GetData()

	if stereoData, ok := data.(*waveform.StereoData); ok {
		fmt.Println(stereoData.RSample)
		fmt.Println(stereoData.LSample)
	} else if monoData, ok := data.(*waveform.MonoData); ok {
		fmt.Println(monoData.Sample)
	}
}

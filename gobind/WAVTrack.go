package main

import (
	"fmt"
	"time"
)

type WAVTrack struct {
	BaseTrack
	sampleRateHz int
}

// =============== WAVTrack Unique Methods  ====================
func (w *WAVTrack) SampleRate() int {
	return w.sampleRateHz
}

// ================ Overwrite Base :Play(): Method ===============
func (t *WAVTrack) Play() {
	fmt.Printf("\n🎶 WAV# playing %s by %s (duration: %v)\n", t.title, t.artist, t.duration)
}

// ============ Constructor Create :WAVTrack: ===============
func NewWAVTrack(title, artist string, duration time.Duration, sampleRate int) *WAVTrack {
	return &WAVTrack{
		BaseTrack: BaseTrack{
			title:    title,
			artist:   artist,
			duration: duration,
			format:   WAV,
		},
		sampleRateHz: sampleRate,
	}
}

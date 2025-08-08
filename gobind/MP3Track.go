package main

import (
	"fmt"
	"time"
)

type MP3Track struct {
	BaseTrack
	bitrateKbps int
}

// ================== MP3 Specific Methods =======================
func (mp3 *MP3Track) Bitrate() int {
	return mp3.bitrateKbps
}

// ==================== Override Play Method =======================
func (t *MP3Track) Play() {
	fmt.Printf("\n🎶 MP3# playing %s by %s (duration: %v)\n", t.title, t.artist, t.duration)
}

// ============ MP3 Constructor ============
func NewMP3Track(title, artist string, duration time.Duration, bitrate int) *MP3Track {
	return &MP3Track{
		BaseTrack: BaseTrack{
			title:    title,
			artist:   artist,
			duration: duration,
			format:   MP3, // hardcoded
		},
		bitrateKbps: bitrate,
	}
}

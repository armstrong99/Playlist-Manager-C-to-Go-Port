package main

import (
	"fmt"
	"time"
)

const (
	MP3 = "mp3"
	WAV = "wav"
)

type Track interface {
	Title() string
	Artist() string
	Duration() time.Duration
	Format() string
	Play()
}

type BaseTrack struct {
	title    string
	artist   string
	duration time.Duration
	format   string
}

func (t *BaseTrack) Title() string {
	return t.title
}

func (t *BaseTrack) Artist() string {
	return t.artist
}

func (t *BaseTrack) Duration() time.Duration {
	return t.duration
}
func (t *BaseTrack) Format() string {
	return t.format
}

func (t *BaseTrack) Play() {
	fmt.Printf("🎶 playing %s by %s (duration: %v)", t.title, t.artist, t.duration)
}

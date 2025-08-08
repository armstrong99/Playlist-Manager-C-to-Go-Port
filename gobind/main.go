package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	var tracks []Track
	homePlaylist := NewPlaylist("Feeling good")
	workPlaylist := NewPlaylist("At work")
	exercisePlaylist := NewPlaylist("At the gym")

	playlist, err := exercisePlaylist.LoadFromDisk("At the gym")

	if err != nil {
		fmt.Printf("Error Playing from disk: %v", err)

	} else {
		playlist.PlayAll()
	}

	for i := range 10 {
		if i%2 == 0 {
			track := NewMP3Track(
				fmt.Sprintf("Track-%v", i),
				fmt.Sprintf("Artist-%v", i),
				time.Duration((i+3)*int(time.Second)),
				rand.Intn(i+1-0+i)+i,
			)
			tracks = append(tracks, track)
		} else {
			track := NewWAVTrack(
				fmt.Sprintf("W:Track-%v", i),
				fmt.Sprintf("W:Artist-%v", i),
				time.Duration((i+3)*int(time.Second)),
				rand.Intn(i+1-0+i)+i,
			)
			tracks = append(tracks, track)
		}
	}

	for i := range 5 {
		exercisePlaylist.AddTrack(tracks[i])

		if i%2 == 0 {
			homePlaylist.AddTrack(tracks[i])
			continue
		}
		workPlaylist.AddTrack(tracks[i])
	}

	homePlaylist.PlayAll()

	workPlaylist.SaveToDisk("")
	exercisePlaylist.SaveToDisk("")
}

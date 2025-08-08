package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	SAVE_FOLDER = "golang-playlists"
)

type Playlist struct {
	name   string
	tracks []Track
}

func (p *Playlist) ensureFolderExists() error {
	fullPath := filepath.Join(".", SAVE_FOLDER)

	info, err := os.Stat(fullPath)

	if os.IsExist(err) && !info.IsDir() {
		return os.Mkdir(SAVE_FOLDER, os.ModePerm)
	}

	if os.IsNotExist(err) {
		return os.Mkdir(SAVE_FOLDER, os.ModePerm)
	}

	// some other error
	return err
}

func (p *Playlist) AddTrack(track Track) {

	p.tracks = append(p.tracks, track)
}

func (p *Playlist) RemoveTrack(index int) {
	if index < len(p.tracks) && index >= 0 {
		p.tracks = slices.Delete(p.tracks, index, index+1)
	}
}

func (p *Playlist) PlayAll() {
	if len(p.tracks) < 1 {
		fmt.Printf("No track in this playlist")
		return
	}
	fmt.Printf("\n============ 🎸🎺🥁 Playing playlist: %s (duration: %v)========= \n", p.name, p.GetTotalDuration())

	for _, t := range p.tracks {
		fmt.Printf("\n🎶 🎸  Now playing: %s, by %s (%v | %s | %s)\n", t.Title(), t.Artist(), t.Duration(), t.Format(), p.name)
		t.Play()
		time.Sleep(t.Duration())
		fmt.Printf("\nFinished playing: %s by %s\n", t.Title(), t.Artist())
		time.Sleep(500 * time.Millisecond)
		t = nil
	}
}

func (p *Playlist) ShuffleTracks() {
	if len(p.tracks) < 1 {
		p.tracks = []Track{}
		return
	}
	rand.Shuffle(len(p.tracks), func(i, j int) {
		p.tracks[i], p.tracks[j] = p.tracks[j], p.tracks[i]
	})
}

func (p *Playlist) GetTotalDuration() time.Duration {
	if len(p.tracks) < 1 {
		fmt.Printf("No track in this playlist")
		return 0
	}
	var totalTime time.Duration

	for _, t := range p.tracks {

		totalTime += t.Duration()
	}

	return totalTime
}

func (p *Playlist) SaveToDisk(filenameHint string) error {
	if err := p.ensureFolderExists(); err != nil {
		return fmt.Errorf("failed to create folder: %v", err)
	}

	safeName := p.name
	if filenameHint != "" {
		safeName = filenameHint
	}

	playListPath := filepath.Join(".", SAVE_FOLDER, safeName+".txt")

	file, err := os.Create(playListPath)

	if err != nil {
		return fmt.Errorf("seems file already exist with this path(%s): %v", playListPath, err)
	}

	defer file.Close()

	if _, err := fmt.Fprintln(file, p.name); err != nil {
		return err
	}

	for _, track := range p.tracks {
		extra := "0"

		switch t := track.(type) {
		case *MP3Track:
			extra = strconv.Itoa(t.Bitrate())
		case *WAVTrack:
			extra = strconv.Itoa(t.SampleRate())
		}

		_, wErr := fmt.Fprintf(file, "%s|%s|%v|%s|%s\n",
			track.Title(),
			track.Artist(),
			track.Duration().Milliseconds(),
			track.Format(),
			extra,
		)
		if wErr != nil {
			return wErr
		}

	}
	fmt.Printf("\nPlaylist saved to: %s\n", playListPath)
	return nil
}

func (p *Playlist) GetSavedPlaylists() ([]string, error) {
	var names []string

	if er := p.ensureFolderExists(); er != nil {
		return []string{}, er
	}

	entries, err := os.ReadDir(SAVE_FOLDER)

	if err != nil {
		return []string{}, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".txt") {
			name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
			names = append(names, name)
		}
	}

	return names, nil
}

func (p *Playlist) LoadFromDisk(name string) (*Playlist, error) {
	if er := p.ensureFolderExists(); er != nil {
		return nil, er
	}

	playlistPath := filepath.Join(".", SAVE_FOLDER, name+".txt")

	file, err := os.Open(playlistPath)

	if err != nil {
		return nil, fmt.Errorf("seems this file may not exists: %v", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		return nil, errors.New("invalid playlist: empty file")
	}

	playlistName := strings.TrimSpace(scanner.Text())
	if playlistName == "" {
		return nil, errors.New("invalid playlist file or name")
	}

	yourPlayList := &Playlist{name: playlistName}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")

		if len(parts) != 5 {
			continue
		}

		title := parts[0]
		artist := parts[1]

		duration, err := strconv.Atoi(parts[2])
		if err != nil {
			fmt.Printf("Duration for this track is undefined, skipping track ->")
			continue
		}

		format := parts[3]
		extra, err := strconv.Atoi(parts[4])
		if err != nil {
			fmt.Printf("Failed to read metadata for this track, skipping track ->")
			continue
		}

		var track Track

		switch format {
		case MP3:
			track = &MP3Track{BaseTrack: BaseTrack{
				title:    title,
				format:   format,
				artist:   artist,
				duration: time.Duration(duration) * time.Millisecond,
			},
				bitrateKbps: extra,
			}
		case WAV:
			track = &WAVTrack{BaseTrack: BaseTrack{
				title:    title,
				format:   format,
				artist:   artist,
				duration: time.Duration(duration) * time.Millisecond,
			},
				sampleRateHz: extra,
			}
		default:
			track = &BaseTrack{title: title, artist: artist, format: format, duration: time.Duration(duration) * time.Millisecond}
		}

		yourPlayList.tracks = append(yourPlayList.tracks, track)

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return yourPlayList, nil
}

func (p *Playlist) Name() string {
	return p.name
}

func (p *Playlist) Size() int {
	return len(p.tracks)
}

func (p *Playlist) GetFolder() string {
	return SAVE_FOLDER
}

// ============ Constructor :New:Playlist: =====================
func NewPlaylist(name string) *Playlist {
	p := &Playlist{
		name:   name,
		tracks: []Track{},
	}
	p.ensureFolderExists()

	return p
}

package player

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"go.senan.xyz/taglib"
)

type Player struct {
	Streamer    beep.Streamer
	Format      beep.Format
	Control     *beep.Ctrl
	Metadata    map[string][]string
	CurrentFile string
}

func NewPlayer(filePath string) (*Player, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file %s: %v", filePath, err)
	}

	var streamer beep.Streamer
	var format beep.Format

	// Deteksi format audio berdasarkan ekstensi
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".mp3":
		streamer, format, err = mp3.Decode(file)
	case ".wav":
		streamer, format, err = wav.Decode(file)
	case ".flac":
		streamer, format, err = flac.Decode(file)
	default:
		file.Close()
		return nil, fmt.Errorf("format tidak didukung: %s", ext)
	}
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("gagal decode %s: %v", filePath, err)
	}

	// Inisialisasi speaker
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("gagal inisialisasi speaker: %v", err)
	}

	// Baca metadata
	metadata, _ := taglib.ReadTags(filePath)

	// Buat kontrol untuk play/pause
	ctrl := &beep.Ctrl{Streamer: streamer, Paused: true}
	speaker.Play(ctrl)

	return &Player{
		Streamer:    streamer,
		Format:      format,
		Control:     ctrl,
		Metadata:    metadata,
		CurrentFile: filePath,
	}, nil
}

func (p *Player) Play() {
	speaker.Lock()
	p.Control.Paused = false
	speaker.Unlock()
}

func (p *Player) Pause() {
	speaker.Lock()
	p.Control.Paused = true
	speaker.Unlock()
}

func (p *Player) Close() {
	speaker.Clear()
}

func (p *Player) GetMetadata() (title, artist, album string) {
	if p.Metadata == nil {
		return "Unknown", "Unknown", "Unknown"
	}
	title = "Unknown"
	if titles, ok := p.Metadata["TITLE"]; ok && len(titles) > 0 {
		title = titles[0]
	}
	artist = "Unknown"
	if artists, ok := p.Metadata["ARTIST"]; ok && len(artists) > 0 {
		artist = artists[0]
	}
	album = "Unknown"
	if albums, ok := p.Metadata["ALBUM"]; ok && len(albums) > 0 {
		album = albums[0]
	}
	return
}

func (p *Player) LoadNewTrack(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("gagal membuka file %s: %v", filePath, err)
	}

	var streamer beep.Streamer
	var format beep.Format

	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".mp3":
		streamer, format, err = mp3.Decode(file)
	case ".wav":
		streamer, format, err = wav.Decode(file)
	case ".flac":
		streamer, format, err = flac.Decode(file)
	default:
		file.Close()
		return fmt.Errorf("format tidak didukung: %s", ext)
	}
	if err != nil {
		file.Close()
		return fmt.Errorf("gagal decode %s: %v", filePath, err)
	}

	speaker.Lock()
	p.Control.Streamer = streamer
	p.Format = format
	p.CurrentFile = filePath
	p.Metadata, _ = taglib.ReadTags(filePath)
	speaker.Unlock()

	return nil
}

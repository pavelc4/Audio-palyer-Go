package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"audio-player-go/config"
	"audio-player-go/player"
	"audio-player-go/scanner"
	"audio-player-go/ui"
	"github.com/gdamore/tcell/v2"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Gagal memuat konfigurasi: %v", err)
		cfg = &config.Config{} // Gunakan config kosong jika gagal
	}

	var audioPath string

	if cfg.LastDir == "" {
		fmt.Printf("Masukkan path folder audio (contoh: /home/user/musik atau C:\\Music) [%s]: ", cfg.LastDir)
		scannerReader := bufio.NewScanner(os.Stdin)
		scannerReader.Scan()
		audioPath = strings.TrimSpace(scannerReader.Text())

		if audioPath == "" {
			log.Fatal("Path folder audio tidak boleh kosong")
		}
	} else {
		audioPath = cfg.LastDir
	}

	// Pindai folder untuk file audio
	tracks, err := scanner.ScanAudioFiles(audioPath)
	if err != nil {
		log.Fatalf("Gagal memindai folder: %v", err)
	}
	if len(tracks) == 0 {
		fmt.Printf("Tidak ada file audio (.mp3, .wav, .flac) ditemukan di %s\n", audioPath)
		os.Exit(1)
	}

	// Simpan direktori terakhir
	cfg.LastDir = audioPath

	if err := cfg.Save(); err != nil {
		log.Printf("Gagal menyimpan konfigurasi: %v", err)
	}

	// Inisialisasi UI terminal
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Gagal inisialisasi UI: %v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("Gagal inisialisasi screen: %v", err)
	}
	defer screen.Fini()

	// Inisialisasi player dengan track pertama
	audioPlayer, err := player.NewPlayer(tracks[0])
	if err != nil {
		log.Fatalf("Gagal inisialisasi player: %v", err)
	}
	defer audioPlayer.Close()

	// Jalankan UI
	ui.Run(screen, audioPlayer, tracks)
}

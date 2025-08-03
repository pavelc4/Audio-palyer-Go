package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ScanAudioFiles(rootPath string) ([]string, error) {
	var audioFiles []string

	// Bersihkan path input
	rootPath = strings.TrimSpace(rootPath)
	if rootPath == "" {
		return nil, fmt.Errorf("path tidak boleh kosong")
	}

	// Pastikan folder ada
	info, err := os.Stat(rootPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("folder %s tidak ditemukan", rootPath)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s bukan folder", rootPath)
	}

	// Pindai folder dan subfolder
	err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// Cek ekstensi audio
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".mp3" || ext == ".wav" || ext == ".flac" {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			audioFiles = append(audioFiles, absPath)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("gagal memindai folder: %v", err)
	}

	return audioFiles, nil
}

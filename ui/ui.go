package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"audio-player-go/player"
	"github.com/gdamore/tcell/v2"
)

func Run(screen tcell.Screen, audioPlayer *player.Player, tracks []string) {
	currentTrack := 0
	commandMode := false
	commandBuffer := ""
	isPlaying := false
	eventChan := make(chan tcell.Event)

	go func() {
		for {
			eventChan <- screen.PollEvent()
		}
	}()

	for {
		// Bersihkan layar
		screen.Clear()

		// Tampilkan UI
		drawText(screen, 0, 0, "Music Player by pavelc4")
		drawText(screen, 0, 2, "Tracks:")
		for i, track := range tracks {
			prefix := "  "
			if i == currentTrack {
				prefix = "> "
			}
			drawText(screen, 0, 3+i, prefix+filepath.Base(track))
		}

		// Tampilkan metadata
		title, artist, album := audioPlayer.GetMetadata()
		drawText(screen, 0, len(tracks)+4, fmt.Sprintf("Judul: %s", title))
		drawText(screen, 0, len(tracks)+5, fmt.Sprintf("Artis: %s", artist))
		drawText(screen, 0, len(tracks)+6, fmt.Sprintf("Album: %s", album))

		// Tampilkan mode command
		if commandMode {
			drawText(screen, 0, len(tracks)+9, fmt.Sprintf("Command: %s", commandBuffer))
		} else {
			drawText(screen, 0, len(tracks)+9, "Gunakan arrow key, Enter untuk play/pause, Spasi untuk pause, // untuk command mode")
		}

		// Tampilkan mode command
		if commandMode {
			drawText(screen, 0, len(tracks)+9, fmt.Sprintf("Command: %s", commandBuffer))
		} else {
			drawText(screen, 0, len(tracks)+9, "Gunakan arrow key, Enter untuk play/pause, Spasi untuk pause, // untuk command mode")
		}

		screen.Show()

		ev := <-eventChan
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if commandMode {
				if ev.Key() == tcell.KeyEnter {
					cmd := strings.TrimSpace(commandBuffer)
					if cmd == "q" || cmd == "quit" {
						return
					}
					commandBuffer = ""
					commandMode = false
				} else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
					if len(commandBuffer) > 0 {
						commandBuffer = commandBuffer[:len(commandBuffer)-1]
					}
				} else if ev.Rune() != 0 {
					commandBuffer += string(ev.Rune())
				}
			} else {
				switch ev.Key() {
				case tcell.KeyUp:
					if currentTrack > 0 {
						currentTrack--
					}
				case tcell.KeyDown:
					if currentTrack < len(tracks)-1 {
						currentTrack++
					}
				case tcell.KeyEnter:
					if audioPlayer.CurrentFile == tracks[currentTrack] {
						if isPlaying {
							audioPlayer.Pause()
						} else {
							audioPlayer.Play()
						}
						isPlaying = !isPlaying
					} else {
						go func() {
							if err := audioPlayer.LoadNewTrack(tracks[currentTrack]); err != nil {
								// Handle error
							} else {
								isPlaying = true
								audioPlayer.Play()
							}
						}()
					}
				case tcell.KeyRune:
					if ev.Rune() == ' ' {
						audioPlayer.Pause()
						isPlaying = false
					} else if ev.Rune() == '/' {
						if len(commandBuffer) == 0 {
							commandBuffer = "/"
						} else if commandBuffer == "/" {
							commandMode = true
							commandBuffer = ""
						}
					}
				case tcell.KeyEscape:
					commandBuffer = ""
				}
			}
		}
	}
}

func drawText(screen tcell.Screen, x, y int, text string) {
	for i, ch := range text {
		screen.SetContent(x+i, y, ch, nil, tcell.StyleDefault)
	}
}

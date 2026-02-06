package main

import (
	"golang.org/x/term"
	"os"
	"time"
)

const downLine = '│'
const sideLine = '─'
const BLOCK = '█'
const EMPTY = ' '

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	input := make(chan byte, 60)
	go captureInput(input)

	position := piecePosition{}
	spawnPiece(&position)

	rotationMatrix := createRotationMatrix(90)
	screen := initScreen()
	drawScreenToTerminal(screen)
	raw := loadConfig("controls.json")
	inpuutMapping := convert(raw)

	tickTime := time.Duration(400)
	ticker := time.NewTicker(tickTime * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case _ = <-ticker.C:
			fallPiece(&position, screen)
			drawScreenToTerminal(screen)
		default:
			time.Sleep(10 * time.Millisecond)
			handleInput(inpuutMapping, statePlaying, input, &position, &rotationMatrix, screen)
		}
	}
}

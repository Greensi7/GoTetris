package main

import (
	"bufio"
	"math/rand"
	//"errors"
	"fmt"
	"math"
	"os"
	"time"
	"unicode"

	"golang.org/x/term"
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

	position := piecePosition{
		cords: [4]coordinate{
			{y: -2, x: 5},
			{y: -1, x: 5},
			{y: 0, x: 5},
			{y: 0, x: 6},
		},
	}
	rotationMatrix := createRotationMatrix(90)
	fmt.Println(rotationMatrix)
	var screen [maxVertical][maxHorizontal]rune
	temp := 0
	clearScreen(&screen)
	var oldPosition *piecePosition = nil
	for {
		select {
		case i := <-input:
			if i == 3 {
				panic(1)
			}
			lower := byte(unicode.ToLower(rune(i)))
			f, ok := mapping[lower]
			if ok {
				pieceCopy := position
				oldPosition = &pieceCopy
				f(&rotationMatrix, &position)
				drawPiece(&position, oldPosition, &screen)
				drawScreenToTerminal(&screen)
			}

		default:
		}
		temp++
		if temp%8 == 0 {
			temp = 1
			pieceCopy := position
			oldPosition = &pieceCopy
			fallPiece(&position)
			drawPiece(&position, oldPosition, &screen)
			drawScreenToTerminal(&screen)
		}
		time.Sleep(30 * time.Millisecond)
	}
}

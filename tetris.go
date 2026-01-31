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
const minHorizontal = 0
const minVertical = 0
const maxHorizontal = 13
const maxVertical = 23

type gameAction func(*[2][2]int, *piecePosition)

type coordinate struct {
	x int
	y int
}

type piecePosition struct {
	cords [4]coordinate
}

func createRotationMatrix(angle float64) [2][2]int {
	//only use for int rotation matrixes
	// can't implement checking for int matrix because floats
	// 90 180 270 will break on any other input most likely
	rotationMatrix := [2][2]int{}
	sinAngle := math.Round(math.Sin(angle * (math.Pi / 180.0)))
	cosAngle := math.Round(math.Cos(angle * (math.Pi / 180.0)))
	rotationMatrix[0][0] = int(cosAngle)
	rotationMatrix[0][1] = int(sinAngle)
	rotationMatrix[1][0] = -int(sinAngle)
	rotationMatrix[1][1] = int(cosAngle)
	return rotationMatrix
}
func pushFromSide(position *piecePosition) {
	min := position.cords[0].x
	max := position.cords[0].x
	minFlag := false
	maxFlag := false
	for _, cord := range position.cords {
		if cord.x >= (maxHorizontal - 1) {
			maxFlag = true
		}
		if cord.x <= minHorizontal {
			minFlag = true
		}
		if cord.x > max {
			max = cord.x
		}
		if cord.x < min {
			min = cord.x
		}
	}
	if maxFlag {
		moveHorizontal(position, -int(math.Abs(float64(maxHorizontal-1-max))))
	}
	if minFlag {
		moveHorizontal(position, int(math.Abs(float64(minHorizontal-min))))
	}
}

func moveHorizontal(position *piecePosition, offset int) piecePosition {
	result := *position
	for index := range position.cords {
		result.cords[index].x += offset
	}
	return result
}

func isOnBoard(cord *coordinate) bool {
	if cord == nil {
		return false
	}
	if cord.x < (minHorizontal+1) || cord.x > (maxHorizontal-2) ||
		cord.y > (maxVertical-2) {
		return false
	}
	return true
}

func drawPiece(position *piecePosition, oldPosition *piecePosition, screen *[maxVertical][maxHorizontal]rune) {
	if position == nil || screen == nil {
		panic(1)
	}

	if oldPosition != nil {
		for _, pieceBlockCord := range oldPosition.cords {
			if pieceBlockCord.y >= (maxVertical - 1) {
				panic(1)
			}
			screen[pieceBlockCord.y][pieceBlockCord.x] = EMPTY
		}
	}

	for _, pieceBlockCord := range position.cords {

		if !isOnBoard(&pieceBlockCord) {
			panic(1)
		}

		if pieceBlockCord.y >= 0 {
			screen[pieceBlockCord.y][pieceBlockCord.x] = BLOCK
		}
	}
}

func clearScreen(screen *[maxVertical][maxHorizontal]rune) {
	for index := range screen {
		for index2 := range screen[index] {
			if (index % (maxVertical - 1)) == 0 {
				screen[index][index2] = sideLine
			} else if (index2 % (maxHorizontal - 1)) == 0 {
				screen[index][index2] = downLine
			} else {
				screen[index][index2] = ' '
			}
		}
	}
}

func inputMoveLeft(_ *[2][2]int, position *piecePosition) {
	moved := moveHorizontal(position, -1)
	position.cords = moved.cords
}

func inputMoveRight(_ *[2][2]int, position *piecePosition) {
	moved := moveHorizontal(position, 1)
	position.cords = moved.cords
}

func inputInterupt(_ *[2][2]int, _ *piecePosition) {
	panic(1)
}

func rotatePiece(rotationMatrix *[2][2]int, position *piecePosition) {
	offset := position.cords[1]
	for i, cord := range position.cords {
		position.cords[i].x = rotationMatrix[0][0]*(cord.x-offset.x) + rotationMatrix[0][1]*(cord.y-offset.y) + offset.x
		position.cords[i].y = rotationMatrix[1][0]*(cord.x-offset.x) + rotationMatrix[1][1]*(cord.y-offset.y) + offset.y
	}
}

func fallPiece(position *piecePosition) {
	for index := range position.cords {
		position.cords[index].y += 1
	}
}

func captureInput(ch chan byte) {
	defer close(ch)
	reader := bufio.NewReaderSize(os.Stdin, 1)
	for {
		input, err := reader.ReadByte()
		if err != nil {
			break
		}
		ch <- input
	}
}

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
	mapping := map[byte]gameAction{
		3:   inputInterupt,
		'k': rotatePiece,
		'h': inputMoveLeft,
		'l': inputMoveRight,
	}
	rotationMatrix := createRotationMatrix(90)
	fmt.Println(rotationMatrix)
	var screen [maxVertical][maxHorizontal]rune
	for {
		clearScreen(&screen)
		drawPiece(&position, nil, &screen)
		drawScreenToTerminal(&screen)
		select {
		case i := <-input:
			if i == 3 {
				panic(1)
			}
			lower := byte(unicode.ToLower(rune(i)))
			f, ok := mapping[lower]
			if ok {
				f(&rotationMatrix, &position)
			}

		default:
		}
		fallPiece(&position)
		time.Sleep(100 * time.Millisecond)
	}
}
func drawScreenToTerminal(screen *[maxVertical][maxHorizontal]rune) {
	clearTerminal()
	for rowIndex, row := range screen {
		for columnIndex, pixel := range row {
			if pixel == EMPTY {
				if (columnIndex+1)%2 == 0 && (rowIndex)%2 == 0 {
					fmt.Printf("%s%c%c%s", "\x1b[1;30m", BLOCK, BLOCK, "\x1b[1;39m")
				} else {
					fmt.Printf("%s%c%c%s", "\x1b[1;90m", BLOCK, BLOCK, "\x1b[1;39m")
				}
			} else if pixel == BLOCK {
				min := 31
				max := 36
				num := rand.Intn(max-min+1) + min
				blockColor := fmt.Sprintf("\x1b[1;%dm", num)
				fmt.Printf("%s%c%c%s", blockColor, BLOCK, BLOCK, "\x1b[1;39m")
			} else {
				fmt.Printf("%c%c", pixel, pixel)
			}
		}
		fmt.Printf("\r\n")
	}
}

func clearTerminal() {
	fmt.Print("\033[H\033[2J")
}

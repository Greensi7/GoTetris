package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"os"
	"unicode"
)

// TODO implement unique validator check
type RawPlaying struct {
	MoveRight   string `json:"right" validate:"len=1,ascii,required"`
	MoveLeft    string `json:"left" validate:"len=1,ascii,required"`
	RotatePiece string `json:"rotate" validate:"len=1,ascii,required"`
	Quit        string `json:"quit" validate:"len=1,ascii,required"`
	Pause       string `json:"pause" validate:"len=1,ascii,required"`
}

type RawMapping struct {
	Playing RawPlaying `json:"Playing" validate:"required"`
}

type gameAction func(*[2][2]int, *piecePosition)
type gameState string

const (
	statePlaying gameState = "Playing"
	statePaused  gameState = "Paused"
)

var playingMapping = map[string]gameAction{
	"Quit":        inputInterupt,
	"RotatePiece": rotatePiece,
	"MoveLeft":    inputMoveLeft,
	"MoveRight":   inputMoveRight,
}

func loadConfig(fileName string) RawMapping {
	raw, err := os.ReadFile(fileName)
	if err != nil {
		errMessage := fmt.Sprintf("Please create %s in project root directory.", fileName)
		fmt.Println(errMessage)
		panic(1)
	}
	rawMapping := RawMapping{}
	json.Unmarshal(raw, &rawMapping)
	validate := validator.New()
	err = validate.Struct(rawMapping)
	if err != nil {
		fmt.Println(err)
		panic(1)
	}
	return rawMapping
}

func convert(input RawMapping) map[gameState]map[byte]gameAction {
	return map[gameState]map[byte]gameAction{
		statePlaying: {
			byte(input.Playing.MoveLeft[0]):    playingMapping["MoveLeft"],
			byte(input.Playing.MoveRight[0]):   playingMapping["MoveRight"],
			byte(input.Playing.RotatePiece[0]): playingMapping["RotatePiece"],
			byte(input.Playing.Quit[0]):        playingMapping["Quit"],
		},
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



type inputHandler func(inputMapping map[byte]gameAction,
	ch chan byte,
	position *piecePosition,
	rotationMatrix *[2][2]int,
	screen [][]rune)

var inputHandlerMapping = map[gameState]inputHandler{
	statePlaying: handleInputPlaying,
	statePaused:  handleInputPaused,
}

func handleInput(inputMapping map[gameState]map[byte]gameAction,
	state gameState,
	ch chan byte,
	position *piecePosition,
	rotationMatrix *[2][2]int,
	screen [][]rune) {
	f, ok := inputHandlerMapping[state]
	if !ok {
		panic(1)
	}
	f(inputMapping[state], ch, position, rotationMatrix, screen)
}

func handleInputPlaying(inputMapping map[byte]gameAction,
	input chan byte,
	position *piecePosition,
	rotationMatrix *[2][2]int,
	screen [][]rune) {
	select {
	case i := <-input:
		f, ok := inputMapping[i]
		if !ok {
			lower := byte(unicode.ToLower(rune(i)))
			f, ok = inputMapping[lower]
		}
		if ok {
			eraserPiece(position, screen)
			positionCopy := *position
			f(rotationMatrix, &positionCopy)
			if isValidPos(&positionCopy, screen) {
				*position = positionCopy
				drawPiece(&positionCopy, screen)
				drawScreenToTerminal(screen)
			} else {
				drawPiece(position, screen)
			}
		}
	default:
	}
}

func handleInputPaused(inputMapping map[byte]gameAction,
	input chan byte,
	_ *piecePosition,
	_ *[2][2]int,
	screen [][]rune) {
}

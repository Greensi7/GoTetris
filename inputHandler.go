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
			byte(input.Playing.Quit[0]):        playingMapping["RotatePiece"],
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

func handleInput(inputMapping map[gameState]map[byte]gameAction,
	input chan byte,
	state gameState,
	position piecePosition,
	rotationMatrix [2][2]int) *piecePosition {
	select {
	case i := <-input:
		f, ok := inputMapping[state][i]
		if !ok {
			lower := byte(unicode.ToLower(rune(i)))
			f, ok = inputMapping[state][lower]
		}
		if ok {
			positionCopy := position
			f(&rotationMatrix, &positionCopy)
			return &positionCopy
		}
	default:
	}
	return nil
}

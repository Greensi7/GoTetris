package main

import (
	"bufio"
	"os"
	"fmt"
	"encoding/json"
	"github.com/go-playground/validator/v10"
)

type RawPlaying struct {
	MoveRight string `json:"right" validate:"len=1,ascii,required"`
	MoveLeft string `json:"left" validate:"len=1,ascii,required"`
	RotatePiece string `json:"rotate" validate:"len=1,ascii,required"`
	Quit string `json:"quit" validate:"len=1,ascii,required"`
	Pause string `json:"pause" validate:"len=1,ascii,required"`
}

type RawMapping struct{
	Playing RawPlaying `json:"Playing" validate:"required"`
}


type gameAction func(*[2][2]int, *piecePosition)
var mapping = map[string]gameAction{
	"Quit":   inputInterupt,
	"RotatePiece": rotatePiece,
	"MoveLeft": inputMoveLeft,
	"MoveRight": inputMoveRight,
}

func loadConfig(){
	raw, err := os.ReadFile("controls.json")
	if err != nil{
		fmt.Println("Please create controls.json in root directory(probably GoTetris).")
		panic(1)
	}
	rawMapping := RawMapping{}
	json.Unmarshal(raw, &rawMapping)
	validate := validator.New()
	err = validate.Struct(rawMapping)
	if err != nil{
		fmt.Println(err)
		panic(1)
	}
	fmt.Println(rawMapping)
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

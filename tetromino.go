package main

import (
	"math"
)

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
		if cord.x >= (height - 1) {
			maxFlag = true
		}
		if cord.x <= minWidth {
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
		moveHorizontal(position, -int(math.Abs(float64(max-1-height))))
	}
	if minFlag {
		moveHorizontal(position, int(math.Abs(float64(min-minWidth+1))))
	}
}

func moveHorizontal(position *piecePosition, offset int) {
	for index := range position.cords {
		position.cords[index].x += offset
	}
}

func inputMoveLeft(_ *[2][2]int, position *piecePosition) {
	moveHorizontal(position, -1)
}

func inputMoveRight(_ *[2][2]int, position *piecePosition) {
	moveHorizontal(position, 1)
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
	pushFromSide(position)
}

func spawnPiece(position *piecePosition){
	newPosition := piecePosition{
		cords: [4]coordinate{
			{y: -2, x: 5},
			{y: -1, x: 5},
			{y: 0, x: 5},
			{y: 1, x: 5},
		},
	}
	position.cords = newPosition.cords
}

func fallPiece(position *piecePosition, screen [][]rune) {
	positionCopy := *position
	for i := range positionCopy.cords {
		positionCopy.cords[i].y += 1
	}
	eraserPiece(position, screen)
	if !isValidPos(&positionCopy, screen) {
		drawPiece(position, screen)
		if isEnd(position){
			panic(1)
		}
		clearRows(position, screen)
		spawnPiece(position)
	}

	for i := range position.cords {
		position.cords[i].y += 1
	}
	drawPiece(position, screen)
}

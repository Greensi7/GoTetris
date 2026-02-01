package main

const width = 13
const height = 23


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
			if !isOnBoard(&pieceBlockCord) {
				panic(1)
			}
			if pieceBlockCord.y > 0 {
				screen[pieceBlockCord.y][pieceBlockCord.x] = EMPTY
			}
			if pieceBlockCord.y == 0 {
				screen[pieceBlockCord.y][pieceBlockCord.x] = sideLine
			}
		}
		oldPosition = nil
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

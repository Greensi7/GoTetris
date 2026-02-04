package main

//const width = 13
//const height = 23

const height = 24
const width = 13
const minWidth = 0
const minHeight = 0


func isOnBoard(cord *coordinate) bool {
	if cord == nil {
		return false
	}
	if cord.x < (minWidth+1) || cord.x > (width-2) ||
		cord.y > (height-2) {
		return false
	}
	return true
}

func eraserPiece(position *piecePosition, screen [][]rune){
	if position == nil || screen == nil {
		panic(1)
	}

	if position != nil {
		for _, pieceBlockCord := range position.cords {
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
		position = nil
	}
}

func isValidPos(position *piecePosition, screen [][]rune)bool{
	for _, pieceBlockCord := range position.cords {

		if !isOnBoard(&pieceBlockCord) {
			panic(1)
		}

		if pieceBlockCord.y >= 0 && screen[pieceBlockCord.y][pieceBlockCord.x] == BLOCK{
			return false
		}
	}
	return true
}

func drawPiece(position *piecePosition, screen [][]rune) {
	if position == nil || screen == nil {
		panic(1)
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

func clearScreen(screen [][]rune) {
	for index := range screen {
		for index2 := range screen[index] {
			if (index % (height - 1)) == 0 {
				screen[index][index2] = sideLine
			} else if (index2 % ( - 1)) == 0 {
				screen[index][index2] = downLine
			} else {
				screen[index][index2] = ' '
			}
		}
	}
}

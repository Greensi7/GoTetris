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

func eraserPiece(position *piecePosition, screen [][]rune) {
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

func isEnd(position *piecePosition) bool {
	for _, pieceBlockCord := range position.cords {

		if pieceBlockCord.y <= 0 {
			return true
		}
	}
	return false
}

func isValidPos(position *piecePosition, screen [][]rune) bool {
	for _, pieceBlockCord := range position.cords {

		if !isOnBoard(&pieceBlockCord) {
			return false
		}

		if pieceBlockCord.y >= 0 && screen[pieceBlockCord.y][pieceBlockCord.x] == BLOCK {
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

func clearRows(position *piecePosition, screen [][]rune) {
	for _, cord := range position.cords {
		if !isOnBoard(&cord) {
			panic(1)
		}

		if cord.y > 0 && isFull(cord.y, screen) {
			screen = append((screen)[:cord.y], (screen)[cord.y+1:]...)
			screen = append((screen)[:1], append([][]rune{nil}, (screen)[1:]...)...)
			initRow(screen, 1)
			drawScreenToTerminal(screen)
		}
	}
}

func isFull(row int, screen [][]rune) bool {
	if row < 0 || row > (height-1) {
		panic(1)
	}

	blocksCount := 0
	for _, segment := range screen[row] {
		if segment == BLOCK {
			blocksCount++
		}
	}
	return (width - 2) == blocksCount
}

func initScreen() [][]rune {
	screen := make([][]rune, height)

	for row := range screen {
		initRow(screen, row)
	}
	return screen
}

func initRow(screen [][]rune, row int) {
	screen[row] = make([]rune, width)
	for col := range screen[row] {
		if (row % (height - 1)) == 0 {
			screen[row][col] = sideLine
		} else if (col % (width - 1)) == 0 {
			screen[row][col] = downLine
		} else {
			screen[row][col] = EMPTY
		}
	}
}

package main

import(
	"fmt"
	"math/rand"
)

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

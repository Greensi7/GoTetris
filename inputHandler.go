package main
import(
	"bufio"
	"os"
)

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

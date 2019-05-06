package main

import (
	"fmt"
	"os"

	"github.com/matheusd/nonota"
	nonotaui "github.com/matheusd/nonota/ui"
)

func main() {
	filename := "nonota-board.yml"
	board, err := nonota.BoardFromFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ui := nonotaui.New(board, filename)

	err = ui.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/matheusd/nonota"
	nonotaui "github.com/matheusd/nonota/ui"
)

type opts struct {
	Previous bool `long:"previous" description:"See the board for the previous month"`
}

func getCmdOpts() *opts {
	cmdOpts := &opts{}
	parser := flags.NewParser(cmdOpts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		e, ok := err.(*flags.Error)
		if ok && e.Type == flags.ErrHelp {
			os.Exit(0)
		}
		fmt.Printf("Argument error: %v\n", e)
		os.Exit(1)
	}

	return cmdOpts
}

func main() {
	opts := getCmdOpts()

	refTime := time.Now()
	if opts.Previous {
		// Go back a day prior to the start of the current period
		// to get a date in the previous billing period
		refTime = nonota.StartOfBilling(nonota.StartOfBilling(refTime).Add(time.Hour * -24))
	}

	filename := "nonota-board.yml"
	board, err := nonota.BoardFromFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ui := nonotaui.New(board, filename, refTime)

	err = ui.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

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
	Filename string `short:"f" long:"filename" description:"Filename of the board to use"`
	Previous bool   `long:"previous" description:"See the board for the previous month"`
	Date     string `long:"date" description:"Year and month to generate billing (in the YYYY-MM format)"`
}

func getCmdOpts() *opts {
	cmdOpts := &opts{
		Filename: "nonota-board.yml",
	}
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
	} else if opts.Date != "" {
		ym, err := time.Parse("2006-01", opts.Date)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		refTime = nonota.StartOfBilling(ym)
	}

	board, err := nonota.BoardFromFile(opts.Filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ui := nonotaui.New(board, opts.Filename, refTime)

	err = ui.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

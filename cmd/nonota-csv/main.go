package main

import (
	"github.com/matheusd/nonota"
	"os"
	"time"
	"fmt"

	flags "github.com/jessevdk/go-flags"
)

type opts struct {
	Current bool `long:"current" description:"Generate for the current month"`
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

	ref := time.Now()
	if !opts.Current {
		// Go back a day prior to the start of the current period
		// to get a date in the previous billing period
		ref = nonota.StartOfBilling(nonota.StartOfBilling(ref).Add(time.Hour*-24))
	}
	start := nonota.StartOfBilling(ref)
	end := nonota.EndOfBilling(ref)
	dtFormat := "2006-01-02 15:04:05"

	filename := "nonota-board.yml"
	board, err := nonota.BoardFromFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tasks := make([]*nonota.Task, 0)
	for _, l := range board.Lists {
		for _, t := range l.Tasks {
			taskTime := t.TotalTime(start, end)
			if taskTime <= 0 {
				continue
			}

			tasks = append(tasks, t)
		}
	}

	var totTime time.Duration
	// Collected all relevant tasks. Output csv.
	for _, t := range tasks {
		csvFmt := "\"%s\",%.2f\n"
		taskTime := t.TotalTime(start, end)
		fmt.Printf(csvFmt, t.Title, taskTime.Hours())
		totTime += taskTime
	}

	fmt.Fprintf(os.Stderr, "\nGenerated CSV between %s and %s\n",
		start.Format(dtFormat), end.Format(dtFormat))
	fmt.Fprintf(os.Stderr, "Total computed time: %s (%.2f hours)\n", totTime,
		totTime.Hours())
}
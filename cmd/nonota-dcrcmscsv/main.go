package main

import (
	"fmt"
	"github.com/matheusd/nonota"
	"os"
	"regexp"
	"strings"
	"time"

	flags "github.com/jessevdk/go-flags"
)

type opts struct {
	Current bool    `long:"current" description:"Generate for the current month"`
	Rate    float64 `long:"rate" description:"Contractor rate in USD/hour"`
	Domain  string  `long:"domain" description:"Default domain for expenses"`
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

var reExtractSub *regexp.Regexp

func init() {
	var err error
	reExtractSub, err = regexp.Compile("#(\\S*)\\s?")
	if err != nil {
		panic(fmt.Errorf("error compiling regexp: %s", err))
	}
}

func extractSubdomain(s string) string {
	match := reExtractSub.FindAllStringSubmatch(s, -1)
	if len(match) < 1 {
		return ""
	}
	if len(match[0]) < 2 {
		return ""
	}
	return match[0][1]
}

func quote(s string) string {
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\"", "'")
	return s
}

func main() {
	opts := getCmdOpts()

	if opts.Domain == "" {
		fmt.Println("Specify the default billing domain")
		os.Exit(1)
	}

	if opts.Rate <= 0 {
		fmt.Println("Specify the rate in USD/hour")
		os.Exit(1)
	}

	ref := time.Now()
	if !opts.Current {
		// Go back a day prior to the start of the current period
		// to get a date in the previous billing period
		ref = nonota.StartOfBilling(nonota.StartOfBilling(ref).Add(time.Hour * -24))
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
	var totExpense float64

	// Collected all relevant tasks. Output csv.
	for _, t := range tasks {
		// TODO: Extract domain from task list
		typ := 1
		domain := opts.Domain
		subdomain := extractSubdomain(t.Title)
		descr := quote(t.Title)
		if t.Description != "" {
			descr += "\\n\\n" + quote(t.Description)
		}
		token := ""
		taskTime := t.TotalTime(start, end)
		labor := taskTime.Hours()
		expense := labor * opts.Rate

		csvFmt := "%d\t%s\t%s\t%s\t%s\t%.2f\n"
		fmt.Printf(csvFmt, typ, domain, subdomain, descr, token, labor)
		totTime += taskTime
		totExpense += expense
	}

	fmt.Fprintf(os.Stderr, "\nGenerated CSV between %s and %s\n",
		start.Format(dtFormat), end.Format(dtFormat))
	fmt.Fprintf(os.Stderr, "Total computed time: %s (%.2f hours)\n", totTime,
		totTime.Hours())
	fmt.Fprintf(os.Stderr, "Total Expense: $ %.2f\n", totExpense)
}

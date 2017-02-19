package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dustin/octoprint"
)

var (
	showStateFlags = flag.NewFlagSet("state", flag.ExitOnError)

	stateHist = showStateFlags.Int("history", 0, "how much history to grab")
	stateFmt  = showStateFlags.String("format", "plain", "plain | csv")
)

func maybePrintTemp(prefix string, t *octoprint.PrinterTempState) {
	if t == nil {
		return
	}
	fmt.Printf("%v: %v (Target %v)\n", prefix, t.Actual, t.Target)
}

func showStatePlain(st *octoprint.PrinterState) {
	fmt.Printf("State: %v\n", st.State)
	maybePrintTemp("\tBed", st.Temperature.Bed)
	maybePrintTemp("\tTool0", st.Temperature.Tool0)
	maybePrintTemp("\tTool1", st.Temperature.Tool1)
	for _, e := range st.Temperature.History {
		fmt.Printf("\t\t%v (%v)\n", e.Time(), time.Since(e.Time()))
		maybePrintTemp("\t\t\tBed", e.Bed)
		maybePrintTemp("\t\t\tTool0", e.Tool0)
		maybePrintTemp("\t\t\tTool1", e.Tool1)
	}
}

func appendTime(a []string, t *octoprint.PrinterTempState) []string {
	if t == nil {
		return append(a, "0", "0")
	}

	return append(a, fmt.Sprint(t.Target), fmt.Sprint(t.Actual))
}

func showStateCSV(st *octoprint.PrinterState) {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	w.Write([]string{"ts", "bedtarget", "bedactual", "tool0target", "tool0actual", "tool1target", "tool1actual"})

	for _, e := range st.Temperature.History {
		vals := []string{e.Time().Format(time.RFC3339Nano)}
		vals = appendTime(vals, e.Bed)
		vals = appendTime(vals, e.Tool0)
		vals = appendTime(vals, e.Tool1)
		w.Write(vals)
	}
}

func showState(ctx context.Context, c *octoprint.Client, args []string) {
	st, err := c.PrinterState(ctx, *stateHist)
	if err != nil {
		log.Fatalf("Error getting state: %v", err)
	}

	if *stateFmt == "csv" {
		showStateCSV(st)
		return
	}

	showStatePlain(st)
}

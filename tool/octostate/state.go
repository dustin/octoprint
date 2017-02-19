package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/dustin/octoprint"
)

var (
	showStateFlags = flag.NewFlagSet("state", flag.ExitOnError)

	stateHist = showStateFlags.Int("history", 0, "how much history to grab")
)

func maybePrintTemp(prefix string, t *octoprint.PrinterTempState) {
	if t == nil {
		return
	}
	fmt.Printf("%v: %v (Target %v)\n", prefix, t.Actual, t.Target)
}

func showState(ctx context.Context, c *octoprint.Client, args []string) {
	st, err := c.PrinterState(ctx, *stateHist)
	if err != nil {
		log.Fatalf("Error getting state: %v", err)
	}

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

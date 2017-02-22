package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/dustin/octoprint"
)

var (
	showStateFlags = flag.NewFlagSet("state", flag.ExitOnError)

	stateHist = showStateFlags.Int("history", 0, "how much history to grab")
	stateFmt  = showStateFlags.String("format", "plain", "plain | csv")
	stateTail = showStateFlags.Bool("tail", false, "continuously watch state")
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
	var pst *octoprint.PrinterState
	var jst *octoprint.JobState

	g := errgroup.Group{}

	g.Go(func() error {
		var err error
		pst, err = c.PrinterState(ctx, *stateHist)
		return err
	})

	g.Go(func() error {
		var err error
		jst, err = c.JobState(ctx)
		return err
	})

	if err := g.Wait(); err != nil {
		log.Fatalf("Error getting initial state: %v", err)
	}

	if *stateFmt == "csv" {
		showStateCSV(pst)
		return
	}

	showStatePlain(pst)

	if *stateTail {
		for range time.Tick(5 * time.Second) {
			oldn := len(pst.Temperature.History)

			g := errgroup.Group{}

			g.Go(func() error {
				return c.UpdatePrinterState(ctx, pst)
			})

			g.Go(func() error {
				var err error
				jst, err = c.JobState(ctx)
				return err
			})

			if err := g.Wait(); err != nil {
				log.Printf("Error updating: %v", err)
				continue
			}

			for i := oldn; i < len(pst.Temperature.History); i++ {
				e := pst.Temperature.History[i]
				fmt.Printf("\t\t%v (%v)\n", e.Time(), time.Since(e.Time()))
				maybePrintTemp("\t\t\tBed", e.Bed)
				maybePrintTemp("\t\t\tTool0", e.Tool0)
				maybePrintTemp("\t\t\tTool1", e.Tool1)
			}
			fmt.Printf("State: %v\n", pst.State)
			if jst.State == "Printing" {
				fmt.Printf("Job State: %v %.0f%% %v done, %v to go\n", jst.State, jst.Progress.Completion,
					time.Second*time.Duration(jst.Progress.PrintTime), time.Second*time.Duration(jst.Progress.PrintTimeLeft))
			}
		}
	}
}

package octoprint

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"sort"
	"time"
)

type PrinterTempState struct {
	Actual float64
	Offset int
	Target float64
}

func (p PrinterTempState) String() string {
	return fmt.Sprintf("{%.2f, target: %.1f}", p.Actual, p.Target)
}

type HistoricalPrinterTempEntry struct {
	Bed       *PrinterTempState
	TimeStamp int `json:"time"`
	Tool0     *PrinterTempState
	Tool1     *PrinterTempState
}

func (h HistoricalPrinterTempEntry) Time() time.Time {
	return time.Unix(int64(h.TimeStamp), 0)
}

func (h HistoricalPrinterTempEntry) String() string {
	parts := []string{"bed=%v"}
	args := []interface{}{h.Time(), h.Bed}
	if h.Tool0 != nil {
		parts = append(parts, "tool0=%v")
		args = append(args, h.Tool0)
	}
	if h.Tool1 != nil {
		parts = append(parts, "tool1=%v")
		args = append(args, h.Tool1)
	}

	return fmt.Sprintf("{@%v: "+strings.Join(parts, ", ")+"}", args...)
}

type PrinterStateState struct {
	Flags struct {
		ClosedOrError bool
		Error         bool
		Operational   bool
		Paused        bool
		Printing      bool
		Ready         bool
		SDReady       bool
	}
	Text string
}

func (p PrinterStateState) String() string {
	var flagStr []string
	if p.Flags.ClosedOrError {
		flagStr = append(flagStr, "ClosedOrError")
	}
	if p.Flags.Error {
		flagStr = append(flagStr, "Error")
	}
	if p.Flags.Operational {
		flagStr = append(flagStr, "Operational")
	}
	if p.Flags.Paused {
		flagStr = append(flagStr, "Paused")
	}
	if p.Flags.Printing {
		flagStr = append(flagStr, "Printing")
	}
	if p.Flags.Ready {
		flagStr = append(flagStr, "Ready")
	}
	if p.Flags.SDReady {
		flagStr = append(flagStr, "SDReady")
	}
	return fmt.Sprintf("%v (flags: %v)", p.Text, strings.Join(flagStr, "|"))
}

type PrinterState struct {
	SD struct {
		Ready bool
	}
	State       PrinterStateState
	Temperature struct {
		Bed     *PrinterTempState
		History []HistoricalPrinterTempEntry
		Tool0   *PrinterTempState
		Tool1   *PrinterTempState `json:",omitempty"`
	}
}

// Merge the history in this PrinterState with the given history,
// keeping all distinct timestamps with ascended ordering.
func (p *PrinterState) MergeHistory(h []HistoricalPrinterTempEntry) {
	seen := map[int]bool{}

	var a []HistoricalPrinterTempEntry
	for _, es := range [][]HistoricalPrinterTempEntry{p.Temperature.History, h} {
		for _, e := range es {
			if seen[e.TimeStamp] {
				continue
			}
			seen[e.TimeStamp] = true
			a = append(a, e)
		}
	}

	sort.Slice(a, func(i, j int) bool { return a[i].TimeStamp < a[j].TimeStamp })
	p.Temperature.History = a
}

// PrinterState fetches the current state of the printer.
func (c *Client) PrinterState(ctx context.Context, history int) (*PrinterState, error) {
	args := ""
	if history > 0 {
		args = "history=true&limit=" + strconv.Itoa(history)
	}
	st := &PrinterState{}
	if err := c.fetchJSON(ctx, "/api/printer", args, st); err != nil {
		return nil, err
	}

	return st, nil
}

// UpdateState updates all of the values of given state, including merging in
// new history with existing history.
func (c *Client) UpdatePrinterState(ctx context.Context, st *PrinterState) error {
	nHist := 60
	if len(st.Temperature.History) > 0 {
		latestTs := st.Temperature.History[len(st.Temperature.History)-1].Time()
		nHist = int((time.Since(latestTs) / 5) + 1)
	}
	nst, err := c.PrinterState(ctx, nHist)
	if err != nil {
		return err
	}
	nst.MergeHistory(st.Temperature.History)
	st.SD = nst.SD
	st.State = nst.State
	st.Temperature = nst.Temperature

	return nil
}

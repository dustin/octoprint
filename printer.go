package octoprint

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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

package octoprint

import (
	"fmt"
	"strings"

	"time"
)

type PrinterTempState struct {
	Actual float64
	Offset int
	Target int
}

func (p PrinterTempState) String() string {
	return fmt.Sprintf("{%.2f, target: %v}", p.Actual, p.Target)
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

type PrinterState struct {
	SD struct {
		Ready bool
	}
	State struct {
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
	Temperature struct {
		Bed     *PrinterTempState
		History []HistoricalPrinterTempEntry
		Tool0   *PrinterTempState
		Tool1   *PrinterTempState `json:",omitempty"`
	}
}

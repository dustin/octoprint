package octoprint

import (
	"encoding/json"
	"testing"
)

var printerStateSample = `{
  "sd": {
    "ready": false
  },
  "state": {
    "flags": {
      "closedOrError": false,
      "error": false,
      "operational": true,
      "paused": false,
      "printing": false,
      "ready": true,
      "sdReady": false
    },
    "text": "Operational"
  },
  "temperature": {
    "bed": {
      "actual": 18.33,
      "offset": 0,
      "target": 0
    },
    "history": [
      {
        "bed": {
          "actual": 18.15,
          "target": 0
        },
        "time": 1487476988,
        "tool0": {
          "actual": 17.59,
          "target": 0
        }
      },
      {
        "bed": {
          "actual": 18.15,
          "target": 0
        },
        "time": 1487476993,
        "tool0": {
          "actual": 17.59,
          "target": 0
        }
      },
      {
        "bed": {
          "actual": 18.15,
          "target": 0
        },
        "time": 1487476998,
        "tool0": {
          "actual": 17.59,
          "target": 0
        }
      },
      {
        "bed": {
          "actual": 18.15,
          "target": 0
        },
        "time": 1487477003,
        "tool0": {
          "actual": 17.59,
          "target": 0
        }
      },
      {
        "bed": {
          "actual": 18.33,
          "target": 0
        },
        "time": 1487477008,
        "tool0": {
          "actual": 17.59,
          "target": 0
        }
      },
      {
        "bed": {
          "actual": 18.33,
          "target": 0
        },
        "time": 1487477013,
        "tool0": {
          "actual": 17.59,
          "target": 0
        }
      },
      {
        "bed": {
          "actual": 18.33,
          "target": 0
        },
        "time": 1487477018,
        "tool0": {
          "actual": 17.59,
          "target": 0
        }
      },
      {
        "bed": {
          "actual": 18.33,
          "target": 0
        },
        "time": 1487477023,
        "tool0": {
          "actual": 17.59,
          "target": 0
        }
      },
      {
        "bed": {
          "actual": 18.33,
          "target": 0
        },
        "time": 1487477028,
        "tool0": {
          "actual": 17.59,
          "target": 0
        }
      },
      {
        "bed": {
          "actual": 18.33,
          "target": 0
        },
        "time": 1487477033,
        "tool0": {
          "actual": 17.59,
          "target": 0
        }
      }
    ],
    "tool0": {
      "actual": 17.59,
      "offset": 0,
      "target": 0
    }
  }
}`

func TestPrinterStateParsing(t *testing.T) {
	s := &PrinterState{}
	if err := json.Unmarshal([]byte(printerStateSample), s); err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", s)
}

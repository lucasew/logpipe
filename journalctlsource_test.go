package logpipe

import (
	"testing"
	"time"

	"github.com/lucasew/gocfg"
)

func TestJournalctlSource(t *testing.T) {
    sink := NewLogPipeTestingSink()
    source, err := NewJournalctlSource(gocfg.NewMapSectionProvider())
    if err != nil {
        t.Fatal(err)
    }
    lp := NewLogPipe()
    lp.RegisterSource("journalctl", source)
    lp.RegisterSink("echo", sink)
    ticker := time.Tick(2*time.Second)
    testTimeEnded := false
    for !testTimeEnded {
        select {
        case <-ticker:
            testTimeEnded = true
        default:
            lp.Tick()
        }
    }
    if sink.GetTextSize() == 0 {
        t.Fail()
    }
}

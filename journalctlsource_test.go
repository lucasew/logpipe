package logpipe

import (
	"testing"
	"github.com/lucasew/gocfg"
	"time"
)

func TestJournalctlSource(t *testing.T) {
    sink := NewLogPipeTestingSink()
    source, _ := NewJournalctlSource(gocfg.NewMapSectionProvider())
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

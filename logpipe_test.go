package logpipe

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

type LogPipeTestingSource chan string

func NewLogPipeTestingSource() LogPipeTestingSource {
    return LogPipeTestingSource(make(chan string, 10))
}

func (ts LogPipeTestingSource) GetSource() <-chan string {
    return ts
}

func (ts LogPipeTestingSource) Emit(message string) {
    ts<-message
}


type LogPipeTestingSink struct {
    sync.Mutex
    texts []string
    ch chan string
    started bool
}

func (ts *LogPipeTestingSink) GetTextIdx(idx int) string {
    ts.Lock()
    defer ts.Unlock()
    return ts.texts[idx]
}

func (ts *LogPipeTestingSink) GetTextSize() int {
    ts.Lock()
    defer ts.Unlock()
    return len(ts.texts)
}

func (ts *LogPipeTestingSink) GetSink() chan<- string {
    if !ts.started {
        log.Printf("started sink")
        go func(ts *LogPipeTestingSink) {
            for text := range ts.ch {
                ts.Lock()
                log.Printf("sink %s", text)
                ts.texts = append(ts.texts, text)
                ts.Unlock()
            }
        }(ts)
        ts.started = true
    }
    return ts.ch
}

func NewLogPipeTestingSink() *LogPipeTestingSink {
    return &LogPipeTestingSink{
        texts: make([]string, 0, 10),
        ch: make(chan string, 10),
        started: false,
    }
}


func TestBroadcast(t *testing.T) {
    var source1 Source
    var sink1 Sink
    var sink2 Sink
    nop(source1, sink1, sink2) // just to avoid the unused error
    tso1 := NewLogPipeTestingSource()
    tsi1 := NewLogPipeTestingSink()
    tsi2 := NewLogPipeTestingSink()
    source1 = tso1 // Static analysis
    sink1 = tsi1
    sink2 = tsi2
    lp := NewLogPipe()
    lp.RegisterSource("sample", tso1)
    lp.RegisterSink("out1", tsi1)
    lp.RegisterSink("out2", tsi2)
    for i := 0; i < 10; i++ {
        tso1.Emit(fmt.Sprintf("%d", i))
    }
    for i := 0; i < 10; {
        if lp.Tick() {
            i++
        }
    }
    time.Sleep(time.Second) // delay to other goroutines settle
    for i := 0; i < 10; i++ {
        expected := fmt.Sprintf("%d", i)
        got_a := tsi1.GetTextIdx(i)
        got_b := tsi2.GetTextIdx(i)
        if expected != got_a {
            t.Errorf("unexpected issue at %dnth message on first sink: expected: '%s', got '%s'", i + 1, expected, got_a)
        }
        if expected != got_b {
            t.Errorf("unexpected issue at %dnth on second sink message: expected: '%s', got '%s'", i + 1, expected, got_b)
        }
    }
}

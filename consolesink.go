package logpipe

import (
	"log"

	"github.com/lucasew/gocfg"
)

type ConsoleSink struct {
    ch chan string
    started bool
}

// NewConsoleSink provides a simple sink that outputs messages to standard
// logging output via an internal channel loop. It does not enforce
// any specific rate limits.
func NewConsoleSink(cfg gocfg.Section) (Sink, error) {
    return &ConsoleSink{
        ch: make(chan string, 1),
        started: false,
    }, nil
}

func (t *ConsoleSink) GetSink() chan <- string {
    if !t.started {
        go func () {
            for {
                select {
                case msg := <-t.ch:
                    log.Printf("terminal: %s", msg)
                }
            }
        }()
        t.started = true
    }
    return t.ch
}

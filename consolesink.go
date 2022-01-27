package logpipe

import (
	"log"

	"github.com/lucasew/gocfg"
)

type ConsoleSink struct {
    ch chan string
    started bool
}

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

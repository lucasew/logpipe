package logpipe

import (
	"sync"
)

func nop(_ ...interface {}) {
}

type Source interface {
    GetSource() <-chan string // Sources can only provide lines
}

type Sink interface {
    GetSink() chan <- string // Sinks can only consume lines
}


type LogPipe struct {
    sync.RWMutex
    sources map[string]Source
    sinks map[string]Sink
}

func NewLogPipe() *LogPipe {
    return &LogPipe{
        sources: map[string]Source{},
        sinks: map[string]Sink{},
    }
}

func (l *LogPipe) RegisterSource(name string, source Source) {
    l.Lock()
    defer l.Unlock()
    l.sources[name] = source
}

func (l *LogPipe) RegisterSink(name string, sink Sink) {
    l.Lock()
    defer l.Unlock()
    l.sinks[name] = sink
}

func (l *LogPipe) broadcast(source string, message string) {
    l.RLock()
    defer l.RUnlock()
    for _, sink := range l.sinks {
        sink.GetSink()<-message
    }
}

func (l *LogPipe) Tick() (hasMessage bool) {
    hasMessage = false
    l.RLock()
    defer l.RUnlock()
    for k, v := range l.sources {
        select {
            case msg := <-v.GetSource():
                l.broadcast(k, msg)
                hasMessage = true
            default:
                continue
        }
    }
    return
}

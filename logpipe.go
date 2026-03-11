package logpipe

import (
	"sync"
)

func nop(_ ...interface {}) {
}

// Source represents a log provider in the pipeline.
// Implementations should continuously emit log lines (e.g., from journalctl or a file)
// to the channel returned by GetSource.
type Source interface {
    GetSource() <-chan string // Sources can only provide lines
}

// Sink represents a log consumer in the pipeline.
// Implementations should read from the channel returned by GetSink
// and forward the logs to external services (e.g., Telegram, Discord).
type Sink interface {
    GetSink() chan <- string // Sinks can only consume lines
}

// LogPipe is the central hub connecting Sources and Sinks.
// It manages the registration of components and handles the broadcasting
// of messages from all registered sources to all registered sinks.
// It uses a RWMutex to ensure thread-safe concurrent modifications to the registry.
type LogPipe struct {
    sync.RWMutex
    sources map[string]Source
    sinks map[string]Sink
}

// NewLogPipe initializes and returns a new LogPipe instance
// with empty source and sink registries.
func NewLogPipe() *LogPipe {
    return &LogPipe{
        sources: map[string]Source{},
        sinks: map[string]Sink{},
    }
}

// RegisterSource adds a new named Source to the pipeline.
// It acquires a write lock to ensure thread safety while modifying the registry.
func (l *LogPipe) RegisterSource(name string, source Source) {
    l.Lock()
    defer l.Unlock()
    l.sources[name] = source
}

// RegisterSink adds a new named Sink to the pipeline.
// It acquires a write lock to ensure thread safety while modifying the registry.
func (l *LogPipe) RegisterSink(name string, sink Sink) {
    l.Lock()
    defer l.Unlock()
    l.sinks[name] = sink
}

// broadcast sends the given message to all registered sinks.
// It acquires a read lock to allow concurrent broadcasting while preventing
// registry modifications during the iteration.
// Note: This operation can block if a sink's channel is full.
func (l *LogPipe) broadcast(source string, message string) {
    l.RLock()
    defer l.RUnlock()
    for _, sink := range l.sinks {
        sink.GetSink()<-message
        // log.Printf("broadcasting to '%s': %s", k, message)
    }
}

// Tick performs a non-blocking read operation across all registered sources.
// If a message is available from any source, it broadcasts it to all sinks
// and returns true. If no messages are ready, it returns false immediately.
// This allows the main loop to periodically poll for updates without getting stuck.
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

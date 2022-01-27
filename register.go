package logpipe

import "github.com/lucasew/gocfg"

var REGISTERED_SOURCES = map[string](func(gocfg.Section) (Source, error)){}

var REGISTERED_SINKS = map[string](func(gocfg.Section) (Sink, error)){}

func init() {
    REGISTERED_SOURCES["journalctl"] = NewJournalctlSource
    REGISTERED_SINKS["telegram"] = NewTelegramSink
    REGISTERED_SINKS["discord"] = NewDiscordSink
    REGISTERED_SINKS["console"] = NewConsoleSink
}

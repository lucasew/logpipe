package logpipe

import "github.com/lucasew/gocfg"

// REGISTERED_SOURCES contains a mapping of source names to their respective factory functions.
// These functions accept a configuration section and return a new Source instance or an error.
var REGISTERED_SOURCES = map[string](func(gocfg.Section) (Source, error)){}

// REGISTERED_SINKS contains a mapping of sink names to their respective factory functions.
// These functions accept a configuration section and return a new Sink instance or an error.
var REGISTERED_SINKS = map[string](func(gocfg.Section) (Sink, error)){}

// init populates the global registries for built-in sources and sinks
// with their corresponding factory functions (e.g., journalctl, telegram, discord).
// This enables dynamic instantiation based on configuration files upon package load.
func init() {
    REGISTERED_SOURCES["journalctl"] = NewJournalctlSource
    REGISTERED_SINKS["telegram"] = NewTelegramSink
    REGISTERED_SINKS["discord"] = NewDiscordSink
    REGISTERED_SINKS["console"] = NewConsoleSink
}

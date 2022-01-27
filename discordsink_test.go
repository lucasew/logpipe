package logpipe

import (
	"os"
	"testing"

	"github.com/lucasew/gocfg"
)

func TestDiscordSink(t *testing.T) {
    cfg := gocfg.NewConfig()
    webhook := os.Getenv("DISCORD_WEBHOOK")
    if webhook == "" {
        t.Skipf("DISCORD_WEBHOOK not defined, skipping discord tests")
    }
    cfg.RawSet("discord", "webhook", webhook)
    sink, err := NewDiscordSink(cfg["discord"])
    if err != nil {
        t.Error(err)
    }
    source := NewLogPipeTestingSource()
    lp := NewLogPipe()
    lp.RegisterSource("std", source)
    lp.RegisterSink("discord", sink)
    source.Emit("test successful")
    lp.Tick()
}

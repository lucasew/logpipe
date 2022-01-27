package logpipe

import (
	"os"
	"testing"
	"time"

	"github.com/lucasew/gocfg"
)

func TestTelegramSink(t *testing.T) {
    cfg := gocfg.NewConfig()
    chat_id := os.Getenv("TELEGRAM_CHAT_ID")
    if chat_id == "" {
        t.Skipf("TELEGRAM_CHAT_ID not defined, skipping telegram tests")
    }
    cfg.RawSet("telegram", "chat_id", chat_id)

    token := os.Getenv("TELEGRAM_TOKEN")
    if token == "" {
        t.Skipf("TELEGRAM_TOKEN not defined, skipping telegram tests")
    }
    cfg.RawSet("telegram", "token", token)
    sink, err := NewTelegramSink(cfg["telegram"])
    if err != nil {
        t.Error(err)
    }
    source := NewLogPipeTestingSource()
    lp := NewLogPipe()
    lp.RegisterSource("std", source)
    lp.RegisterSink("telegram", sink)
    source.Emit("test successful")
    lp.Tick()
    time.Sleep(time.Second)
}

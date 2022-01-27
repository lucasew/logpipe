package logpipe

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/lucasew/gocfg"
)

type discordSink struct {
    cfg gocfg.Section
    webhook *url.URL
    ch chan string
    started bool
}

func NewDiscordSink(cfg gocfg.Section) (Sink, error) {
    u, err := url.Parse(cfg["webhook"])
    if err != nil {
        return nil, err
    }
    return &discordSink{
        webhook: u,
        ch: make(chan string, 20),
        cfg: cfg,
    }, nil
}

func (d *discordSink) GetSink() chan<- string {
    if !d.started {
        go func() {
            timer := time.NewTimer(100*time.Millisecond)
            for {
                select {
                    case msg := <- d.ch:
                        params := url.Values{}
                        params.Add("content", msg)
                        encoded := params.Encode()
                        req, err := http.NewRequest("POST", d.webhook.String(), bytes.NewBufferString(encoded))
                        if err != nil {
                            panic(err)
                        }
                        req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
                        try := 0
                        retry:
                        try++
                        if try >= 3 {
                            continue
                        }
                        <-timer.C
                        res, err := http.DefaultClient.Do(req)
                        io.Copy(os.Stdout, res.Body)
                        if err != nil {
                            log.Printf("retry: %s", err.Error())
                            goto retry
                        }
                    case <-timer.C:
                        continue
                }
            }
        }()
        d.started = true
    }
    return d.ch
}

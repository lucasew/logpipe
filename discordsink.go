package logpipe

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/lucasew/gocfg"
)

type discordSink struct {
    cfg gocfg.Section
    webhook *url.URL
    ch chan string
    started bool
}

var (
    ErrDiscordSinkNoWebhookProvided = errors.New("no webhook was provided")
)

func NewDiscordSink(cfg gocfg.Section) (Sink, error) {
    webhook, ok := cfg["webhook"]
    if !ok {
        return nil, ErrDiscordSinkNoWebhookProvided
    }
    u, err := url.Parse(webhook)
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
            timer := time.Tick(200*time.Millisecond)
            for msg := range d.ch {
                // log.Printf("discord %s", msg)
                params := url.Values{}
                params.Add("content", msg)
                encoded := params.Encode()
                req, err := http.NewRequest("POST", d.webhook.String(), bytes.NewBufferString(encoded))
                if err != nil {
                    panic(err)
                }
                req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
                <-timer
                _, err = http.DefaultClient.Do(req)
                // io.Copy(os.Stdout, res.Body)
            }
        }()
        d.started = true
    }
    return d.ch
}

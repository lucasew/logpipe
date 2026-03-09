package logpipe

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/lucasew/gocfg"
)

type discordSink struct {
	cfg     gocfg.SectionProvider
	webhook *url.URL
	ch      chan string
	started bool
}

var (
	ErrDiscordSinkNoWebhookProvided = errors.New("no webhook was provided")
)

func NewDiscordSink(cfg gocfg.SectionProvider) (Sink, error) {
	webhook := cfg.RawGet("webhook")
	ok := cfg.RawHasKey("webhook")
	if !ok {
		return nil, ErrDiscordSinkNoWebhookProvided
	}
	u, err := url.Parse(webhook)
	if err != nil {
		return nil, err
	}
	return &discordSink{
		webhook: u,
		ch:      make(chan string, 20),
		cfg:     cfg,
	}, nil
}

func (d *discordSink) GetSink() chan<- string {
	if !d.started {
		go func() {
			ticker := time.NewTicker(200 * time.Millisecond)
			defer ticker.Stop()
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
				<-ticker.C
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Printf("error calling discord webhook: %s", err)
				} else if res != nil {
					if closeErr := res.Body.Close(); closeErr != nil {
						log.Printf("error closing response body: %s", closeErr)
					}
				}
				// io.Copy(os.Stdout, res.Body)
			}
		}()
		d.started = true
	}
	return d.ch
}

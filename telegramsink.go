package logpipe

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/lucasew/gocfg"
)

type telegramSink struct {
	cfg     gocfg.SectionProvider
	ch      chan string
	started bool
}

func NewTelegramSink(cfg gocfg.SectionProvider) (Sink, error) {
	return &telegramSink{
		cfg:     cfg,
		ch:      make(chan string, 10),
		started: false,
	}, nil
}

func (d *telegramSink) GetSink() chan<- string {
	if !d.started {
		go func() {
			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()
			for msg := range d.ch {
				// log.Printf("telegram: %s", msg)
				params := url.Values{}
				params.Add("text", msg)
				params.Add("chat_id", d.cfg.RawGet("chat_id"))
				// TODO: add more parameters as defined in https://core.telegram.org/bots/api#sendmessage
				encoded := params.Encode()
				url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?%s", d.cfg.RawGet("token"), encoded)
				req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte{}))
				if err != nil {
					panic(err)
				}
				<-ticker.C
				res, err := http.DefaultClient.Do(req)
				if res != nil {
					if err := res.Body.Close(); err != nil {
						log.Printf("error closing response body: %s", err)
					}
				}
				// io.Copy(os.Stdout, res.Body)
				if err != nil {
					log.Printf("retry: %s", err.Error())
				}
			}
		}()
		d.started = true
	}
	return d.ch
}

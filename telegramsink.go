package logpipe

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/lucasew/gocfg"
)

type telegramSink struct {
    cfg gocfg.Section
    ch chan string
    started bool
}

func NewTelegramSink(cfg gocfg.Section) (Sink, error) {
    return &telegramSink{
        cfg: cfg,
        ch: make(chan string, 10),
        started: false,
    }, nil
}

func (d *telegramSink) GetSink() chan<- string {
    if !d.started {
        go func() {
            timer := time.Tick(100*time.Millisecond)
            for msg := range d.ch {
                // log.Printf("telegram: %s", msg)
                params := url.Values{}
                params.Add("text", msg)
                params.Add("chat_id", d.cfg["chat_id"])
                // TODO: add more parameters as defined in https://core.telegram.org/bots/api#sendmessage
                encoded := params.Encode()
                url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?%s", d.cfg["token"], encoded)
                req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte{}))
                if err != nil {
                    ReportError(err)
                    continue
                }
                <-timer
                res, err := http.DefaultClient.Do(req)
                if err != nil {
                    ReportError(err)
                    continue
                }
                res.Body.Close()
                // io.Copy(os.Stdout, res.Body)
            }
        }()
        d.started = true
    }
    return d.ch
}

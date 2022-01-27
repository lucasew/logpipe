package logpipe

import (
	"bytes"
	"fmt"
	// "io"
	"log"
	"net/http"
	"net/url"
	// "os"
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
            timer := time.NewTimer(100*time.Millisecond)
            for {
                select {
                    case msg := <- d.ch:
                        params := url.Values{}
                        params.Add("text", msg)
                        params.Add("chat_id", d.cfg["chat_id"])
                        // TODO: add more parameters as defined in https://core.telegram.org/bots/api#sendmessage
                        encoded := params.Encode()
                        url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?%s", d.cfg["token"], encoded)
                        req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte{}))
                        if err != nil {
                            panic(err)
                        }
                        try := 0
                        retry:
                        try++
                        if try >= 3 {
                            continue
                        }
                        <-timer.C
                        res, err := http.DefaultClient.Do(req)
                        nop(res)
                        // io.Copy(os.Stdout, res.Body)
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

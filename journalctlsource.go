package logpipe

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
	"text/template"

	"github.com/lucasew/gocfg"
)

type journalctlSource struct {
    ch chan string
    cfg gocfg.Section
    started bool
    template *template.Template
}

func NewJournalctlSource(cfg gocfg.Section) (Source, error) {
    tmplStr, ok := cfg["format"]
    if !ok {
        tmplStr = "#{{._HOSTNAME}} {{.__REALTIME_TIMESTAMP}} ({{._SYSTEMD_CGROUP}}): {{.MESSAGE}}"
    }
    tmpl, err := template.New("msg").Parse(tmplStr)
    if err != nil {
        return nil, err
    }
    return &journalctlSource{
        ch: make(chan string, 10),
        cfg: cfg,
        template: tmpl,
    }, nil
}

func (j *journalctlSource) GetSource() <-chan string {
    if !j.started {
        go func () {
            cmd := exec.Command("journalctl", "--no-pager", "--output=json", "-f", "--utc")
            stdout, err := cmd.StdoutPipe()
            if err != nil {
                panic(err)
            }
            scanner := bufio.NewScanner(stdout)
            err = cmd.Start()
            if err != nil {
                panic(err)
            }
            for scanner.Scan() {
                val := map[string]string{}
                if scanner.Err() != nil {
                    panic(scanner.Err())
                }
                line := scanner.Text()
                err = json.Unmarshal([]byte(line), &val)
                if err != nil {
                    panic(err)
                }
                buf := bytes.NewBuffer([]byte{})
                err := j.template.Execute(buf, val)
                if err != nil {
                    log.Printf("error(source/journalctl): %s", err.Error())
                }
                j.ch <- buf.String()
            }
        }()
        j.started = true
    }
    return j.ch
}

package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/lucasew/gocfg"
	"github.com/lucasew/logpipe"
)

var CONFIG_FILE string
var cfg = gocfg.NewConfig()
var lp = logpipe.NewLogPipe()

func init () {
    flag.StringVar(&CONFIG_FILE, "c", "/etc/logpipe.ini", "Configuration file to use")
    flag.Parse()
    f, err := os.Open(CONFIG_FILE)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    err = cfg.InjestReader(f)
    if err != nil {
        panic(err)
    }
    for k, v := range cfg {
        if strings.HasPrefix(k, "source.") {
            parts := strings.Split(k, ".")
            if len(parts) != 2 {
                log.Fatalf("invalid source definition in section '%s'", k)
            }
            name := parts[1]
            sourceType, ok := v["type"]
            if !ok {
                log.Fatalf("no source type was provided in section '%s'", k)
            }
            newSource, ok := logpipe.REGISTERED_SOURCES[sourceType]
            if !ok {
                log.Fatalf("undefined source type '%s' in section '%s'", name, k)
            }
            source, err := newSource(v)
            if err != nil {
                log.Fatalf("error initializing source '%s' in section '%s': %s", name, k, err.Error())
            }
            lp.RegisterSource(name, source)
            continue
        }
        if strings.HasPrefix(k, "sink.") {
            parts := strings.Split(k, ".")
            if len(parts) != 2 {
                log.Fatalf("invalid sink definition in section '%s'", k)
            }
            name := parts[1]
            sourceType, ok := v["type"]
            if !ok {
                log.Fatalf("no source type was provided in section '%s'", k)
            }
            newSink, ok := logpipe.REGISTERED_SINKS[sourceType]
            if !ok {
                log.Fatalf("undefined sink type '%s' in section '%s'", name, k)
            }
            sink, err := newSink(v)
            if err != nil {
                log.Fatalf("error initializing '%s' in sink '%s': %s", name, k, err.Error())
            }
            lp.RegisterSink(name, sink)
            continue
        }
        if k == "env" {
            for vark, varv := range v {
                err := os.Setenv(vark, varv)
                if err != nil {
                    panic(err)
                }
            }
        }
    }
}

func main() {
    for {
        lp.Tick()
    }
}

package main

import (
	"errors"
	"flag"
	"log/slog"
	"os"
	"strings"

	"github.com/itsokto/cloud-geoip/writer"
)

var targets = []struct {
	Name string
	ASNs []string
}{
	{Name: "akamai", ASNs: []string{"AS-AKAMAI"}},
	{Name: "alibaba", ASNs: []string{"AS37963", "AS45102", "AS24429", "AS134963", "AS203513"}},
	{Name: "tencent", ASNs: []string{"AS45090", "AS132203", "AS133478", "AS137876"}},
	{Name: "cognosphere", ASNs: []string{"AS203923"}},
	{Name: "ucloud", ASNs: []string{"AS135377", "AS139327"}},
}

func main() {
	outputDir := flag.String("output", "output", "output directory")
	noV4 := flag.Bool("no-v4", false, "skip IPv4")
	noV6 := flag.Bool("no-v6", false, "skip IPv6")
	noAggregate := flag.Bool("no-aggregate", false, "disable prefix aggregation")
	verbose := flag.Bool("v", false, "verbose logging")
	flag.Parse()

	if *verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	irr := newIRR()
	ripe := newRIPEStat()

	var entries []writer.Entry
	for _, t := range targets {
		slog.Info("target", "name", t.Name, "asns", strings.Join(t.ASNs, " "))

		var prefixes []string
		asns, err := irr.resolveASNs(t.ASNs)
		if err == nil {
			prefixes, err = ripe.announcedPrefixesFor(asns)
		}
		if err == nil && len(prefixes) == 0 {
			err = errors.New("no announced prefixes")
		}
		if err == nil {
			prefixes, err = filterPrefixes(prefixes, *noV4, *noV6)
		}
		if err == nil && !*noAggregate {
			prefixes, err = aggregatePrefixes(prefixes)
		}
		if err != nil {
			slog.Error("target failed", "name", t.Name, "err", err)
			os.Exit(1)
		}

		entries = append(entries, writer.Entry{Name: t.Name, Prefixes: prefixes})
		slog.Info("collected", "name", t.Name, "prefixes", len(prefixes))
	}

	writers := []writer.Writer{
		&writer.PlainWriter{},
		&writer.SRSWriter{},
		&writer.DatWriter{Filename: "cloud-geoip.dat"},
	}

	for _, w := range writers {
		if err := w.Write(*outputDir, entries); err != nil {
			slog.Error("write failed", "err", err)
			os.Exit(1)
		}
	}
}

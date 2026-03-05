package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/itsokto/cloud-geoip/writer"
)

var targets = []struct {
	Name string
	ASNs []string
}{
	{"akamai", []string{"AS-AKAMAI"}},
	{"alibaba", []string{"AS37963", "AS45102", "AS24429"}},
	{"cognosphere", []string{"AS135377"}},
}

func main() {
	outputDir := flag.String("output", "output", "output directory")
	noV4 := flag.Bool("no-v4", false, "skip IPv4")
	noV6 := flag.Bool("no-v6", false, "skip IPv6")
	noAggregate := flag.Bool("no-aggregate", false, "disable prefix aggregation")
	sources := flag.String("S", "", "IRR sources (passed to bgpq4 -S)")
	host := flag.String("h", "", "IRR server (passed to bgpq4 -h)")
	flag.Parse()

	log.SetFlags(0)

	var extra []string
	if *sources != "" {
		extra = append(extra, "-S", *sources)
	}
	if *host != "" {
		extra = append(extra, "-h", *host)
	}
	if !*noAggregate {
		extra = append(extra, "-A")
	}

	var entries []writer.Entry
	for _, t := range targets {
		fmt.Fprintf(os.Stderr, "\n=== %s (%s) ===\n", t.Name, strings.Join(t.ASNs, " "))

		prefixes, err := queryPrefixes(t.ASNs, extra, *noV4, *noV6)
		if err != nil {
			log.Fatalf("%s: %v", t.Name, err)
		}

		entries = append(entries, writer.Entry{Name: t.Name, Prefixes: prefixes})
		fmt.Fprintf(os.Stderr, "  %d prefixes\n", len(prefixes))
	}

	writers := []writer.Writer{
		&writer.PlainWriter{},
		&writer.SRSWriter{},
		&writer.DatWriter{},
	}

	for _, w := range writers {
		if err := w.Write(*outputDir, entries); err != nil {
			log.Fatalf("write: %v", err)
		}
	}
}

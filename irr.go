package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	defaultIRRServer = "whois.radb.net:43"
	irrTimeout       = 60 * time.Second
)

type irrClient struct {
	server string
}

func newIRR() *irrClient {
	return &irrClient{server: defaultIRRServer}
}

func isASSet(name string) bool {
	return strings.HasPrefix(name, "AS-") || strings.Contains(name, ":")
}

func (c *irrClient) expandASSet(name string) ([]string, error) {
	conn, err := net.DialTimeout("tcp", c.server, irrTimeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(irrTimeout)); err != nil {
		return nil, err
	}

	if _, err := fmt.Fprintf(conn, "!!\n!i%s,1\n!q\n", name); err != nil {
		return nil, err
	}

	r := bufio.NewReader(conn)
	header, err := r.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	header = strings.TrimSpace(header)

	if !strings.HasPrefix(header, "A") {
		return nil, fmt.Errorf("%s: IRR returned %q", name, header)
	}
	n, err := strconv.Atoi(header[1:])
	if err != nil {
		return nil, fmt.Errorf("%s: bad IRR header %q", name, header)
	}

	buf := make([]byte, n)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}

	members := strings.Fields(string(buf))
	if len(members) == 0 {
		return nil, fmt.Errorf("%s: no members", name)
	}
	slog.Info("expanded as-set", name, "asns", len(members))
	return members, nil
}

func (c *irrClient) resolveASNs(names []string) ([]string, error) {
	var out []string
	seen := map[string]bool{}
	for _, name := range names {
		asns := []string{name}
		if isASSet(name) {
			var err error
			if asns, err = c.expandASSet(name); err != nil {
				return nil, err
			}
		}
		for _, asn := range asns {
			if !seen[asn] {
				seen[asn] = true
				out = append(out, asn)
			}
		}
	}
	return out, nil
}

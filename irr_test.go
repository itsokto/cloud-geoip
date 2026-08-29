package main

import (
	"bufio"
	"fmt"
	"net"
	"reflect"
	"strings"
	"testing"
)

func stubIRR(t *testing.T, replies map[string]string) *irrClient {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func() {
				defer conn.Close()
				s := bufio.NewScanner(conn)
				for s.Scan() {
					q := s.Text()
					if !strings.HasPrefix(q, "!i") {
						continue
					}
					name := strings.TrimSuffix(strings.TrimPrefix(q, "!i"), ",1")
					members, ok := replies[name]
					if !ok {
						fmt.Fprint(conn, "D\n")
						continue
					}
					fmt.Fprintf(conn, "A%d\n%s\nC\n", len(members)+1, members)
				}
			}()
		}
	}()

	c := newIRR()
	c.server = ln.Addr().String()
	return c
}

func TestIsASSet(t *testing.T) {
	for name, want := range map[string]bool{
		"AS-AKAMAI":          true,
		"AS1234:AS-CUSTOMER": true,
		"AS203923":           false,
		"AS45090":            false,
	} {
		if got := isASSet(name); got != want {
			t.Errorf("isASSet(%q) = %v, want %v", name, got, want)
		}
	}
}

func TestResolveASNs(t *testing.T) {
	c := stubIRR(t, map[string]string{"AS-TEST": "AS1 AS2 AS3"})

	got, err := c.resolveASNs([]string{"AS-TEST", "AS3", "AS9"})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"AS1", "AS2", "AS3", "AS9"}; !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	if _, err := c.resolveASNs([]string{"AS-MISSING"}); err == nil {
		t.Error("expected error for unknown as-set")
	}
}

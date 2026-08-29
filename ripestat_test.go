package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func stubRIPEStat(t *testing.T, prefixes map[string][]string) *ripeStat {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body []string
		for _, p := range prefixes[r.URL.Query().Get("resource")] {
			body = append(body, fmt.Sprintf(`{"prefix":%q}`, p))
		}
		fmt.Fprintf(w, `{"data":{"prefixes":[%s]}}`, strings.Join(body, ","))
	}))
	t.Cleanup(srv.Close)

	c := newRIPEStat()
	c.baseURL = srv.URL
	return c
}

func TestAnnouncedPrefixesFor(t *testing.T) {
	c := stubRIPEStat(t, map[string][]string{
		"AS1": {"1.2.3.0/24"},
		"AS2": {"2001:db8::/32"},
		"AS3": nil, // dormant as-set member
	})

	got, err := c.announcedPrefixesFor([]string{"AS1", "AS3", "AS2"})
	if err != nil {
		t.Fatal(err)
	}
	// Order follows the input, not whichever request finished first.
	if want := []string{"1.2.3.0/24", "2001:db8::/32"}; !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAnnouncedPrefixesError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := newRIPEStat()
	c.baseURL = srv.URL
	if _, err := c.announcedPrefixes("AS1"); err == nil {
		t.Error("expected error on 500")
	}
}

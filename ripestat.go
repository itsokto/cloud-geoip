package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/sync/errgroup"
)

const defaultRIPEStatURL = "https://stat.ripe.net/data/announced-prefixes/data.json"

var httpClient = &http.Client{Timeout: 60 * time.Second}

type ripeStat struct {
	baseURL     string
	concurrency int
}

func newRIPEStat() *ripeStat {
	return &ripeStat{baseURL: defaultRIPEStatURL, concurrency: 3}
}

func (c *ripeStat) announcedPrefixes(asn string) ([]string, error) {
	endpoint := fmt.Sprintf("%s?sourceapp=cloud-geoip&resource=%s", c.baseURL, url.QueryEscape(asn))

	resp, err := httpClient.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", asn, resp.Status)
	}

	var result struct {
		Data struct {
			Prefixes []struct {
				Prefix string `json:"prefix"`
			} `json:"prefixes"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("%s: %w", asn, err)
	}

	prefixes := make([]string, len(result.Data.Prefixes))
	for i, p := range result.Data.Prefixes {
		prefixes[i] = p.Prefix
	}
	return prefixes, nil
}

func (c *ripeStat) announcedPrefixesFor(asns []string) ([]string, error) {
	results := make([][]string, len(asns))

	var g errgroup.Group
	g.SetLimit(c.concurrency)
	for i, asn := range asns {
		g.Go(func() error {
			prefixes, err := c.announcedPrefixes(asn)
			results[i] = prefixes
			return err
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	var all []string
	for _, r := range results {
		all = append(all, r...)
	}
	return all, nil
}

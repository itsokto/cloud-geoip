package main

import (
	"net/netip"

	"go4.org/netipx"
)

func filterPrefixes(prefixes []string, noV4, noV6 bool) ([]string, error) {
	kept := make([]string, 0, len(prefixes))
	for _, p := range prefixes {
		prefix, err := netip.ParsePrefix(p)
		if err != nil {
			return nil, err
		}
		if (prefix.Addr().Is4() && noV4) || (prefix.Addr().Is6() && noV6) {
			continue
		}
		kept = append(kept, p)
	}
	return kept, nil
}

func aggregatePrefixes(prefixes []string) ([]string, error) {
	var b netipx.IPSetBuilder
	for _, p := range prefixes {
		prefix, err := netip.ParsePrefix(p)
		if err != nil {
			return nil, err
		}
		b.AddPrefix(prefix)
	}
	set, err := b.IPSet()
	if err != nil {
		return nil, err
	}
	merged := make([]string, 0, len(prefixes))
	for _, p := range set.Prefixes() {
		merged = append(merged, p.String())
	}
	return merged, nil
}

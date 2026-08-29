package main

import (
	"reflect"
	"testing"
)

func TestFilterPrefixes(t *testing.T) {
	all := []string{"1.2.3.0/24", "2.16.0.0/13", "2001:db8::/32"}

	got, err := filterPrefixes(all, false, true)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"1.2.3.0/24", "2.16.0.0/13"}; !reflect.DeepEqual(got, want) {
		t.Errorf("no-v6: got %v, want %v", got, want)
	}

	got, err = filterPrefixes(all, true, false)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"2001:db8::/32"}; !reflect.DeepEqual(got, want) {
		t.Errorf("no-v4: got %v, want %v", got, want)
	}

	if _, err := filterPrefixes([]string{"not-a-cidr"}, false, false); err == nil {
		t.Error("expected error on malformed CIDR")
	}
}

func TestAggregatePrefixes(t *testing.T) {
	got, err := aggregatePrefixes([]string{
		"2.16.0.0/13", "2.16.0.0/24", "2.16.1.0/24",
		"9.9.9.0/25", "9.9.9.128/25",
		"2001:db8::/33", "2001:db8:8000::/33",
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"2.16.0.0/13", "9.9.9.0/24", "2001:db8::/32"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

package main

import (
	"testing"
	"time"
)

func TestParseDateDefaultsDateOnlyToNoonInLocation(t *testing.T) {
	loc, err := time.LoadLocation(defaultTimeZone)
	if err != nil {
		t.Fatal(err)
	}

	got, err := parseDate("2026-06-18", loc)
	if err != nil {
		t.Fatal(err)
	}

	want := time.Date(2026, 6, 18, 12, 0, 0, 0, loc)
	if !got.Equal(want) {
		t.Fatalf("parseDate() = %s, want %s", got, want)
	}
}

func TestParseDatePreservesExplicitTime(t *testing.T) {
	loc, err := time.LoadLocation(defaultTimeZone)
	if err != nil {
		t.Fatal(err)
	}

	got, err := parseDate("2026-06-18 09:30", loc)
	if err != nil {
		t.Fatal(err)
	}

	want := time.Date(2026, 6, 18, 9, 30, 0, 0, loc)
	if !got.Equal(want) {
		t.Fatalf("parseDate() = %s, want %s", got, want)
	}
}

func TestPrepareDataUsesTwentyFourHourTime(t *testing.T) {
	data := PrepareData(MoonDatePhase{
		NextNewMoon: time.Date(2026, 7, 14, 14, 45, 0, 0, time.UTC),
	})

	if data.NextNewMoon != "July 14, 2026 at 14:45" {
		t.Fatalf("NextNewMoon = %q", data.NextNewMoon)
	}
}

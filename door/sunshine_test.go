package door

import (
	"strconv"
	"testing"
	"time"
)

func TestWinter(t *testing.T) {
	lat, _ := strconv.ParseFloat("51.613078", 64)
	long, _ := strconv.ParseFloat("-0.165323", 64)
	start, _ := time.ParseDuration("6h30m")
	end, _ := time.ParseDuration("22h")

	// Sunrise: 2017-11-13 07:14:28 +0000 UTC
	// Sunset: 2017-11-13 16:12:10 +0000 UTC

	now, _ := time.Parse(time.RFC3339, "2017-11-13T06:29:59Z")
	if adjust(now, lat, long, start, end) {
		t.Errorf("Expected false\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-11-13T06:30:01Z")
	if !adjust(now, lat, long, start, end) {
		t.Errorf("Expected true\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-11-13T07:14:27Z")
	if !adjust(now, lat, long, start, end) {
		t.Errorf("Expected true\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-11-13T07:15:00Z")
	if adjust(now, lat, long, start, end) {
		t.Errorf("Expected false\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-11-13T12:00:00Z")
	if adjust(now, lat, long, start, end) {
		t.Errorf("Expected false\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-11-13T16:11:00Z")
	if adjust(now, lat, long, start, end) {
		t.Errorf("Expected false\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-11-13T16:14:00Z")
	if !adjust(now, lat, long, start, end) {
		t.Errorf("Expected true\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-11-13T21:59:00Z")
	if !adjust(now, lat, long, start, end) {
		t.Errorf("Expected true\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-11-13T22:01:00Z")
	if adjust(now, lat, long, start, end) {
		t.Errorf("Expected false\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-11-13T23:59:59Z")
	if adjust(now, lat, long, start, end) {
		t.Errorf("Expected false\n")
	}
}

func TestSummer(t *testing.T) {
	lat, _ := strconv.ParseFloat("51.613078", 64)
	long, _ := strconv.ParseFloat("-0.165323", 64)
	start, _ := time.ParseDuration("6h30m")
	end, _ := time.ParseDuration("22h")

	// Sunrise: 2017-06-21 03:41:25 +0000 UTC
	// Sunset: 2017-06-21 20:21:02 +0000 UTC

	now, _ := time.Parse(time.RFC3339, "2017-06-21T03:52:59Z")
	if adjust(now, lat, long, start, end) {
		t.Errorf("Expected false\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-06-21T06:29:59Z")
	if adjust(now, lat, long, start, end) {
		t.Errorf("Expected false\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-06-21T06:30:01Z")
	if adjust(now, lat, long, start, end) {
		t.Errorf("Expected false\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-06-21T12:00:00Z")
	if adjust(now, lat, long, start, end) {
		t.Errorf("Expected false\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-06-21T20:20:00Z")
	if adjust(now, lat, long, start, end) {
		t.Errorf("Expected false\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-06-21T20:22:00Z")
	if !adjust(now, lat, long, start, end) {
		t.Errorf("Expected true\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-06-21T21:59:00Z")
	if !adjust(now, lat, long, start, end) {
		t.Errorf("Expected true\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-06-21T22:01:00Z")
	if adjust(now, lat, long, start, end) {
		t.Errorf("Expected false\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-06-21T23:59:59Z")
	if adjust(now, lat, long, start, end) {
		t.Errorf("Expected false\n")
	}
}

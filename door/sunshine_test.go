package door

import (
	"strconv"
	"testing"
	"time"

	"github.com/cpucycle/astrotime"
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
	now, _ = time.Parse(time.RFC3339, "2017-11-13T07:14:29Z")
	if adjust(now, lat, long, start, end) {
		t.Errorf("Expected false\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-11-13T12:00:00Z")
	if adjust(now, lat, long, start, end) {
		t.Errorf("Expected false\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-11-13T16:12:09Z")
	if adjust(now, lat, long, start, end) {
		t.Errorf("Expected false\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-11-13T16:12:11Z")
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

	// Sunrise: 2017-06-21 03:41:21 +0000 UTC
	// Sunset: 2017-06-21 20:20:59 +0000 UTC

	now, _ := time.Parse(time.RFC3339, "2017-06-21T03:41:20Z")
	if adjust(now, lat, long, start, end) {
		t.Errorf("Expected false\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-06-21T03:41:22Z")
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
	now, _ = time.Parse(time.RFC3339, "2017-06-21T20:20:58Z")
	if adjust(now, lat, long, start, end) {
		t.Errorf("Expected false\n")
	}
	now, _ = time.Parse(time.RFC3339, "2017-06-21T20:21:00Z")
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

// There appears to be an issue when calculating sunrise/set. The calculation is sensitive
// to the time component of the time instant passed into `CalcSunrise` and `CalcSunset`;
// if for a given day a pass I in multiple instants, each with a slightly different time
// component, the calculates for sunrise/set for that day will vary per instant. Here is
// an example tests case. To work around this I'm simply setting the time component of the
// input instant to `00:00:00+0ns`
func TestDrift(t *testing.T) {
	lat, _ := strconv.ParseFloat("56.613078", 64)
	long, _ := strconv.ParseFloat("-2.165323", 64)

	now1, _ := time.Parse(time.RFC3339, "2017-11-13T01:00:00Z")
	sunrise1 := astrotime.CalcSunrise(now1, lat, long)
	sunset1 := astrotime.CalcSunset(now1, lat, long)

	now2, _ := time.Parse(time.RFC3339, "2017-11-13T02:00:00Z")
	sunrise2 := astrotime.CalcSunrise(now2, lat, long)
	sunset2 := astrotime.CalcSunset(now2, lat, long)

	if !sunrise1.Equal(sunrise2) {
		t.Errorf("Expected sunrise instants to match but %v != %v\n", sunrise1, sunrise2)
	}
	if !sunset1.Equal(sunset2) {
		t.Errorf("Expected sunset instants to match but %v != %v\n", sunset1, sunset2)
	}
}

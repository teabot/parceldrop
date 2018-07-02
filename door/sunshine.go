package door

import (
	"log"
	"strconv"
	"time"

	"github.com/cpucycle/astrotime"
)

var lastState bool

func CheckSunRise(latitude, longitude string, dayStart, dayEnd time.Duration, change func(bool, bool)) {
	lat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		log.Fatalf("Invalid latitude: %v\n", latitude)
	}
	long, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		log.Fatalf("Invalid longitude: %v\n", longitude)
	}

	go func() {
		change(false, false)
		for {
			nextState, night := adjust(time.Now(), lat, long, dayStart, dayEnd)
			if nextState != lastState {
				log.Printf("SUN: Changed light state to %v\n", nextState)
			}
			lastState = nextState
			change(nextState, night)
			time.Sleep(60 * time.Second)
		}
	}()
}

func adjust(now time.Time, latitude, longitude float64, dayStart, dayEnd time.Duration) (bool, bool) {
	nowDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	start := nowDay.Add(dayStart)
	end := nowDay.Add(dayEnd)
	sunrise := astrotime.CalcSunrise(nowDay, latitude, longitude)
	sunset := astrotime.CalcSunset(nowDay, latitude, longitude)

	// fmt.Printf("Sunrise %v\n", sunrise)
	// fmt.Printf("Sunset  %v\n", sunset)
	// fmt.Printf("Now     %v\n", now)

	if now.After(start) && now.Before(end) {
		if now.Before(sunrise) || now.After(sunset) {
			return true, true
		}
	}
	return false, now.Before(sunrise) || now.After(sunset)
}

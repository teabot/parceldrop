package codebook

import "log"
import "time"

// IsoTimestamp x
type IsoTimestamp string

// Day 3 letter x
type Day string

// IsoTime x
type IsoTime string

// CodeType
type CodeType string

// AccessCode x
type AccessCode struct {
	Digits string
	Name   string
	Types  []CodeType
	// Duration
	ValidityHours int64
	FirstUse      IsoTimestamp
	// Interval
	ValidFrom  IsoTimestamp
	Expiration IsoTimestamp
	// Count
	MaxUsage int
	Usage    int
	// Day
	Days      []Day
	StartTime IsoTime
	EndTime   IsoTime
}

const (
	Active    CodeType = "active"
	Duration  CodeType = "duration"
	Interval  CodeType = "interval"
	Count     CodeType = "count"
	DayFilter CodeType = "day"

	ISO8601 = "2006-01-02T15:04:05-0700"
)

var masterCode string

func Initialise(adminCode, defaultCode string) error {
	if len(adminCode) > 0 {
		masterCode = adminCode
	} else {
		masterCode = defaultCode
	}

	err := OpenStore()
	if err != nil {
		return err
	}
	return nil
}

func Check(digits string, now time.Time) bool {
	if digits == masterCode {
		log.Printf("CODEBOOK: Matched admin code: %v\n", digits)
		return true
	}
	code, _ := GetAccessCode(digits)
	if code == nil {
		log.Printf("CODEBOOK: Code not found: %v\n", digits)
		return false
	}
	if !code.hasType(Active) {
		log.Printf("CODEBOOK: Code not active: %v\n", digits)
		return false
	}
	if code.hasType(Count) && (code.Usage < 0 || code.MaxUsage < 1 || code.Usage >= code.MaxUsage) {
		log.Printf("CODEBOOK: Usage exceeded: %v, %v>=%v\n", digits, code.Usage, code.MaxUsage)
		return false
	}

	nowStr := now.Format(ISO8601)
	if code.hasType(Duration) && code.Usage > 0 {
		firstUse, err := time.Parse(ISO8601, string(code.FirstUse))
		if err != nil {
			log.Printf("CODEBOOK: Duration: Invalid first use: %v, %v\n", digits, code.FirstUse)
			return false
		}
		duration := time.Duration(code.ValidityHours) * time.Hour
		if firstUse.Add(duration).Before(now) {
			log.Printf("CODEBOOK: Duration expired: %v, %v, %vhrs\n", digits, code.FirstUse, code.ValidityHours)
			return false
		}
	}
	if code.hasType(Interval) {
		from, err := time.Parse(ISO8601, string(code.ValidFrom))
		if err != nil {
			log.Printf("CODEBOOK: Interval: Invalid from: %v, %v\n", digits, code.ValidFrom)
			return false
		}
		to, err := time.Parse(ISO8601, string(code.Expiration))
		if err != nil {
			log.Printf("CODEBOOK: Interval: Invalid to: %v, %v\n", digits, code.Expiration)
			return false
		}
		if now.Before(from) || now.After(to) {
			log.Printf("CODEBOOK: Outside of interval: %v, %v -> %v\n", digits, code.ValidFrom, code.Expiration)
			return false
		}
	}
	if code.hasType(DayFilter) {
		if now.Hour() < 7 || now.Hour() > 21 {
			log.Printf("CODEBOOK: Outside of day pattern: %v\n", digits)
			return false
		}
	}
	if code.Usage == 0 {
		code.FirstUse = IsoTimestamp(nowStr)
	}
	code.Usage++
	code.save()
	log.Printf("Updated code: %v\n", code)
	return true
}

func Close() {
	CloseStore()
}

func (c *AccessCode) hasType(ct CodeType) bool {
	for _, t := range c.Types {
		if t == ct {
			return true
		}
	}
	return false
}

func Rescind(digits *string) error {
	log.Printf("CODEBOOK: Rescinding code: %v\n", digits)
	code, err := GetAccessCode(*digits)
	if err != nil {
		log.Printf("CODEBOOK: Error getting code: %v", digits)
		return err
	}
	if code == nil {
		log.Printf("CODEBOOK: No code to rescind: %v", digits)
		return nil
	}
	if !code.hasType(Active) {
		log.Printf("CODEBOOK: Code already rescinded: %v", code)
		return nil
	}
	code.Types = remove(code.Types, Active)
	err2 := code.save()
	if err2 != nil {
		log.Printf("CODEBOOK: Error saving rescinded code: %v, %v", code, err2)
		return err
	}
	return nil
}

func Update(code *AccessCode) error {
	log.Printf("CODEBOOK: Updating code: %v\n", code)
	err := code.save()
	return err
}

func remove(l []CodeType, item CodeType) []CodeType {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

package codebook

import (
	"testing"
	"time"
)

func TestMaster(t *testing.T) {
	err := Initialise("", "123456", "999999", time.ParseDuration("6h30m"), time.ParseDuration("21h30m"))
	defer Close()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}
	check, silent, name := Check("123456", time.Now().UTC())
	if !check {
		t.Errorf("Expected check to succeed\n")
	}
	if !silent {
		t.Errorf("Expected silent\n")
	}
	if name != "Master" {
		t.Errorf("Expected check to succeed\n")
	}
}

func TestDefault(t *testing.T) {
	err := Initialise("", "", "999999", time.ParseDuration("6h30m"), time.ParseDuration("21h30m"))
	defer Close()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}
	check, silent, name := Check("999999", time.Now().UTC())
	if !check {
		t.Errorf("Expected check to succeed\n")
	}
	if !silent {
		t.Errorf("Expected silent\n")
	}
	if name != "Master" {
		t.Errorf("Expected check to succeed\n")
	}
}

func TestUnknown(t *testing.T) {
	err := Initialise("", "123456", "999999", time.ParseDuration("6h30m"), time.ParseDuration("21h30m"))
	defer Close()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}
	check, _, name := Check("5555", time.Now().UTC())
	if check {
		t.Errorf("Expected check to fail\n")
	}
	if name != "" {
		t.Errorf("Expected check to fail\n")
	}
}

func TestValidCount(t *testing.T) {
	err := Initialise("", "123456", "999999", time.ParseDuration("6h30m"), time.ParseDuration("21h30m"))
	defer Close()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}

	code := AccessCode{
		Types:    []CodeType{"active", "count"},
		Name:     "Test",
		Digits:   "6789",
		MaxUsage: 3,
		Usage:    1,
	}
	code.save()

	check, _, _ := Check("6789", time.Now().UTC())
	if !check {
		t.Errorf("Expected check to succeed\n")
	}

	actual, _ := GetAccessCode("6789")
	if actual.Usage != 2 {
		t.Errorf("Expected usage to be 2\n")
	}
}

func TestSilent(t *testing.T) {
	err := Initialise("", "123456", "999999", time.ParseDuration("6h30m"), time.ParseDuration("21h30m"))
	defer Close()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}

	code := AccessCode{
		Types:  []CodeType{"active", "silent"},
		Name:   "Test",
		Digits: "6789",
	}
	code.save()

	check, silent, _ := Check("6789", time.Now().UTC())
	if !check {
		t.Errorf("Expected check to succeed\n")
	}
	if !silent {
		t.Errorf("Expected silent\n")
	}
}

func TestNotSilent(t *testing.T) {
	err := Initialise("", "123456", "999999", time.ParseDuration("6h30m"), time.ParseDuration("21h30m"))
	defer Close()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}

	code := AccessCode{
		Types:  []CodeType{"active"},
		Name:   "Test",
		Digits: "6789",
	}
	code.save()

	check, silent, _ := Check("6789", time.Now().UTC())
	if !check {
		t.Errorf("Expected check to succeed\n")
	}
	if silent {
		t.Errorf("Expected not silent\n")
	}
}

func TestValidCountFirstUse(t *testing.T) {
	err := Initialise("", "123456", "999999", time.ParseDuration("6h30m"), time.ParseDuration("21h30m"))
	defer Close()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}

	code := AccessCode{
		Types:    []CodeType{"active", "count"},
		Name:     "Test",
		Digits:   "6789",
		MaxUsage: 3,
		Usage:    0,
	}
	code.save()

	now := time.Now().UTC()
	check, _, _ := Check("6789", now)
	if !check {
		t.Errorf("Expected check to succeed\n")
	}

	actual, _ := GetAccessCode("6789")
	if actual.Usage != 1 {
		t.Errorf("Expected usage to be 1\n")
	}
	if actual.FirstUse != IsoTimestamp(now.Format(ISO8601)) {
		t.Errorf("Expected first use to be %v\n", now)
	}
}

func TestInvalidCount(t *testing.T) {
	err := Initialise("", "123456", "999999", time.ParseDuration("6h30m"), time.ParseDuration("21h30m"))
	defer Close()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}

	code := AccessCode{
		Types:    []CodeType{"active", "count"},
		Name:     "Test",
		Digits:   "6789",
		MaxUsage: 3,
		Usage:    4,
	}
	code.save()

	check, _, _ := Check("6789", time.Now().UTC())
	if check {
		t.Errorf("Expected check to fail\n")
	}
}

func TestInvalidType(t *testing.T) {
	err := Initialise("", "123456", "999999", time.ParseDuration("6h30m"), time.ParseDuration("21h30m"))
	defer Close()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}

	code := AccessCode{
		Types:  []CodeType{},
		Name:   "Test",
		Digits: "6789",
	}
	code.save()

	check, _, _ := Check("6789", time.Now().UTC())
	if check {
		t.Errorf("Expected check to fail\n")
	}
}

func TestInactiveType(t *testing.T) {
	err := Initialise("", "123456", "999999", time.ParseDuration("6h30m"), time.ParseDuration("21h30m"))
	defer Close()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}

	code := AccessCode{
		Types:  []CodeType{"count"},
		Name:   "Test",
		Digits: "6789",
	}
	code.save()

	check, _, _ := Check("6789", time.Now().UTC())
	if check {
		t.Errorf("Expected check to fail\n")
	}
}

func TestDurationInvalidFirstUse(t *testing.T) {
	err := Initialise("", "123456", "999999", time.ParseDuration("6h30m"), time.ParseDuration("21h30m"))
	defer Close()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}

	code := AccessCode{
		Types:    []CodeType{"active", "duration"},
		Name:     "Test",
		Digits:   "6789",
		FirstUse: "jsjsjsjsj",
		Usage:    1,
	}
	code.save()

	check, _, _ := Check("6789", time.Now().UTC())
	if check {
		t.Errorf("Expected check to fail\n")
	}
}

func TestDurationExpired(t *testing.T) {
	err := Initialise("", "123456", "999999", time.ParseDuration("6h30m"), time.ParseDuration("21h30m"))
	defer Close()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}

	now := time.Now().UTC()
	code := AccessCode{
		Types:         []CodeType{"active", "duration"},
		Name:          "Test",
		Digits:        "6789",
		FirstUse:      IsoTimestamp(now.Add(-2 * time.Hour).Format(ISO8601)),
		ValidityHours: 1,
		Usage:         1,
	}
	code.save()

	check, _, _ := Check("6789", time.Now().UTC())
	if check {
		t.Errorf("Expected check to fail\n")
	}
}

func TestDurationOk(t *testing.T) {
	err := Initialise("", "123456", "999999", time.ParseDuration("6h30m"), time.ParseDuration("21h30m"))
	defer Close()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}

	now := time.Now().UTC()
	code := AccessCode{
		Types:         []CodeType{"active", "duration"},
		Name:          "Test",
		Digits:        "6789",
		FirstUse:      IsoTimestamp(now.Add(-1 * time.Hour).Format(ISO8601)),
		ValidityHours: 2,
		Usage:         1,
	}
	code.save()

	check, _, _ := Check("6789", time.Now().UTC())
	if !check {
		t.Errorf("Expected check to succeed\n")
	}

	actual, _ := GetAccessCode("6789")
	if actual.Usage != 2 {
		t.Errorf("Expected usage to be 2\n")
	}
}

func TestIntervalInvalidFrom(t *testing.T) {
	err := Initialise("", "123456", "999999", time.ParseDuration("6h30m"), time.ParseDuration("21h30m"))
	defer Close()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}

	now := time.Now().UTC()
	code := AccessCode{
		Types:      []CodeType{"active", "interval"},
		Name:       "Test",
		Digits:     "6789",
		ValidFrom:  "jsjsjsjsj",
		Expiration: IsoTimestamp(now.Add(1 * time.Hour).Format(ISO8601)),
	}
	code.save()

	check, _, _ := Check("6789", time.Now().UTC())
	if check {
		t.Errorf("Expected check to fail\n")
	}
}

func TestIntervalInvalidTo(t *testing.T) {
	err := Initialise("", "123456", "999999", time.ParseDuration("6h30m"), time.ParseDuration("21h30m"))
	defer Close()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}

	now := time.Now().UTC()
	code := AccessCode{
		Types:      []CodeType{"active", "interval"},
		Name:       "Test",
		Digits:     "6789",
		Expiration: "jsjsjsjsj",
		ValidFrom:  IsoTimestamp(now.Add(-1 * time.Hour).Format(ISO8601)),
	}
	code.save()

	check, _, _ := Check("6789", time.Now().UTC())
	if check {
		t.Errorf("Expected check to fail\n")
	}
}

func TestIntervalOutside(t *testing.T) {
	err := Initialise("", "123456", "999999", time.ParseDuration("6h30m"), time.ParseDuration("21h30m"))
	defer Close()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}

	now := time.Now().UTC()
	code := AccessCode{
		Types:      []CodeType{"active", "interval"},
		Name:       "Test",
		Digits:     "6789",
		ValidFrom:  IsoTimestamp(now.Add(-1 * time.Hour).Format(ISO8601)),
		Expiration: IsoTimestamp(now.Add(1 * time.Hour).Format(ISO8601)),
	}
	code.save()

	check, _, _ := Check("6789", now.Add(2*time.Hour))
	if check {
		t.Errorf("Expected check to fail\n")
	}
	check, _, _ = Check("6789", now.Add(-2*time.Hour))
	if check {
		t.Errorf("Expected check to fail\n")
	}
}

func TestIntervalInside(t *testing.T) {
	err := Initialise("", "123456", "999999", time.ParseDuration("6h30m"), time.ParseDuration("21h30m"))
	defer Close()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}

	now := time.Now().UTC()
	code := AccessCode{
		Types:      []CodeType{"active", "interval"},
		Name:       "Test",
		Digits:     "6789",
		ValidFrom:  IsoTimestamp(now.Add(-1 * time.Hour).Format(ISO8601)),
		Expiration: IsoTimestamp(now.Add(1 * time.Hour).Format(ISO8601)),
	}
	code.save()

	check, _, _ := Check("6789", now)
	if !check {
		t.Errorf("Expected check to succeed\n")
	}
}

func TestDayInside(t *testing.T) {
	err := Initialise("", "123456", "999999", time.ParseDuration("6h30m"), time.ParseDuration("21h30m"))
	defer Close()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}

	code := AccessCode{
		Types:    []CodeType{"active", "count", "day"},
		Name:     "Test",
		Digits:   "6789",
		Usage:    0,
		MaxUsage: 1,
	}
	code.save()

	now, _ := time.Parse(ISO8601, "2017-11-09T11:10:03+0000")
	check, _, _ := Check("6789", now)
	if !check {
		t.Errorf("Expected check to succeed\n")
	}
}

func TestDayOutside(t *testing.T) {
	err := Initialise("", "123456", "999999", time.ParseDuration("6h30m"), time.ParseDuration("21h30m"))
	defer Close()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}

	code := AccessCode{
		Types:    []CodeType{"active", "count", "day"},
		Name:     "Test",
		Digits:   "6789",
		Usage:    0,
		MaxUsage: 1,
	}
	code.save()

	now1, _ := time.Parse(ISO8601, "2017-11-09T06:29:03+0000")
	check, _, _ := Check("6789", now1)
	if check {
		t.Errorf("Expected check to fail\n")
	}
	now2, _ := time.Parse(ISO8601, "2017-11-09T21:31:03+0000")
	check, _, _ = Check("6789", now2)
	if check {
		t.Errorf("Expected check to fail\n")
	}
}

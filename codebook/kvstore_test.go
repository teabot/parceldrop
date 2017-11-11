package codebook

import (
	"testing"
)

func TestOpenClose(t *testing.T) {
	err := OpenStore()
	defer CloseStore()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}
}

func TestSaveAndGet(t *testing.T) {
	err := OpenStore()
	defer CloseStore()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}

	code := AccessCode{
		Types:         []CodeType{"active", "duration", "interval", "count"},
		Name:          "Test",
		Digits:        "123456",
		ValidityHours: 4,
		FirstUse:      "2006-01-02T15:04:05-0700",
		ValidFrom:     "2006-02-02T15:04:05-0700",
		Expiration:    "2006-03-02T15:04:05-0700",
		MaxUsage:      3,
		Usage:         1,
	}
	code.save()

	actual, err := GetAccessCode("123456")
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}
	if "123456" != actual.Digits {
		t.Errorf("Expected Digits %v, got %v\n", code.Digits, actual.Digits)
	}
	if "Test" != actual.Name {
		t.Errorf("Expected Name %v, got %v\n", code.Name, actual.Name)
	}
	if 4 != actual.ValidityHours {
		t.Errorf("Expected ValidityHours %v, got %v\n", code.ValidityHours, actual.ValidityHours)
	}
	if "2006-01-02T15:04:05-0700" != actual.FirstUse {
		t.Errorf("Expected FirstUse %v, got %v\n", code.FirstUse, actual.FirstUse)
	}
	if "2006-02-02T15:04:05-0700" != actual.ValidFrom {
		t.Errorf("Expected ValidFrom %v, got %v\n", code.ValidFrom, actual.ValidFrom)
	}
	if "2006-03-02T15:04:05-0700" != actual.Expiration {
		t.Errorf("Expected Expiration %v, got %v\n", code.Expiration, actual.Expiration)
	}
	if 3 != actual.MaxUsage {
		t.Errorf("Expected MaxUsage %v, got %v\n", code.MaxUsage, actual.MaxUsage)
	}
	if 1 != actual.Usage {
		t.Errorf("Expected Usage %v, got %v\n", code.Usage, actual.Usage)
	}
	if !actual.hasType("active") {
		t.Errorf("Expected active code\n")
	}
	if !actual.hasType("duration") {
		t.Errorf("Expected duration code\n")
	}
	if !actual.hasType("interval") {
		t.Errorf("Expected interval code\n")
	}
	if !actual.hasType("count") {
		t.Errorf("Expected count code\n")
	}
}

func TestSaveAndGetOther(t *testing.T) {
	err := OpenStore()
	defer CloseStore()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}

	code := AccessCode{
		Types:         []CodeType{"active", "duration", "interval", "count"},
		Name:          "Test",
		Digits:        "123456",
		ValidityHours: 4,
		FirstUse:      "2006-01-02T15:04:05-0700",
		ValidFrom:     "2006-02-02T15:04:05-0700",
		Expiration:    "2006-03-02T15:04:05-0700",
		MaxUsage:      3,
		Usage:         1,
	}
	code.save()

	actual, err := GetAccessCode("4444")
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}
	if actual != nil {
		t.Errorf("Expected nil\n")
	}
}

func TestGetOther(t *testing.T) {
	err := OpenStore()
	defer CloseStore()
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}

	actual, err := GetAccessCode("4444")
	if err != nil {
		t.Errorf("Expected no error, got %v\n", err)
	}
	if actual != nil {
		t.Errorf("Expected nil\n")
	}
}

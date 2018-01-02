package sms

import (
	"testing"
)

func TestRedactMoreThan5(t *testing.T) {
	if redactCode("465757") != "4****7" {
		t.Errorf("Expected 4****7\n")
	}
}

func TestRedact5(t *testing.T) {
	if redactCode("46575") != "4****" {
		t.Errorf("Expected 4****\n")
	}
}

func TestRedactLessThan5(t *testing.T) {
	if redactCode("4652") != "****" {
		t.Errorf("Expected ****\n")
	}
	if redactCode("465") != "****" {
		t.Errorf("Expected ****\n")
	}
	if redactCode("46") != "****" {
		t.Errorf("Expected ****\n")
	}
	if redactCode("4") != "****" {
		t.Errorf("Expected ****\n")
	}
	if redactCode("") != "****" {
		t.Errorf("Expected ****\n")
	}
}

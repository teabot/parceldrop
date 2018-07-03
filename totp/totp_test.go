package totp

import (
	"testing"
	"time"
)

const (
	ISO8601 = "2006-01-02T15:04:05-0700"
)

func TestTotp(t *testing.T) {
	//Tue Jul  3 06:23:07 BST 2018
	now, _ := time.Parse(ISO8601, "2018-07-03T08:01:26+0100")
	period, _ := time.ParseDuration("300s")
	Initialise("6s634i355bxjnxmfyecw6ktt4v4sydywl7cggjxb3toye5p4wegsqvhb", period)
	check := Validate("504110", now)
	if !check {
		t.Errorf("Expected check to succeed\n")
	}
	check = Validate("504111", now)
	if check {
		t.Errorf("Expected check to fail\n")
	}
}

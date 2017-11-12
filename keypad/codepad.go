package keypad

import (
	"bytes"
	"log"
	"time"
)

// SubmissionType x
type SubmissionType string

const (
	// User x
	User SubmissionType = "user"
	// Partial y
	Partial = "partial"
	// Final z
	Final = "final"
)

// Code is
type Code struct {
	Digits    string
	Submitted SubmissionType
}

// ScanCodes x
func ScanCodes(device string, autoClear time.Duration, maxLength int, codeFn func(Code), timeoutFn func(string), errFn func(error)) {
	var keys = make(chan keyEvent)

	go ScanKeys(device, keys, errFn)

	var scanBuffer bytes.Buffer
	var autoClearTimer *time.Timer

	reset := func() {
		lastCode := scanBuffer.String()
		scanBuffer.Reset()
		timeoutFn(lastCode)
	}

	for {
		select {
		case k := <-keys:
			if autoClearTimer != nil {
				autoClearTimer.Stop()
			}
			if k.Enter {
				if scanBuffer.Len() > 0 {
					log.Printf("KEYPAD: User submitted code: %v\n", scanBuffer.String())
					codeFn(Code{
						Digits:    scanBuffer.String(),
						Submitted: User,
					})
				}
				scanBuffer.Reset()
			} else if k.Clear {
				log.Println("KEYPAD: Manual clear")
				scanBuffer.Reset()
			} else {
				scanBuffer.WriteString(k.Digit)
				if scanBuffer.Len() >= maxLength {
					codeFn(Code{
						Digits:    scanBuffer.String(),
						Submitted: Final,
					})
					scanBuffer.Reset()
				} else {
					autoClearTimer = startTimer(autoClear, reset)
					codeFn(Code{
						Digits:    scanBuffer.String(),
						Submitted: Partial,
					})
				}
			}
		}
	}
}

func startTimer(duration time.Duration, exec func()) *time.Timer {
	countDown := time.NewTimer(duration)
	go func() {
		<-countDown.C
		exec()
	}()
	return countDown
}

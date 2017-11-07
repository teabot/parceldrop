// Copyright Banrai LLC. All rights reserved. Use of this source code is
// governed by the license that can be found in the LICENSE file.

// Package scanner provides functions for reading barcode scans from
// usb-connected barcode scanner devices as if they were keyboards, i.e.,
// by using the corresponding '/dev/input/event' device, inspired by this
// post on linuxquestions.org:
//
// http://www.linuxquestions.org/questions/programming-9/read-from-a-usb-barcode-scanner-that-simulates-a-keyboard-495358/#post2767643
//
// Also found important Go-specific information by reviewing the code from
// this repo on github:
//
// https://github.com/gvalkov/golang-evdev

package main

import (
	"log"
	"time"

	"github.com/teabot/parceldrop/door"
	"github.com/teabot/parceldrop/keypad"
	"github.com/teabot/parceldrop/sms"
)

const (
	doorOpenDuration       time.Duration = 60 * time.Second
	resetToIdleDuration    time.Duration = 2 * time.Second
	resetCodeInputDuration time.Duration = 5 * time.Second

	maxCodeLength int = 6

	defaultCode string = "1234"
)

var timer *time.Timer

func main() {
	door.Initialise()
	sms.Initialise()

	codeFn := func(code keypad.Code) {
		if door.Locked() {
			unschedule()
			door.Wait()
			log.Printf("Code: %v, Submitted: %v\n", code.Digits, code.Submitted)
			if code.Digits == defaultCode {
				validCode(code.Digits)
			} else {
				if code.Submitted == keypad.Final || code.Submitted == keypad.User {
					invalidCode(code.Digits)
				}
			}
		}
	}

	keyTimeoutFn := func(digits string) {
		log.Println("Auto clear")
		codeFn(keypad.Code{
			Digits:    digits,
			Submitted: keypad.User,
		})
	}

	errorFn := func(e error) {
		panic(e)
	}

	keypad.ScanCodes(keypad.DefaultDevice, resetCodeInputDuration, maxCodeLength, codeFn, keyTimeoutFn, errorFn)
}

func validCode(digits string) {
	log.Printf("Opened with code: %v\n", digits)
	door.Unlock()
	sms.SendCorrectCode(digits)
	schedule(resetToLocked())
}

func invalidCode(digits string) {
	log.Printf("Invalid code: %v\n", digits)
	door.Reject()
	sms.SendInvalidCode(digits)
	schedule(resetToIdle())
}

func schedule(nextTimer *time.Timer) {
	unschedule()
	timer = nextTimer
}

func unschedule() {
	if timer != nil {
		timer.Stop()
	}
}

func resetToLocked() *time.Timer {
	resetToLocked := time.NewTimer(doorOpenDuration)
	go func() {
		select {
		case <-resetToLocked.C:
			if door.Open() {
				sms.SendDoorNotClosed()
				log.Println("Door not closed")
			}
		default:
			for door.Open() {
				time.Sleep(500 * time.Millisecond)
			}
			log.Println("Detected door close")
		}
		door.Lock()
	}()
	return resetToLocked
}

func resetToIdle() *time.Timer {
	resetToIdle := time.NewTimer(resetToIdleDuration)
	go func() {
		<-resetToIdle.C
		door.Lock()
	}()
	return resetToIdle
}

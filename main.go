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
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/teabot/parceldrop/control"

	"github.com/teabot/parceldrop/codebook"

	"github.com/teabot/parceldrop/door"
	"github.com/teabot/parceldrop/keypad"
	"github.com/teabot/parceldrop/sms"
)

const (
	doorOpenDuration       = 60 * time.Second
	resetToIdleDuration    = 3 * time.Second
	resetCodeInputDuration = 5 * time.Second

	maxCodeLength = 6
	defaultCode   = "1234"
)

var timer *time.Timer

func main() {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	codebook.Initialise(os.Getenv("ADMIN_CODE"), defaultCode)
	sms.Initialise(strings.Split(os.Getenv("SMS_DESTINATIONS"), ","))
	door.Initialise(overrideOpen)
	control.InitialiseSqs(os.Getenv("AWS_SQS_QUEUE"), overrideOpen)

	go func() {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v", sig)
		fmt.Println("Wait for 2 second to finish processing")
		codebook.CloseStore()
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

	codeFn := func(code keypad.Code) {
		if door.Locked() {
			unscheduleAnyEvents()
			door.Wait()
			log.Printf("MAIN: Code: %v, Submitted: %v\n", code.Digits, code.Submitted)
			if codebook.Check(code.Digits, time.Now().UTC()) {
				validCode(code.Digits)
			} else {
				if code.Submitted == keypad.Final || code.Submitted == keypad.User {
					invalidCode(code.Digits)
				}
			}
		}
	}

	keyTimeoutFn := func(digits string) {
		log.Println("MAIN: Auto clear")
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
	log.Printf("MAIN: Opened with code: %v\n", digits)
	door.Unlock()
	sms.SendCorrectCode(digits)
	scheduleEvent(resetToLocked())
}

func overrideOpen(overrideType string) {
	scheduleEvent(resetToLocked())
	log.Printf("MAIN: Opened with override: %v\n", overrideType)
	door.Unlock()
	sms.SendOverrideOpen(overrideType)
}

func invalidCode(digits string) {
	scheduleEvent(resetToIdle())
	log.Printf("MAIN: Invalid code: %v\n", digits)
	door.Reject()
	sms.SendInvalidCode(digits)
}

func scheduleEvent(nextTimer *time.Timer) {
	unscheduleAnyEvents()
	timer = nextTimer
}

func unscheduleAnyEvents() {
	if timer != nil {
		timer.Stop()
		timer = nil
		log.Println("MAIN: Unscheduled timer")
	}
}

func resetToLocked() *time.Timer {
	resetToLocked := time.NewTimer(doorOpenDuration)
	go checkDoorClosed(resetToLocked)
	log.Println("MAIN: Scheduled resetToLocked")
	return resetToLocked
}

func checkDoorClosed(resetToLocked *time.Timer) {
	checkContactClosed := make(chan bool)
	go contactCheck(checkContactClosed)

	select {
	case <-resetToLocked.C:
		log.Println("MAIN: resetToLocked returned")
		if door.Open() {
			sms.SendDoorNotClosed()
			log.Println("MAIN: Door not closed")
		}
	case <-checkContactClosed:
		unscheduleAnyEvents()
	}
	door.Lock()
}

func contactCheck(check chan bool) {
	for door.Open() {
		time.Sleep(500 * time.Millisecond)
	}
	log.Println("MAIN: Detected door close")
	check <- true
}

func resetToIdle() *time.Timer {
	resetToIdle := time.NewTimer(resetToIdleDuration)
	go func() {
		<-resetToIdle.C
		door.Lock()
	}()
	log.Println("MAIN: Scheduled resetToIdle")
	return resetToIdle
}

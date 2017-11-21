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
	openWaitDuration       = 60 * time.Second
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

	door.CheckSunRise(
		os.Getenv("LATITUDE"),
		os.Getenv("LONGITUDE"),
		os.Getenv("DAY_START"),
		os.Getenv("DAY_END"),
		door.SetDarkOutside)

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
			valid, name := codebook.Check(code.Digits, time.Now().UTC())
			if valid {
				validCode(code.Digits, name)
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

func validCode(digits, name string) {
	log.Printf("MAIN: Unlocked with code: %v\n", digits)
	door.Unlock()
	sms.SendCorrectCode(digits, name)
	scheduleEvent(waitForDoorToBeOpened())
}

func overrideOpen(overrideType string) {
	scheduleEvent(waitForDoorToBeOpened())
	log.Printf("MAIN: Unlocked with override: %v\n", overrideType)
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

func waitForDoorToBeOpened() *time.Timer {
	resetToLocked := time.NewTimer(openWaitDuration)
	go checkDoor(resetToLocked, door.Open)
	log.Println("MAIN: Scheduled resetToLocked (check=open)")
	return resetToLocked
}

func resetToLocked() *time.Timer {
	resetToLocked := time.NewTimer(doorOpenDuration)
	go checkDoor(resetToLocked, door.Closed)
	log.Println("MAIN: Scheduled resetToLocked (check=closed)")
	return resetToLocked
}

func checkDoor(checkDuration *time.Timer, expectedState door.ContactState) {
	contact := make(chan bool)
	stop := make(chan bool)
	defer close(contact)
	defer close(stop)

	go contactCheck(contact, stop, expectedState)

	select {
	case <-checkDuration.C:
		log.Println("MAIN: checkDuration returned")
		switch expectedState {
		case door.Closed:
			if door.State() == door.Open {
				sms.SendDoorNotClosed()
				log.Println("MAIN: Door not closed")
			}
		case door.Open:
			if door.State() == door.Closed {
				sms.SendDoorNotOpened()
				log.Println("MAIN: Door never opened")
			}
		}
	case <-contact:
		unscheduleAnyEvents()
	}
	if expectedState == door.Open {
		resetToLocked()
	} else {
		door.Lock()
	}
}

func contactCheck(check chan bool, stop chan bool, expectedState door.ContactState) {
	stopped := false
	for door.State() != expectedState && !stopped {
		select {
		default:
			time.Sleep(100 * time.Millisecond)
		case <-stop:
			stopped = true
			return
		}
	}
	if expectedState == door.Open {
		log.Println("MAIN: Detected door open")
	} else {
		log.Println("MAIN: Detected door close")
	}
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

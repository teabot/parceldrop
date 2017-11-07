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

package keypad

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	shiftCode     = 42
	clearCode     = 9
	enterCode     = 4
	eventCaptures = 16
	// DefaultDevice location on the Pi
	DefaultDevice = "/dev/input/event0"
)

// InputEvent is a Go implementation of the native linux device
// input_event struct, as described in the kernel documentation
// (https://www.kernel.org/doc/Documentation/input/input.txt),
// with a big assist from https://github.com/gvalkov/golang-evdev
type inputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

type keyEvent struct {
	Digit   string
	Enter   bool
	Clear   bool
	Timeout bool
}

var eventSize = int(unsafe.Sizeof(inputEvent{}))

// KEYCODES is the map of hex found in the InputEvent.Code field, and
// its corresponding char (string) representation
// [source: Vojtech Pavlik (author of the Linux Input Drivers project),
// via linuxquestions.org user bricedebrignaisplage]
var KEYCODES = map[byte]string{
	0x02: "1",
	0x03: "2",
	0x04: "3",
	0x05: "4",
	0x06: "5",
	0x07: "6",
	0x08: "7",
	0x09: "8",
	0x0a: "9",
	0x0b: "0",
}

// lookupKeyCode finds the corresponding string for the given hex byte,
// returning "-" as the default if not found
func lookupKeyCode(b byte) string {
	val, exists := KEYCODES[b]
	if exists {
		return val
	}
	fmt.Printf("unknown: %x\n", b)
	return "-"
}

// read takes the open scanner device pointer and returns a list of
// inputEvent captures, corresponding to input (scan) events
func read(dev *os.File) ([]inputEvent, error) {
	events := make([]inputEvent, eventCaptures)
	buffer := make([]byte, eventSize*eventCaptures)
	_, err := dev.Read(buffer)
	if err != nil {
		return events, err
	}
	b := bytes.NewBuffer(buffer)
	err = binary.Read(b, binary.LittleEndian, &events)
	if err != nil {
		return events, err
	}
	// remove trailing structures
	for i := range events {
		if events[i].Time.Sec == 0 {
			events = append(events[:i])
			break
		}
	}
	return events, err
}

// ScanKeys takes a linux input device string pointing to the scanner
// to read from, invokes the given function on the resulting barcode string
// when complete, or the errfn on error, then goes back to read/scan again
func ScanKeys(device string, keys chan<- keyEvent, errFn func(error)) {
	scanner, err := os.Open(device)
	if err != nil {
		// invoke the function which handles scanner errors
		errFn(err)
	}
	defer scanner.Close()

	shift := false
	for {
		events, scanErr := read(scanner)
		if scanErr != nil {
			// invoke the function which handles scanner errors
			errFn(scanErr)
		}
		for i := range events {
			if events[i].Type == 1 && events[i].Value == 1 {
				if events[i].Code == shiftCode {
					shift = true
				} else {
					if shift {
						if events[i].Code == enterCode {
							keys <- keyEvent{
								Enter: true,
							}
						} else if events[i].Code == clearCode {
							keys <- keyEvent{
								Clear: true,
							}
						}
					} else {
						if events[i].Code != 0 {
							keys <- keyEvent{
								Digit: lookupKeyCode(byte(events[i].Code)),
							}
						}
					}
					shift = false
				}
			}
		}
	}
}

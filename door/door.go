package door

import (
	"fmt"
	"log"
	"time"

	"github.com/luismesas/goPi/piface"
	"github.com/luismesas/goPi/spi"
)

type inputPin int
type outputPin int
type ContactState bool

const (
	red    outputPin = 4
	green  outputPin = 5
	blue   outputPin = 6
	white  outputPin = 7
	latch  outputPin = 3
	lights outputPin = 2

	contact  inputPin = 0
	override inputPin = 1

	Open   ContactState = true
	Closed ContactState = false

	// when door is open
	contactOpen = 0

	// when door is closed
	manualOverrideEngaged = 0
)

var locked = true
var pfd *piface.PiFaceDigital
var openFn func(string)
var darkOutsideInHours = false
var darkOutside = false

// Initialise x
func Initialise(overrideFn func(string)) error {
	openFn = overrideFn
	pfd = piface.NewPiFaceDigital(spi.DEFAULT_HARDWARE_ADDR, spi.DEFAULT_BUS, spi.DEFAULT_CHIP)
	err := pfd.InitBoard()
	if err != nil {
		fmt.Printf("DOOR: Error on init board: %s", err)
		return err
	}
	Lock()
	go checkOverride()
	return nil
}

func State() ContactState {
	if pfd.Switches[contact].Value() == contactOpen {
		return Open
	}
	return Closed
}

// Locked x
func Locked() bool {
	return locked
}

// Unlock x
func Unlock() {
	log.Println("DOOR: Latch activated")
	log.Println("DOOR: LED: Green")
	pfd.Leds[red].AllOff()
	pfd.Leds[blue].AllOff()
	pfd.Leds[white].AllOff()

	pfd.Leds[green].AllOn()
	if darkOutside {
		pfd.Leds[lights].AllOn()
	}

	pfd.Leds[latch].AllOn()
	time.Sleep(1 * time.Second)
	pfd.Leds[latch].AllOff()
	locked = false
	return
}

// Reject x
func Reject() {
	log.Println("DOOR: LED: Red")
	pfd.Leds[green].AllOff()
	pfd.Leds[blue].AllOff()
	pfd.Leds[white].AllOff()
	pfd.Leds[latch].AllOff()

	pfd.Leds[red].AllOn()
	return
}

// Wait x
func Wait() {
	log.Println("DOOR: LED: Blue")
	pfd.Leds[red].AllOff()
	pfd.Leds[green].AllOff()
	// pfd.Leds[white].AllOff()
	pfd.Leds[latch].AllOff()

	pfd.Leds[blue].AllOn()
	return
}

// Lock x
func Lock() {
	log.Println("DOOR: Locked")
	pfd.Leds[red].AllOff()
	pfd.Leds[green].AllOff()
	pfd.Leds[blue].AllOff()
	pfd.Leds[latch].AllOff()
	pfd.Leds[lights].AllOff()

	resetToLight()
	locked = true
	return
}

func SetDarkOutside(darkInHours, dark bool) {
	darkOutsideInHours = darkInHours
	darkOutside = dark
	if locked {
		resetToLight()
	}
}

func resetToLight() {
	if darkOutsideInHours {
		//log.Println("DOOR: LED: White")
		pfd.Leds[white].AllOn()
	} else {
		// log.Println("DOOR: LED: Off")
		pfd.Leds[white].AllOff()
	}
}

// This does not work properly
func checkOverride() {
	for {
		if pfd.Switches[override].Value() == manualOverrideEngaged && locked {
			log.Println("DOOR: Manual override, unlocking")
			openFn("button")
		}
		time.Sleep(200 * time.Millisecond)
	}
}

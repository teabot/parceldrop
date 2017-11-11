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

const (
	red   outputPin = 4
	green outputPin = 5
	blue  outputPin = 6
	white outputPin = 7
	latch outputPin = 3

	contact  inputPin = 0
	override inputPin = 1

	open   = 1
	closed = 0
)

var locked = true
var pfd *piface.PiFaceDigital

// Initialise x
func Initialise() error {
	pfd = piface.NewPiFaceDigital(spi.DEFAULT_HARDWARE_ADDR, spi.DEFAULT_BUS, spi.DEFAULT_CHIP)
	err := pfd.InitBoard()
	if err != nil {
		fmt.Printf("Error on init board: %s", err)
		return err
	}
	Lock()
	go checkOverride()
	return nil
}

// Open x
func Open() bool {
	return pfd.Switches[contact].Value() == open
}

// Closed x
func Closed() bool {
	return !Open()
}

// Locked x
func Locked() bool {
	return locked
}

// Unlock x
func Unlock() {
	log.Println("Door: latch")
	log.Println("LED: Green")
	pfd.Leds[red].AllOff()
	pfd.Leds[blue].AllOff()
	pfd.Leds[white].AllOff()

	pfd.Leds[green].AllOn()
	pfd.Leds[latch].AllOn()
	time.Sleep(1 * time.Second)
	pfd.Leds[latch].AllOff()
	locked = false
	return
}

// Reject x
func Reject() {
	log.Println("LED: Red")
	pfd.Leds[green].AllOff()
	pfd.Leds[blue].AllOff()
	pfd.Leds[white].AllOff()
	pfd.Leds[latch].AllOff()

	pfd.Leds[red].AllOn()
	return
}

// Wait x
func Wait() {
	log.Println("LED: Blue")
	pfd.Leds[red].AllOff()
	pfd.Leds[green].AllOff()
	pfd.Leds[white].AllOff()
	pfd.Leds[latch].AllOff()

	pfd.Leds[blue].AllOn()
	return
}

// Lock x
func Lock() {
	log.Println("LED: White")
	pfd.Leds[red].AllOff()
	pfd.Leds[green].AllOff()
	pfd.Leds[blue].AllOff()
	pfd.Leds[latch].AllOff()

	pfd.Leds[white].AllOn()
	locked = true
	return
}

func checkOverride() {
	for {
		if pfd.Switches[override].Value() == closed && locked {
			log.Println("Manual override, unlocking")
			Unlock()
		}
		time.Sleep(200 * time.Millisecond)
	}
}

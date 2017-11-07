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
	red   outputPin = 3
	green outputPin = 4
	blue  outputPin = 5
	white outputPin = 6
	latch outputPin = 7

	contact inputPin = 1
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
	return nil
}

// Open x
func Open() bool {
	// return pfd.Switches[contact].Value() != 0
	return false
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

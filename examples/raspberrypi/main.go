package main

import (
	"fmt"
	"log"
	"os"

	"github.com/adafruit/ILI9341_GoLang"
	"periph.io/x/conn/v3/driver"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

const (
	// ILI9341 colors
	ILI9341_BLACK   = 0x0000
	ILI9341_RED     = 0xF800
)

// PeriphPin is an implementation of the Pin interface using periph.io.
type PeriphPin struct {
	pin gpio.PinIO
}

// Set sets the pin state.
func (p *PeriphPin) Set(high bool) {
	if high {
		p.pin.Out(gpio.High)
	} else {
		p.pin.Out(gpio.Low)
	}
}

func main() {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Use spireg SPI port default.
	p, err := spireg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	// Set the SPI frequency.
	if err := p.LimitSpeed(40 * 1000 * 1000); err != nil {
		log.Fatal(err)
	}

	// Set the SPI mode.
	if err := p.SetMode(spi.Mode0); err != nil {
		log.Fatal(err)
	}

	// Get the GPIO pins for DC and RST.
	dcPin := gpioreg.ByName("GPIO24")
	if dcPin == nil {
		log.Fatal("Failed to find GPIO24")
	}

	rstPin := gpioreg.ByName("GPIO25")
	if rstPin == nil {
		log.Fatal("Failed to find GPIO25")
	}

	// Create the ILI9341 device.
	dev, err := ili9341.New(p, &PeriphPin{pin: dcPin}, &PeriphPin{pin: rstPin})
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the display.
	if err := dev.Begin(); err != nil {
		log.Fatal(err)
	}

	// Fill the screen with black.
	if err := dev.FillRectangle(0, 0, ili9341.TFTWIDTH, ili9341.TFTHEIGHT, ILI9341_BLACK); err != nil {
		log.Fatal(err)
	}

	// Draw a red rectangle.
	if err := dev.FillRectangle(10, 20, 50, 30, ILI9341_RED); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Drew a red rectangle on the screen.")
}

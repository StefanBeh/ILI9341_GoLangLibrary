package main

import (
	"fmt"
	"log"

	"github.com/behling_stefan/ILI9341_GoLangLibrary/ili9341_golangLibrary"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

const (
	// ILI9341 colors
	ILI9341_BLACK = 0x0000
	ILI9341_RED   = 0xF800
	ILI9341_WHITE = 0xFFFF
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

// PeriphSPI is an implementation of the SPI interface using periph.io.
type PeriphSPI struct {
	conn spi.Conn
}

// Tx sends and receives data.
func (s *PeriphSPI) Tx(w, r []byte) error {
	return s.conn.Tx(w, r)
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

	// Connect to the SPI device.
	conn, err := p.Connect(10*physic.MegaHertz, spi.Mode0, 8)
	if err != nil {
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
	dev, err := ili9341.New(&PeriphSPI{conn: conn}, &PeriphPin{pin: dcPin}, &PeriphPin{pin: rstPin})
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the display.
	if err := dev.Begin(); err != nil {
		log.Fatal(err)
	}

	// Set rotation to landscape.
	if err := dev.SetRotation(1); err != nil {
		log.Fatal(err)
	}

	// Fill the screen with black.
	if err := dev.FillRectangle(0, 0, ili9341.TFTHEIGHT, ili9341.TFTWIDTH, ILI9341_BLACK); err != nil {
		log.Fatal(err)
	}

	// Draw a table
	// horizontal lines
	if err := dev.DrawLine(0, 40, 320, 40, ILI9341_WHITE); err != nil {
		log.Fatal(err)
	}
	if err := dev.DrawLine(0, 80, 320, 80, ILI9341_WHITE); err != nil {
		log.Fatal(err)
	}
	if err := dev.DrawLine(0, 120, 320, 120, ILI9341_WHITE); err != nil {
		log.Fatal(err)
	}
	if err := dev.DrawLine(0, 160, 320, 160, ILI9341_WHITE); err != nil {
		log.Fatal(err)
	}
	if err := dev.DrawLine(0, 200, 320, 200, ILI9341_WHITE); err != nil {
		log.Fatal(err)
	}

	// vertical lines
	if err := dev.DrawLine(100, 0, 100, 240, ILI9341_WHITE); err != nil {
		log.Fatal(err)
	}
	if err := dev.DrawLine(220, 0, 220, 240, ILI9341_WHITE); err != nil {
		log.Fatal(err)
	}

	// Draw text
	if err := dev.DrawString(10, 10, "Header 1", ILI9341_WHITE, ILI9341_BLACK, 2); err != nil {
		log.Fatal(err)
	}
	if err := dev.DrawString(110, 10, "Header 2", ILI9341_WHITE, ILI9341_BLACK, 2); err != nil {
		log.Fatal(err)
	}
	if err := dev.DrawString(230, 10, "Header 3", ILI9341_WHITE, ILI9341_BLACK, 2); err != nil {
		log.Fatal(err)
	}

	if err := dev.DrawString(10, 50, "Row 1, Col 1", ILI9341_WHITE, ILI9341_BLACK, 1); err != nil {
		log.Fatal(err)
	}
	if err := dev.DrawString(110, 50, "Row 1, Col 2", ILI9341_WHITE, ILI9341_BLACK, 1); err != nil {
		log.Fatal(err)
	}
	if err := dev.DrawString(230, 50, "Row 1, Col 3", ILI9341_WHITE, ILI9341_BLACK, 1); err != nil {
		log.Fatal(err)
	}

	if err := dev.DrawString(10, 90, "Row 2, Col 1", ILI9341_WHITE, ILI9341_BLACK, 1); err != nil {
		log.Fatal(err)
	}
	if err := dev.DrawString(110, 90, "Row 2, Col 2", ILI9341_WHITE, ILI9341_BLACK, 1); err != nil {
		log.Fatal(err)
	}
	if err := dev.DrawString(230, 90, "Row 2, Col 3", ILI9341_WHITE, ILI9341_BLACK, 1); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Drew a table on the screen.")
}

# ILI9341 GoLang Library

This repository contains a GoLang library for controlling ILI9341 TFT displays. It is a rewrite and adaptation of the popular Adafruit ILI9341 Python library, tailored for Go applications, particularly on single-board computers like the Raspberry Pi.

## Features

- Basic display initialization and control.
- Drawing primitives (points, lines, rectangles, circles, etc.).
- Text rendering.
- Image display.

## Installation

To use this library, you'll need to have Go installed on your system.

```bash
go get github.com/behling_stefan/ILI9341_GoLangLibrary/ili9341_golangLibrary
```

## Wiring

For detailed wiring instructions to connect your ILI9341 display to a Raspberry Pi, please refer to the [wiring guide](wiring.md).

## Usage

An example of how to use this library on a Raspberry Pi can be found in the `examples/raspberrypi` directory.

```go
// See examples/raspberrypi/main.go for a complete example
package main

import (
	"log"
	"time"

	"github.com/behling_stefan/ILI9341_GoLangLibrary/ili9341_golangLibrary"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

func main() {
	// Initialize periph.io host
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Open SPI port
	p, err := spireg.Open("/dev/spidev0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	// Configure SPI
	conn, err := p.Connect(physic.MegaHertz, spi.Mode3, 8)
	if err != nil {
		log.Fatal(err)
	}

	// Create ILI9341 driver
	dcPin := "GPIO24" // Data/Command pin
	rstPin := "GPIO25" // Reset pin
	display, err := ili9341_golangLibrary.NewILI9341(conn, dcPin, rstPin)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize display
	display.Init()
	display.SetRotation(ili9341_golangLibrary.Rotation90)
	display.FillScreen(ili9341_golangLibrary.Black)

	// Draw something
	display.DrawString(10, 10, "Hello, Go!", ili9341_golangLibrary.White, ili9341_golangLibrary.Font8x8)
	display.DrawRect(5, 5, 100, 50, ili9341_golangLibrary.Red)

	time.Sleep(5 * time.Second)
}
```

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## License

Adafruit invests time and resources providing this open source code,
please support Adafruit and open-source hardware by purchasing
products from Adafruit!

Written by Limor Fried/Ladyada for Adafruit Industries.
MIT license, all text above must be included in any redistribution
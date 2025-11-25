module raspberrypi-example

go 1.25.4

replace github.com/adafruit/ILI9341_GoLang => ../../ili9341

require (
	github.com/adafruit/ILI9341_GoLang v0.0.0-00010101000000-000000000000
	periph.io/x/conn/v3 v3.7.2
	periph.io/x/host/v3 v3.8.5
)

module raspberrypi-example

go 1.25.4

replace github.com/behling_stefan/ILI9341_GoLangLibrary => ../..

require (
	github.com/behling_stefan/ILI9341_GoLangLibrary v0.0.0
	periph.io/x/conn/v3 v3.7.2
	periph.io/x/host/v3 v3.8.5
)

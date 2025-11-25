package ili9341

import (
	"time"
)

const (
	TFTWIDTH  = 240
	TFTHEIGHT = 320

	NOP     = 0x00
	SWRESET = 0x01
	RDDID   = 0x04
	RDDST   = 0x09

	SLPIN   = 0x10
	SLPOUT  = 0x11
	PTLON   = 0x12
	NORON   = 0x13

	RDMODE     = 0x0A
	RDMADCTL   = 0x0B
	RDPIXFMT   = 0x0C
	RDIMGFMT   = 0x0D
	RDSELFDIAG = 0x0F

	INVOFF   = 0x20
	INVON    = 0x21
	GAMMASET = 0x26
	DISPOFF  = 0x28
	DISPON   = 0x29

	CASET = 0x2A
	PASET = 0x2B
	RAMWR = 0x2C
	RAMRD = 0x2E

	PTLAR    = 0x30
	VSCRDEF  = 0x33
	MADCTL   = 0x36
	VSCRSADD = 0x37
	PIXFMT   = 0x3A

	FRMCTR1 = 0xB1
	FRMCTR2 = 0xB2
	FRMCTR3 = 0xB3
	INVCTR  = 0xB4
	DFUNCTR = 0xB6

	PWCTR1 = 0xC0
	PWCTR2 = 0xC1
	PWCTR3 = 0xC2
	PWCTR4 = 0xC3
	PWCTR5 = 0xC4
	VMCTR1 = 0xC5
	VMCTR2 = 0xC7

	RDID1 = 0xDA
	RDID2 = 0xDB
	RDID3 = 0xDC
	RDID4 = 0xDD

	GMCTRP1 = 0xE0
	GMCTRN1 = 0xE1

	MADCTL_MY  = 0x80
	MADCTL_MX  = 0x40
	MADCTL_MV  = 0x20
	MADCTL_ML  = 0x10
	MADCTL_BGR = 0x08
	MADCTL_MH  = 0x04
)

var initcmd = []byte{
	0xEF, 3, 0x03, 0x80, 0x02,
	0xCF, 3, 0x00, 0xC1, 0x30,
	0xED, 4, 0x64, 0x03, 0x12, 0x81,
	0xE8, 3, 0x85, 0x00, 0x78,
	0xCB, 5, 0x39, 0x2C, 0x00, 0x34, 0x02,
	0xF7, 1, 0x20,
	0xEA, 2, 0x00, 0x00,
	PWCTR1, 1, 0x23, // Power control VRH[5:0]
	PWCTR2, 1, 0x10, // Power control SAP[2:0];BT[3:0]
	VMCTR1, 2, 0x3e, 0x28, // VCM control
	VMCTR2, 1, 0x86, // VCM control2
	MADCTL, 1, 0x48, // Memory Access Control
	VSCRSADD, 1, 0x00, // Vertical scroll zero
	PIXFMT, 1, 0x55,
	FRMCTR1, 2, 0x00, 0x18,
	DFUNCTR, 3, 0x08, 0x82, 0x27, // Display Function Control
	0xF2, 1, 0x00, // 3Gamma Function Disable
	GAMMASET, 1, 0x01, // Gamma curve selected
	GMCTRP1, 15, 0x0F, 0x31, 0x2B, 0x0C, 0x0E, 0x08, // Set Gamma
	0x4E, 0xF1, 0x37, 0x07, 0x10, 0x03, 0x0E, 0x09, 0x00,
	GMCTRN1, 15, 0x00, 0x0E, 0x14, 0x03, 0x11, 0x07, // Set Gamma
	0x31, 0xC1, 0x48, 0x08, 0x0F, 0x0C, 0x31, 0x36, 0x0F,
	SLPOUT, 0x80, // Exit Sleep
	DISPON, 0x80, // Display on
	0x00, // End of list
}

// SPI is an interface for a SPI bus.
type SPI interface {
	// Tx sends and receives data.
	Tx(w, r []byte) error
}

// Pin is an interface for a GPIO pin.
type Pin interface {
	// Set sets the pin to high or low.
	Set(high bool)
}

// ILI9341 is a handle to the display controller.
type ILI9341 struct {
	spi      SPI
	dc       Pin
	rst      Pin
	width    int
	height   int
	rotation uint8
}

// New creates a new ILI9341 controller.
func New(spi SPI, dc, rst Pin) (*ILI9341, error) {
	return &ILI9341{
		spi:    spi,
		dc:     dc,
		rst:    rst,
		width:  TFTWIDTH,
		height: TFTHEIGHT,
	}, nil
}

// Begin initializes the display.
func (d *ILI9341) Begin() error {
	if d.rst != nil {
		d.rst.Set(true)
		time.Sleep(5 * time.Millisecond)
		d.rst.Set(false)
		time.Sleep(20 * time.Millisecond)
		d.rst.Set(true)
		time.Sleep(150 * time.Millisecond)
	} else {
		// Perform a software reset
		if err := d.sendCommand(SWRESET); err != nil {
			return err
		}
		time.Sleep(150 * time.Millisecond)
	}

	addr := 0
	for {
		cmd := initcmd[addr]
		addr++
		if cmd == 0 {
			break
		}
		x := initcmd[addr]
		addr++
		numArgs := int(x & 0x7F)
		if err := d.sendCommand(cmd, initcmd[addr:addr+numArgs]...); err != nil {
			return err
		}
		addr += numArgs
		if (x & 0x80) != 0 {
			time.Sleep(150 * time.Millisecond)
		}
	}
	d.width = TFTWIDTH
	d.height = TFTHEIGHT

	return nil
}

func (d *ILI9341) sendCommand(cmd byte, data ...byte) error {
	d.dc.Set(false) // command mode
	if err := d.spi.Tx([]byte{cmd}, nil); err != nil {
		return err
	}
	if len(data) > 0 {
		d.dc.Set(true) // data mode
		if err := d.spi.Tx(data, nil); err != nil {
			return err
		}
	}
	return nil
}

// SetRotation sets the rotation of the display.
func (d *ILI9341) SetRotation(r uint8) error {
	d.rotation = r % 4 // can't be higher than 3
	var m byte
	switch d.rotation {
	case 0:
		m = MADCTL_MX | MADCTL_BGR
		d.width = TFTWIDTH
		d.height = TFTHEIGHT
	case 1:
		m = MADCTL_MV | MADCTL_BGR
		d.width = TFTHEIGHT
		d.height = TFTWIDTH
	case 2:
		m = MADCTL_MY | MADCTL_BGR
		d.width = TFTWIDTH
		d.height = TFTHEIGHT
	case 3:
		m = MADCTL_MX | MADCTL_MY | MADCTL_MV | MADCTL_BGR
		d.width = TFTHEIGHT
		d.height = TFTWIDTH
	}
	return d.sendCommand(MADCTL, m)
}

// InvertDisplay enables or disables display color inversion.
func (d *ILI9341) InvertDisplay(invert bool) error {
	if invert {
		return d.sendCommand(INVON)
	}
	return d.sendCommand(INVOFF)
}

// ScrollTo scrolls the display by a given amount.
func (d *ILI9341) ScrollTo(y uint16) error {
	data := []byte{byte(y >> 8), byte(y & 0xff)}
	return d.sendCommand(VSCRSADD, data...)
}

// SetScrollMargins sets the top and bottom scroll margins.
func (d *ILI9341) SetScrollMargins(top, bottom uint16) error {
	if top+bottom <= TFTHEIGHT {
		middle := TFTHEIGHT - (top + bottom)
		data := []byte{
			byte(top >> 8), byte(top & 0xff),
			byte(middle >> 8), byte(middle & 0xff),
			byte(bottom >> 8), byte(bottom & 0xff),
		}
		return d.sendCommand(VSCRDEF, data...)
	}
	return nil
}

// SetAddrWindow sets the address window for RAM access.
func (d *ILI9341) SetAddrWindow(x, y, w, h uint16) error {
	x2 := x + w - 1
	y2 := y + h - 1
	data := []byte{byte(x >> 8), byte(x & 0xff), byte(x2 >> 8), byte(x2 & 0xff)}
	if err := d.sendCommand(CASET, data...); err != nil {
		return err
	}
	data = []byte{byte(y >> 8), byte(y & 0xff), byte(y2 >> 8), byte(y2 & 0xff)}
	if err := d.sendCommand(PASET, data...); err != nil {
		return err
	}
	return d.sendCommand(RAMWR)
}

// FillRectangle fills a rectangle with a given color.
func (d *ILI9341) FillRectangle(x, y, w, h uint16, color uint16) error {
	if err := d.SetAddrWindow(x, y, w, h); err != nil {
		return err
	}
	d.dc.Set(true)
	c := []byte{byte(color >> 8), byte(color)}
	for i := 0; i < int(w*h); i++ {
		if err := d.spi.Tx(c, nil); err != nil {
			return err
		}
	}
	return nil
}

// DrawPixel draws a single pixel at a given position and color.
func (d *ILI9341) DrawPixel(x, y int16, color uint16) error {
	if x < 0 || x >= int16(d.width) || y < 0 || y >= int16(d.height) {
		return nil
	}
	if err := d.SetAddrWindow(uint16(x), uint16(y), 1, 1); err != nil {
		return err
	}
	d.dc.Set(true)
	c := []byte{byte(color >> 8), byte(color)}
	return d.spi.Tx(c, nil)
}

// DrawLine draws a line between two points.
func (d *ILI9341) DrawLine(x0, y0, x1, y1 int16, color uint16) error {
	steep := abs(y1-y0) > abs(x1-x0)
	if steep {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
	}
	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}
	dx := x1 - x0
	dy := abs(y1 - y0)
	err := int16(dx / 2)
	ystep := int16(1)
	if y0 > y1 {
		ystep = -1
	}
	for ; x0 <= x1; x0++ {
		if steep {
			if err := d.DrawPixel(y0, x0, color); err != nil {
				return err
			}
		} else {
			if err := d.DrawPixel(x0, y0, color); err != nil {
				return err
			}
		}
		err -= dy
		if err < 0 {
			y0 += ystep
			err += dx
		}
	}
	return nil
}

func abs(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}

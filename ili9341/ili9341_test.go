package ili9341

import (
	"testing"
)

func TestBegin(t *testing.T) {
	spi := &MockSPI{}
	dc := &MockPin{}
	rst := &MockPin{}

	dev, err := New(spi, dc, rst)
	if err != nil {
		t.Fatal(err)
	}

	if err := dev.Begin(); err != nil {
		t.Fatal(err)
	}

	// In a real test, we would check the SPI data to make sure the correct
	// commands were sent. For now, we just check that the test runs without
	// error.
}

func TestSetRotation(t *testing.T) {
	spi := &MockSPI{}
	dc := &MockPin{}
	rst := &MockPin{}

	dev, err := New(spi, dc, rst)
	if err != nil {
		t.Fatal(err)
	}

	if err := dev.SetRotation(1); err != nil {
		t.Fatal(err)
	}

	expected := []byte{MADCTL, MADCTL_MV | MADCTL_BGR}
	if len(spi.W) != len(expected) {
		t.Fatalf("unexpected spi data length: got %d, want %d", len(spi.W), len(expected))
	}
	for i := range expected {
		if spi.W[i] != expected[i] {
			t.Errorf("unexpected spi data at index %d: got %x, want %x", i, spi.W[i], expected[i])
		}
	}
}

func TestInvertDisplay(t *testing.T) {
	spi := &MockSPI{}
	dc := &MockPin{}
	rst := &MockPin{}

	dev, err := New(spi, dc, rst)
	if err != nil {
		t.Fatal(err)
	}

	if err := dev.InvertDisplay(true); err != nil {
		t.Fatal(err)
	}

	expected := []byte{INVON}
	if len(spi.W) != len(expected) {
		t.Fatalf("unexpected spi data length: got %d, want %d", len(spi.W), len(expected))
	}
	for i := range expected {
		if spi.W[i] != expected[i] {
			t.Errorf("unexpected spi data at index %d: got %x, want %x", i, spi.W[i], expected[i])
		}
	}

	spi.W = nil

	if err := dev.InvertDisplay(false); err != nil {
		t.Fatal(err)
	}

	expected = []byte{INVOFF}
	if len(spi.W) != len(expected) {
		t.Fatalf("unexpected spi data length: got %d, want %d", len(spi.W), len(expected))
	}
}

func TestScrollTo(t *testing.T) {
	spi := &MockSPI{}
	dc := &MockPin{}
	rst := &MockPin{}

	dev, err := New(spi, dc, rst)
	if err != nil {
		t.Fatal(err)
	}

	if err := dev.ScrollTo(100); err != nil {
		t.Fatal(err)
	}

	expected := []byte{VSCRSADD, 0, 100}
	if len(spi.W) != len(expected) {
		t.Fatalf("unexpected spi data length: got %d, want %d", len(spi.W), len(expected))
	}
	for i := range expected {
		if spi.W[i] != expected[i] {
			t.Errorf("unexpected spi data at index %d: got %x, want %x", i, spi.W[i], expected[i])
		}
	}
}

func TestSetScrollMargins(t *testing.T) {
	spi := &MockSPI{}
	dc := &MockPin{}
	rst := &MockPin{}

	dev, err := New(spi, dc, rst)
	if err != nil {
		t.Fatal(err)
	}

	if err := dev.SetScrollMargins(10, 20); err != nil {
		t.Fatal(err)
	}

	middle := TFTHEIGHT - (10 + 20)
	expected := []byte{
		VSCRDEF,
		0, 10,
		byte(middle >> 8), byte(middle & 0xff),
		0, 20,
	}
	if len(spi.W) != len(expected) {
		t.Fatalf("unexpected spi data length: got %d, want %d", len(spi.W), len(expected))
	}
	for i := range expected {
		if spi.W[i] != expected[i] {
			t.Errorf("unexpected spi data at index %d: got %x, want %x", i, spi.W[i], expected[i])
		}
	}
}

func TestSetAddrWindow(t *testing.T) {
	spi := &MockSPI{}
	dc := &MockPin{}
	rst := &MockPin{}

	dev, err := New(spi, dc, rst)
	if err != nil {
		t.Fatal(err)
	}

	if err := dev.SetAddrWindow(10, 20, 30, 40); err != nil {
		t.Fatal(err)
	}

	x2 := 10 + 30 - 1
	y2 := 20 + 40 - 1
	expected := []byte{
		CASET, 0, 10, 0, byte(x2),
		PASET, 0, 20, 0, byte(y2),
		RAMWR,
	}
	if len(spi.W) != len(expected) {
		t.Fatalf("unexpected spi data length: got %d, want %d", len(spi.W), len(expected))
	}
	for i := range expected {
		if spi.W[i] != expected[i] {
			t.Errorf("unexpected spi data at index %d: got %x, want %x", i, spi.W[i], expected[i])
		}
	}
}

func TestFillRectangle(t *testing.T) {
	spi := &MockSPI{}
	dc := &MockPin{}
	rst := &MockPin{}

	dev, err := New(spi, dc, rst)
	if err != nil {
		t.Fatal(err)
	}

	if err := dev.FillRectangle(10, 20, 30, 40, 0x1234); err != nil {
		t.Fatal(err)
	}

	x2 := 10 + 30 - 1
	y2 := 20 + 40 - 1
	expected := []byte{
		CASET, 0, 10, 0, byte(x2),
		PASET, 0, 20, 0, byte(y2),
		RAMWR,
	}
	// 30*40 = 1200 pixels, 2 bytes per pixel
	for i := 0; i < 1200; i++ {
		expected = append(expected, 0x12, 0x34)
	}

	if len(spi.W) != len(expected) {
		t.Fatalf("unexpected spi data length: got %d, want %d", len(spi.W), len(expected))
	}

	for i := 0; i < 11; i++ {
		if spi.W[i] != expected[i] {
			t.Errorf("unexpected spi data at index %d: got %x, want %x", i, spi.W[i], expected[i])
		}
	}
	for i := 11; i < len(expected); i += 2 {
		if spi.W[i] != 0x12 || spi.W[i+1] != 0x34 {
			t.Errorf("unexpected spi data at index %d: got %x%x, want %x%x", i, spi.W[i], spi.W[i+1], 0x12, 0x34)
		}
	}
}

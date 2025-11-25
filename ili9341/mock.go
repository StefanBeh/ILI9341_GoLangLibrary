package ili9341

// MockPin is a mock implementation of the Pin interface.
type MockPin struct {
	State bool
}

// Set sets the pin state.
func (p *MockPin) Set(high bool) {
	p.State = high
}

// MockSPI is a mock implementation of the SPI interface.
type MockSPI struct {
	W []byte
	R []byte
}

// Tx is the mock implementation of the Tx method.
func (s *MockSPI) Tx(w, r []byte) error {
	s.W = append(s.W, w...)
	if r != nil {
		// In a real implementation, we would read from the SPI bus.
		// For now, we just fill the read buffer with zeros.
		for i := range r {
			r[i] = 0
		}
	}
	return nil
}

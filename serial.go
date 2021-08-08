package serial

import (
	"fmt"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

// I2C defines an I2C type to read and write from/to registers.
type I2C struct {
	dev  *i2c.Dev
	bus  i2c.BusCloser
	addr uint16
}

// NewI2C returns a new I2C interface at the specified bus and address.
// If `bus` is set to "", the first available bus is used. The address must
// always be specified.
func NewI2C(bus string, addr uint16) (*I2C, error) {
	if _, err := host.Init(); err != nil {
		return nil, fmt.Errorf("serial: could not initialize host: %w", err)
	}

	b, err := i2creg.Open(bus)
	if err != nil {
		return nil, fmt.Errorf("serial: could not open I2C bus: %w", err)
	}

	dev := &i2c.Dev{
		Addr: addr,
		Bus:  b,
	}

	i2c := &I2C{
		dev:  dev,
		bus:  b,
		addr: addr,
	}

	return i2c, nil
}

// Read reads a single byte from a register.
func (i *I2C) Read(reg byte) (byte, error) {
	b := make([]byte, 1)
	if err := i.dev.Tx([]byte{reg}, b); err != nil {
		return 0, fmt.Errorf("serial: could not read byte from register %x at address %x: %w", reg, i.addr, err)
	}

	return b[0], nil
}

// ReadBytes reads n bytes from a register.
func (i *I2C) ReadBytes(reg byte, n int) ([]byte, error) {
	b := make([]byte, n)
	if err := i.dev.Tx([]byte{reg}, b); err != nil {
		return nil, fmt.Errorf("serial: could not read all %d bytes from register %x at address %x: %w", n, reg, i.addr, err)
	}

	return b, nil
}

// Write writes a byte or bytes to a register.
func (i *I2C) Write(reg byte, data ...byte) error {
	n, err := i.dev.Write(append([]byte{reg}, data...))
	if err != nil {
		return fmt.Errorf("serial: could not write %x to register %x at address %x: %w", data, reg, i.addr, err)
	}
	n-- // remove register write
	if n != len(data) {
		return fmt.Errorf("serial: wrong number of bytes written: want %d, got %d", len(data), n)
	}

	return nil
}

// Close closes the bus used by I2C.
func (i *I2C) Close() {
	i.bus.Close()
}

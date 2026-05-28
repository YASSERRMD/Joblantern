// Package printer drives an ESC/POS thermal receipt printer over a
// generic io.Writer (USB serial, network, or a buffer for tests).
package printer

import (
	"fmt"
	"io"
	"strings"
)

// Command bytes per the ESC/POS reference manual.
var (
	cmdInit    = []byte{0x1B, 0x40}
	cmdBoldOn  = []byte{0x1B, 0x45, 0x01}
	cmdBoldOff = []byte{0x1B, 0x45, 0x00}
	cmdCenter  = []byte{0x1B, 0x61, 0x01}
	cmdLeft    = []byte{0x1B, 0x61, 0x00}
	cmdCut     = []byte{0x1D, 0x56, 0x00}
)

// Printer is a minimal wrapper.
type Printer struct {
	W io.Writer
}

// Init resets the printer.
func (p Printer) Init() error { _, err := p.W.Write(cmdInit); return err }

// Line writes a left-aligned line plus newline.
func (p Printer) Line(s string) error {
	_, _ = p.W.Write(cmdLeft)
	_, err := fmt.Fprintln(p.W, s)
	return err
}

// Header writes a centered, bold header line.
func (p Printer) Header(s string) error {
	_, _ = p.W.Write(cmdCenter)
	_, _ = p.W.Write(cmdBoldOn)
	_, err := fmt.Fprintln(p.W, s)
	_, _ = p.W.Write(cmdBoldOff)
	return err
}

// Cut performs a full cut.
func (p Printer) Cut() error { _, err := p.W.Write(cmdCut); return err }

// Rule writes a horizontal divider.
func (p Printer) Rule() error {
	_, err := fmt.Fprintln(p.W, strings.Repeat("-", 42))
	return err
}

package protocol

import (
	"fmt"

	"github.com/flavioheleno/oled-emulator/device"
)

// SPIBridge emulates SPI communication with the display device
type SPIBridge struct {
	device      device.Device
	dcPin       bool // Data/Command pin state
	csPin       bool // Chip Select pin state
	buffer      []byte
	commandMode bool
	dataBuffer  []byte
	commandCode byte
}

// NewSPIBridge creates a new SPI bridge
func NewSPIBridge(dev device.Device) *SPIBridge {
	return &SPIBridge{
		device:      dev,
		dcPin:       false,
		csPin:       false,
		buffer:      make([]byte, 256),
		commandMode: true,
		dataBuffer:  make([]byte, 0),
	}
}

// SetDC sets the Data/Command pin state
// false = command mode, true = data mode
func (sb *SPIBridge) SetDC(state bool) {
	sb.dcPin = state
}

// SetCS sets the Chip Select pin state
// false = selected, true = not selected
func (sb *SPIBridge) SetCS(state bool) {
	sb.csPin = state
}

// Write sends data over SPI
func (sb *SPIBridge) Write(data []byte) error {
	if sb.csPin {
		// Chip not selected, ignore write
		return nil
	}

	if len(data) == 0 {
		return nil
	}

	if sb.dcPin {
		// Data mode
		return sb.writeData(data)
	}

	// Command mode
	return sb.writeCommand(data)
}

// writeCommand processes command bytes
func (sb *SPIBridge) writeCommand(data []byte) error {
	for _, b := range data {
		if err := sb.device.ProcessCommand(b, sb.dataBuffer); err != nil {
			return fmt.Errorf("command error: %w", err)
		}
		sb.dataBuffer = sb.dataBuffer[:0]
		sb.commandCode = b
	}

	return nil
}

// writeData processes data bytes
func (sb *SPIBridge) writeData(data []byte) error {
	// For SSD1322, data mode typically follows a WriteRAM command
	// The device implementation handles writing to VRAM through SetPixel or similar
	// For now, we'll just acknowledge the data
	// A full implementation would process the data into the display buffer

	return nil
}

// Reset performs a hardware reset sequence
func (sb *SPIBridge) Reset() error {
	sb.dataBuffer = sb.dataBuffer[:0]
	return sb.device.Reset()
}

// ReadData reads from the display (if supported)
// Note: This is a placeholder - real SSD1322 does support reading
func (sb *SPIBridge) ReadData(length int) ([]byte, error) {
	result := make([]byte, length)

	// For now, return zeros - real implementation would read VRAM
	for i := 0; i < length; i++ {
		result[i] = 0
	}

	return result, nil
}

// SendInitSequence sends an initialization sequence
func (sb *SPIBridge) SendInitSequence(sequence []byte) error {
	// Command unlock
	sb.SetDC(false)
	if err := sb.Write([]byte{0xFD}); err != nil {
		return fmt.Errorf("unlock command failed: %w", err)
	}

	sb.SetDC(false)
	if err := sb.Write([]byte{0xB1}); err != nil {
		return err
	}

	// Send initialization sequence
	for i := 0; i < len(sequence); i++ {
		if i%2 == 0 {
			// Command byte
			sb.SetDC(false)
			if err := sb.Write([]byte{sequence[i]}); err != nil {
				return err
			}
		} else {
			// Data byte
			sb.SetDC(false)
			if err := sb.Write([]byte{sequence[i]}); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetDevice returns the underlying device
func (sb *SPIBridge) GetDevice() device.Device {
	return sb.device
}

// GetStatus returns the current bridge status
func (sb *SPIBridge) GetStatus() Status {
	return Status{
		DCPin:       sb.dcPin,
		CSPin:       sb.csPin,
		CommandMode: !sb.dcPin,
		LastCommand: sb.commandCode,
	}
}

// Status holds the current status of the SPI bridge
type Status struct {
	DCPin       bool
	CSPin       bool
	CommandMode bool
	LastCommand byte
}

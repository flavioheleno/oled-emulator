package device

import (
	"fmt"
)

// SSD1322 command codes
const (
	// Fundamental Commands
	CmdSetColumnAddress   = 0x15 // Set column address
	CmdSetRowAddress      = 0x75 // Set row address
	CmdWriteRAM           = 0x5C // Write RAM
	CmdReadRAM            = 0x5D // Read RAM
	CmdSetContrast        = 0xC1 // Set contrast
	CmdMasterContrast     = 0xC7 // Master current control
	CmdSetRemap           = 0xA0 // Set remap and dual COM mode
	CmdSetStartLine       = 0xA1 // Set display start line
	CmdDisplayOffset      = 0xA2 // Set display offset
	CmdDisplayMode        = 0xA4 // Set display mode (normal/entire on)
	CmdInvertDisplay      = 0xA6 // Set normal/inverse display
	CmdSetMultiplexRatio  = 0xCA // Set MUX ratio

	// Display On/Off
	CmdSleepMode          = 0xAE // Sleep mode (display OFF)
	CmdNormalDisplay      = 0xAF // Normal mode (display ON)

	// Scrolling Commands
	CmdHorizontalScroll   = 0x26 // Horizontal scroll setup
	CmdContinuousScroll   = 0x27 // Horizontal scroll setup (continuous)
	CmdDeactivateScroll   = 0x2E // Deactivate scroll
	CmdActivateScroll     = 0x2F // Activate scroll

	// Timing and Driving Scheme Commands
	CmdSetClockDivider    = 0xB3 // Set clock divider ratio
	CmdSetPhaseLength     = 0xB1 // Set phase length
	CmdEnhanceDisplay     = 0xB4 // Display enhancement
	CmdSetPrecharge       = 0xBB // Set second precharge period
	CmdSetVCOMH           = 0xBE // Set V_COMH deselect level

	// Grayscale Table
	CmdGrayscaleTable     = 0xB9 // Set default grayscale table

	// Command Lock
	CmdCommandLock        = 0xFD // Set command lock
)

// SSD1322 display controller emulation
type SSD1322 struct {
	*BaseDevice
	memory               *MemoryHelper
	commandLocked        bool
	displayOn            bool
	dataMode             bool // true = data, false = command
	contrastLevel        byte
	masterCurrentLevel   byte
	invertDisplay        bool
	columnStart          int
	columnEnd            int
	rowStart             int
	rowEnd               int
	currentColumn        int
	currentRow           int
	scrollEnabled        bool
	startLine            int
	displayOffset        int
	multiplexRatio       byte
	clockDivider         byte
	phaseLength          byte
	prechargeVoltage     byte
	vcomhLevel           byte
	remapSettings        byte
	grayscaleTableMode   int // 0 = default, 1 = custom
}

// NewSSD1322 creates a new SSD1322 device
func NewSSD1322(width, height int) *SSD1322 {
	config := Config{
		Width:       width,
		Height:      height,
		ColorDepth:  4,
		PixelFormat: HorizontalNibble,
		ColumnOffset: 28, // SSD1322 has 480 internal columns, display starts at column 28
	}

	baseDevice := NewBaseDevice(config)

	ssd1322 := &SSD1322{
		BaseDevice:       baseDevice,
		memory:           NewMemoryHelper(width, height, HorizontalNibble, 28),
		commandLocked:    true,
		displayOn:        false,
		dataMode:         false,
		contrastLevel:    0x7F,
		masterCurrentLevel: 0x0F,
		invertDisplay:    false,
		columnStart:      0,
		columnEnd:        width - 1,
		rowStart:         0,
		rowEnd:           height - 1,
		currentColumn:    0,
		currentRow:       0,
		scrollEnabled:    false,
		startLine:        0,
		displayOffset:    0,
		multiplexRatio:   0x3F,
		clockDivider:     0x00,
		phaseLength:      0x74,
		prechargeVoltage: 0x3C,
		vcomhLevel:       0x07,
		remapSettings:    0x14,
		grayscaleTableMode: 0,
	}

	return ssd1322
}

// ProcessCommand handles SSD1322 commands
func (ssd *SSD1322) ProcessCommand(cmd byte, data []byte) error {
	// Most commands are locked unless unlocked with CmdCommandLock
	switch cmd {
	case CmdCommandLock:
		// Unlock/lock commands (unlock sequence: 0xFD, 0xB1)
		if len(data) > 0 {
			if data[0] == 0xB1 {
				ssd.commandLocked = false
			} else if data[0] == 0xB0 {
				ssd.commandLocked = true
			}
		}
		return nil

	case CmdNormalDisplay:
		ssd.displayOn = true
		return nil

	case CmdSleepMode:
		ssd.displayOn = false
		return nil

	case CmdWriteRAM:
		// This switches to data mode for RAM writing
		ssd.dataMode = true
		return nil

	case CmdReadRAM:
		ssd.dataMode = true
		return nil
	}

	// Commands that require unlock
	if ssd.commandLocked && cmd != CmdCommandLock {
		// Some commands may still be allowed when locked
	}

	switch cmd {
	case CmdSetColumnAddress:
		if len(data) >= 2 {
			ssd.columnStart = int(data[0])
			ssd.columnEnd = int(data[1])
			ssd.currentColumn = ssd.columnStart
		}
		return nil

	case CmdSetRowAddress:
		if len(data) >= 2 {
			ssd.rowStart = int(data[0])
			ssd.rowEnd = int(data[1])
			ssd.currentRow = ssd.rowStart
		}
		return nil

	case CmdSetContrast:
		if len(data) > 0 {
			ssd.contrastLevel = data[0]
		}
		return nil

	case CmdMasterContrast:
		if len(data) > 0 {
			ssd.masterCurrentLevel = data[0] & 0x0F
		}
		return nil

	case CmdInvertDisplay:
		if len(data) > 0 {
			ssd.invertDisplay = (data[0] & 0x01) != 0
		}
		return nil

	case CmdSetMultiplexRatio:
		if len(data) > 0 {
			ssd.multiplexRatio = data[0]
		}
		return nil

	case CmdSetStartLine:
		if len(data) > 0 {
			ssd.startLine = int(data[0] & 0x7F)
		}
		return nil

	case CmdDisplayOffset:
		if len(data) > 0 {
			ssd.displayOffset = int(data[0])
		}
		return nil

	case CmdSetRemap:
		if len(data) > 0 {
			ssd.remapSettings = data[0]
		}
		return nil

	case CmdSetClockDivider:
		if len(data) > 0 {
			ssd.clockDivider = data[0]
		}
		return nil

	case CmdSetPhaseLength:
		if len(data) > 0 {
			ssd.phaseLength = data[0]
		}
		return nil

	case CmdEnhanceDisplay:
		// Display enhancement - typically ignored for emulation
		return nil

	case CmdSetPrecharge:
		if len(data) > 0 {
			ssd.prechargeVoltage = data[0]
		}
		return nil

	case CmdSetVCOMH:
		if len(data) > 0 {
			ssd.vcomhLevel = data[0]
		}
		return nil

	case CmdGrayscaleTable:
		if len(data) > 0 {
			ssd.grayscaleTableMode = int(data[0])
		}
		return nil

	case CmdDeactivateScroll:
		ssd.scrollEnabled = false
		return nil

	case CmdActivateScroll:
		ssd.scrollEnabled = true
		return nil

	case CmdHorizontalScroll:
		if len(data) >= 5 {
			ssd.scrollEnabled = true
		}
		return nil

	case CmdDisplayMode:
		// 0xA4 = normal, 0xA5 = entire display ON
		return nil

	default:
		// Unknown command - silently ignore
		return nil
	}
}

// WriteData writes pixel data to VRAM at current addressing position
func (ssd *SSD1322) WriteData(data []byte) error {
	if !ssd.dataMode {
		return fmt.Errorf("not in data write mode")
	}

	for _, byteVal := range data {
		// Each byte contains 2 pixels (4-bit each)
		// Convert from VRAM column addressing to display coordinates
		col := ssd.currentColumn
		row := ssd.currentRow

		if col >= ssd.columnStart && col <= ssd.columnEnd &&
			row >= ssd.rowStart && row <= ssd.rowEnd {

			// Get actual display coordinates
			// (accounting for column offset)
			displayCol := col - ssd.columnStart

			if displayCol < ssd.Width() {
				// Write lower nibble (first pixel)
				pixel1 := byteVal & 0x0F
				if err := ssd.memory.SetPixelNibble(ssd.vram, displayCol, row, pixel1); err == nil {
					ssd.MarkDirty(displayCol, row, displayCol, row)
				}

				// Write upper nibble (second pixel)
				displayCol++
				if displayCol < ssd.Width() {
					pixel2 := (byteVal >> 4) & 0x0F
					if err := ssd.memory.SetPixelNibble(ssd.vram, displayCol, row, pixel2); err == nil {
						ssd.MarkDirty(displayCol, row, displayCol, row)
					}
				}

				// Advance to next column pair
				ssd.currentColumn++
				if ssd.currentColumn > ssd.columnEnd {
					ssd.currentColumn = ssd.columnStart
					ssd.currentRow++
					if ssd.currentRow > ssd.rowEnd {
						ssd.currentRow = ssd.rowStart
					}
				}
			}
		}
	}

	return nil
}

// SetPixel implements the Device interface
func (ssd *SSD1322) SetPixel(x, y int, color byte) error {
	if x < 0 || x >= ssd.Width() || y < 0 || y >= ssd.Height() {
		return fmt.Errorf("pixel out of bounds: (%d, %d)", x, y)
	}

	if err := ssd.memory.SetPixelNibble(ssd.vram, x, y, color&0x0F); err != nil {
		return err
	}

	ssd.MarkDirty(x, y, x, y)
	return nil
}

// GetPixel implements the Device interface
func (ssd *SSD1322) GetPixel(x, y int) (byte, error) {
	return ssd.memory.GetPixelNibble(ssd.vram, x, y)
}

// Reset performs a hardware reset
func (ssd *SSD1322) Reset() error {
	// Clear VRAM
	for i := range ssd.vram {
		ssd.vram[i] = 0
	}

	// Reset all settings to default
	ssd.commandLocked = true
	ssd.displayOn = false
	ssd.dataMode = false
	ssd.contrastLevel = 0x7F
	ssd.masterCurrentLevel = 0x0F
	ssd.invertDisplay = false
	ssd.columnStart = 0
	ssd.columnEnd = ssd.Width() - 1
	ssd.rowStart = 0
	ssd.rowEnd = ssd.Height() - 1
	ssd.currentColumn = 0
	ssd.currentRow = 0
	ssd.scrollEnabled = false
	ssd.startLine = 0
	ssd.displayOffset = 0

	ssd.MarkDirty(0, 0, ssd.Width()-1, ssd.Height()-1)
	return nil
}

// IsDisplayOn returns whether the display is powered on
func (ssd *SSD1322) IsDisplayOn() bool {
	return ssd.displayOn
}

// GetContrastLevel returns current contrast
func (ssd *SSD1322) GetContrastLevel() byte {
	return ssd.contrastLevel
}

// IsInverted returns whether display is inverted
func (ssd *SSD1322) IsInverted() bool {
	return ssd.invertDisplay
}

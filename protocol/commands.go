package protocol

import (
	"fmt"
)

// CommandInfo holds information about a command
type CommandInfo struct {
	Code        byte
	Name        string
	Description string
	DataBytes   int
}

// SSD1322Commands defines all SSD1322 commands
var SSD1322Commands = map[byte]CommandInfo{
	// Fundamental Commands
	0x15: {Code: 0x15, Name: "SetColumnAddress", Description: "Set column address", DataBytes: 2},
	0x75: {Code: 0x75, Name: "SetRowAddress", Description: "Set row address", DataBytes: 2},
	0x5C: {Code: 0x5C, Name: "WriteRAM", Description: "Write RAM", DataBytes: 0},
	0x5D: {Code: 0x5D, Name: "ReadRAM", Description: "Read RAM", DataBytes: 0},

	// Fundamental Commands - Contrast
	0xC1: {Code: 0xC1, Name: "SetContrast", Description: "Set contrast", DataBytes: 1},
	0xC7: {Code: 0xC7, Name: "MasterCurrentControl", Description: "Master current control", DataBytes: 1},

	// Display Setup
	0xA0: {Code: 0xA0, Name: "SetRemap", Description: "Set remap and dual COM mode", DataBytes: 1},
	0xA1: {Code: 0xA1, Name: "SetStartLine", Description: "Set display start line", DataBytes: 1},
	0xA2: {Code: 0xA2, Name: "DisplayOffset", Description: "Set display offset", DataBytes: 1},
	0xA4: {Code: 0xA4, Name: "DisplayMode", Description: "Set display mode", DataBytes: 0},
	0xA5: {Code: 0xA5, Name: "EntireDisplayON", Description: "Entire display ON", DataBytes: 0},
	0xA6: {Code: 0xA6, Name: "NormalDisplay", Description: "Normal display", DataBytes: 1},
	0xA7: {Code: 0xA7, Name: "InverseDisplay", Description: "Inverse display", DataBytes: 0},
	0xAE: {Code: 0xAE, Name: "SleepMode", Description: "Sleep mode (display OFF)", DataBytes: 0},
	0xAF: {Code: 0xAF, Name: "NormalMode", Description: "Normal mode (display ON)", DataBytes: 0},

	// MUX Ratio & Timing
	0xCA: {Code: 0xCA, Name: "SetMultiplexRatio", Description: "Set MUX ratio", DataBytes: 1},
	0xB3: {Code: 0xB3, Name: "SetClockDivider", Description: "Set clock divider ratio", DataBytes: 1},
	0xB1: {Code: 0xB1, Name: "SetPhaseLength", Description: "Set phase length", DataBytes: 1},
	0xBB: {Code: 0xBB, Name: "SetPrecharge", Description: "Set second precharge period", DataBytes: 1},
	0xBE: {Code: 0xBE, Name: "SetVCOMH", Description: "Set V_COMH deselect level", DataBytes: 1},

	// Display Enhancement
	0xB4: {Code: 0xB4, Name: "DisplayEnhance", Description: "Display enhancement", DataBytes: 1},
	0xD1: {Code: 0xD1, Name: "DisplayEnhanceB", Description: "Display enhancement B", DataBytes: 1},

	// Scrolling
	0x26: {Code: 0x26, Name: "HorizontalScroll", Description: "Horizontal scroll setup", DataBytes: 5},
	0x27: {Code: 0x27, Name: "ContinuousScroll", Description: "Horizontal scroll setup (continuous)", DataBytes: 5},
	0x2E: {Code: 0x2E, Name: "DeactivateScroll", Description: "Deactivate scroll", DataBytes: 0},
	0x2F: {Code: 0x2F, Name: "ActivateScroll", Description: "Activate scroll", DataBytes: 0},

	// Grayscale
	0xB9: {Code: 0xB9, Name: "GrayscaleTable", Description: "Set default grayscale table", DataBytes: 1},

	// Command Lock
	0xFD: {Code: 0xFD, Name: "CommandLock", Description: "Set command lock", DataBytes: 1},
}

// GetCommandInfo returns information about a command
func GetCommandInfo(code byte) (CommandInfo, error) {
	info, ok := SSD1322Commands[code]
	if !ok {
		return CommandInfo{}, fmt.Errorf("unknown command: 0x%02X", code)
	}
	return info, nil
}

// CommandBuilder helps construct SPI command sequences
type CommandBuilder struct {
	bytes []byte
}

// NewCommandBuilder creates a new command builder
func NewCommandBuilder() *CommandBuilder {
	return &CommandBuilder{
		bytes: make([]byte, 0),
	}
}

// AddCommand adds a command byte
func (cb *CommandBuilder) AddCommand(code byte) *CommandBuilder {
	cb.bytes = append(cb.bytes, code)
	return cb
}

// AddData adds a data byte
func (cb *CommandBuilder) AddData(data byte) *CommandBuilder {
	cb.bytes = append(cb.bytes, data)
	return cb
}

// AddBytes adds multiple bytes
func (cb *CommandBuilder) AddBytes(data ...byte) *CommandBuilder {
	cb.bytes = append(cb.bytes, data...)
	return cb
}

// Build returns the command bytes
func (cb *CommandBuilder) Build() []byte {
	result := make([]byte, len(cb.bytes))
	copy(result, cb.bytes)
	return result
}

// Reset clears the builder
func (cb *CommandBuilder) Reset() *CommandBuilder {
	cb.bytes = cb.bytes[:0]
	return cb
}

// SSD1322InitSequence generates a typical initialization sequence for SSD1322
func SSD1322InitSequence() []byte {
	builder := NewCommandBuilder()

	// Command unlock
	builder.AddCommand(0xFD).AddData(0xB1)

	// Display OFF
	builder.AddCommand(0xAE)

	// Set clock divider
	builder.AddCommand(0xB3).AddData(0x00)

	// Set MUX ratio
	builder.AddCommand(0xCA).AddData(0x3F)

	// Set display offset
	builder.AddCommand(0xA2).AddData(0x00)

	// Set display start line
	builder.AddCommand(0xA1).AddData(0x00)

	// Set remap
	builder.AddCommand(0xA0).AddData(0x14)

	// Set phase length
	builder.AddCommand(0xB1).AddData(0x74)

	// Display enhancement
	builder.AddCommand(0xB4).AddData(0x00)

	// Set contrast
	builder.AddCommand(0xC1).AddData(0x7F)

	// Master current
	builder.AddCommand(0xC7).AddData(0x0F)

	// Set precharge
	builder.AddCommand(0xBB).AddData(0x3C)

	// Set VCOMH
	builder.AddCommand(0xBE).AddData(0x07)

	// Normal display
	builder.AddCommand(0xA6)

	// Column addressing
	builder.AddCommand(0x15).AddData(0x1C).AddData(0x5B)

	// Row addressing
	builder.AddCommand(0x75).AddData(0x00).AddData(0x3F)

	// Display ON
	builder.AddCommand(0xAF)

	return builder.Build()
}

// DrawPixelCommand creates a command sequence to draw a pixel
func DrawPixelCommand(x, y, color byte) []byte {
	builder := NewCommandBuilder()

	// Set column address
	builder.AddCommand(0x15).AddData(x).AddData(x)

	// Set row address
	builder.AddCommand(0x75).AddData(y).AddData(y)

	// Write RAM
	builder.AddCommand(0x5C)

	// Data byte contains 2 pixels (color is 4-bit)
	builder.AddData((color << 4) | color)

	return builder.Build()
}

// FillScreenCommand creates a command sequence to fill the entire screen
func FillScreenCommand(color byte) []byte {
	builder := NewCommandBuilder()

	// Set column address
	builder.AddCommand(0x15).AddData(0x1C).AddData(0x5B)

	// Set row address
	builder.AddCommand(0x75).AddData(0x00).AddData(0x3F)

	// Write RAM
	builder.AddCommand(0x5C)

	// Add color data for entire screen
	// 256x64 display = 480 columns (internal) x 64 rows = 15360 bytes
	for i := 0; i < 7680; i++ {
		builder.AddData((color << 4) | color)
	}

	return builder.Build()
}

// ContrastCommand creates a command to set contrast
func ContrastCommand(level byte) []byte {
	return NewCommandBuilder().
		AddCommand(0xC1).
		AddData(level).
		Build()
}

// InversionCommand creates a command to set display inversion
func InversionCommand(inverted bool) []byte {
	if inverted {
		return NewCommandBuilder().AddCommand(0xA7).Build()
	}
	return NewCommandBuilder().AddCommand(0xA6).Build()
}

// PowerCommand creates a command to control power
func PowerCommand(on bool) []byte {
	if on {
		return NewCommandBuilder().AddCommand(0xAF).Build()
	}
	return NewCommandBuilder().AddCommand(0xAE).Build()
}

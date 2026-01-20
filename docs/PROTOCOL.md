# SSD1322 Protocol Reference

## Overview

The SSD1322 is a 256×64 monochrome/grayscale display driver IC from Solomon Systech. It supports 4-bit grayscale (16 levels) and communicates via SPI, 8080 parallel, or 6800 parallel interfaces.

This emulator focuses on the SPI protocol and command set.

## Display Specifications

- **Display Size**: 256×64 pixels
- **Internal RAM**: 480×64 pixels (larger than display, for smooth scrolling)
- **Color Depth**: 4-bit grayscale (16 levels)
- **Pixel Format**: 2 pixels per byte (horizontal nibble packing)
- **Interface**: SPI, 8080 8-bit, 6800 8-bit
- **Supply Voltage**: 3.3V

## Communication Protocol

### SPI Interface

```
DC (Data/Command) Pin:
  - LOW: Command mode (send commands)
  - HIGH: Data mode (send pixel data)

CS (Chip Select) Pin:
  - LOW: Chip selected
  - HIGH: Chip deselected

CLK: Serial clock
MOSI: Master Out, Slave In (data from MCU to display)
```

### Command Format

Commands are sent with DC=LOW:

```
Byte 0: Command code (0x00-0xFF)
Byte 1-N: Optional data bytes (depending on command)
```

Data is sent with DC=HIGH:

```
Byte 0-N: VRAM data (2 pixels per byte in nibble format)
```

## Command Set

### Fundamental Commands

#### Set Column Address (0x15)
```
Data bytes: 2
Byte 0: Start column (0x1C-0x5B for 256 pixel width)
Byte 1: End column (0x1C-0x5B)

Example: 0x15, 0x1C, 0x5B  // Set to full width
```

#### Set Row Address (0x75)
```
Data bytes: 2
Byte 0: Start row (0x00-0x3F)
Byte 1: End row (0x00-0x3F)

Example: 0x75, 0x00, 0x3F  // Set to full height
```

#### Write RAM (0x5C)
```
Data bytes: Variable (display data)
Starts write to VRAM at current address

Example sequence:
  0x15, 0x1C, 0x5B  // Set column
  0x75, 0x00, 0x3F  // Set row
  0x5C              // Start write
  [pixel data...]   // Send 15360 bytes for full screen
```

#### Read RAM (0x5D)
```
Reads from VRAM at current address
Returns: Pixel data bytes
```

### Contrast & Current Control

#### Set Contrast (0xC1)
```
Data bytes: 1
Byte 0: Contrast level (0x00-0xFF)

Default: 0x7F
Example: 0xC1, 0x80
```

#### Master Current Control (0xC7)
```
Data bytes: 1
Byte 0: Current level (0x00-0x0F)

Default: 0x0F (100%)
Example: 0xC7, 0x08
```

### Display Setup

#### Set Remap & Dual COM Mode (0xA0)
```
Data bytes: 1
Byte 0: Remap settings

Bit 0: Column address remap (0=left-to-right, 1=right-to-left)
Bit 1: Nibble remap (0=normal, 1=swapped)
Bit 2: Row address remap (0=top-to-bottom, 1=bottom-to-top)
Bit 3: Dual COM mode (0=disable, 1=enable)
Bit 4: Reserved
Bit 5: Left/Right swap

Default: 0x14
Example: 0xA0, 0x14
```

#### Set Display Start Line (0xA1)
```
Data bytes: 1
Byte 0: Start line (0x00-0x3F)

Default: 0x00
Example: 0xA1, 0x10
```

#### Set Display Offset (0xA2)
```
Data bytes: 1
Byte 0: Offset value (0x00-0x3F)

Default: 0x00
Example: 0xA2, 0x20
```

#### Display Mode (0xA4)
```
Data bytes: 0
Entire display follows RAM content

Example: 0xA4
```

#### Entire Display ON (0xA5)
```
Data bytes: 0
Entire display ON regardless of RAM

Example: 0xA5
```

#### Normal/Inverse Display (0xA6 / 0xA7)
```
0xA6: Normal display (on=white, off=black)
0xA7: Inverse display (on=black, off=white)

Data bytes: 0
Examples: 0xA6  or  0xA7
```

#### Sleep Mode (0xAE)
```
Data bytes: 0
Turn display OFF (low power)

Example: 0xAE
```

#### Normal Mode (0xAF)
```
Data bytes: 0
Turn display ON

Example: 0xAF
```

### Timing & Driving

#### Set Multiplex Ratio (0xCA)
```
Data bytes: 1
Byte 0: MUX ratio (0x0F-0x7F, typical 0x3F for 64 rows)

Default: 0x3F
Example: 0xCA, 0x3F
```

#### Set Clock Divider (0xB3)
```
Data bytes: 1
Byte 0: Divider (Bits 0-3) and Oscillator Frequency (Bits 4-7)

Default: 0x00
Example: 0xB3, 0x00
```

#### Set Phase Length (0xB1)
```
Data bytes: 1
Byte 0: Phase1 length (Bits 0-3) and Phase2 length (Bits 4-7)

Default: 0x74
Example: 0xB1, 0x74
```

#### Set Precharge Period (0xBB)
```
Data bytes: 1
Byte 0: Precharge voltage level

Default: 0x3C
Example: 0xBB, 0x3C
```

#### Set VCOMH Level (0xBE)
```
Data bytes: 1
Byte 0: VCOMH level (0x00-0xFF)

Default: 0x07
Example: 0xBE, 0x07
```

### Display Enhancement

#### Display Enhancement (0xB4 / 0xD1)
```
Data bytes: 1
Byte 0: Enhancement settings

Enables external VSL supply

Default: 0x00
Examples: 0xB4, 0x00  or  0xD1, 0x00
```

### Scrolling

#### Horizontal Scroll Setup (0x26 / 0x27)
```
0x26: Setup horizontal scroll (left or right)
0x27: Setup horizontal scroll - continuous mode

Data bytes: 5
Byte 0: Row address start
Byte 1: Number of rows to scroll
Byte 2: Scroll speed
Byte 3: Row address end (for fixed mode)
Byte 4: Offset (for continuous mode)

Example: 0x26, 0x00, 0x3F, 0x00, 0x00, 0x00
```

#### Deactivate Scroll (0x2E)
```
Data bytes: 0
Stops scrolling

Example: 0x2E
```

#### Activate Scroll (0x2F)
```
Data bytes: 0
Starts scrolling with previous settings

Example: 0x2F
```

### Grayscale

#### Set Grayscale Table (0xB9)
```
Data bytes: 1 or 64
Byte 0: Grayscale table mode
  - Bit 0: Use default grayscale table if 0
  - Other bits for custom table selection

Example: 0xB9, 0x00  // Use default
```

### Command Lock (0xFD)

#### Unlock/Lock Commands
```
Data bytes: 1
Byte 0: Lock code
  - 0xB1: Unlock (enable most commands)
  - 0xB0: Lock (disable commands)

Default: Locked

Example unlock:
  0xFD, 0xB1  // Unlock
  [send commands]
  0xFD, 0xB0  // Lock
```

## Typical Initialization Sequence

```go
// 1. Unlock commands
0xFD, 0xB1

// 2. Display OFF
0xAE

// 3. Set clock divider
0xB3, 0x00

// 4. Set MUX ratio
0xCA, 0x3F

// 5. Set display offset
0xA2, 0x00

// 6. Set start line
0xA1, 0x00

// 7. Set remap
0xA0, 0x14

// 8. Set phase length
0xB1, 0x74

// 9. Set precharge
0xBB, 0x3C

// 10. Set VCOMH
0xBE, 0x07

// 11. Set contrast
0xC1, 0x7F

// 12. Master current
0xC7, 0x0F

// 13. Normal display
0xA6

// 14. Display ON
0xAF
```

## Data Formats

### Pixel Data Format (Horizontal Nibble)

Each byte contains 2 pixels:

```
Byte: [Pixel N+1 (upper nibble)] [Pixel N (lower nibble)]

Example: 0xF5
  Pixel N = 0x5 (darker)
  Pixel N+1 = 0xF (brighter)
```

### Grayscale Levels

```
0x0 = Black (off)
0x1 = Very dark
0x2 = Dark
0x3 = Dark gray
...
0xE = Light gray
0xF = White (on)
```

## Color/Brightness Mapping

When mapping 8-bit grayscale to 4-bit:

```go
// Convert 8-bit (0-255) to 4-bit (0-15)
bit4 = bit8 >> 4

// Convert 4-bit to 8-bit for display
bit8 = (bit4 << 4) | bit4
```

## Error Handling

Most commands should be followed by a small delay (typically 1-10ms) to allow the controller time to process.

The display maintains internal state. Always:
1. Unlock commands before sending configuration
2. Set addressing mode before writing/reading RAM
3. Use Reset (hardware or software) to return to known state

## VRAM Addressing

VRAM is organized as:
- 480 columns × 64 rows × 4-bit (nibble)
- Display area: columns 28-83, all 64 rows
- Addressing mode: Column-major (set column, then row)

```
Total bytes: (480 × 64) / 2 = 15,360 bytes
```

## Power Consumption

- Display ON: ~50mA (typical)
- Display OFF (sleep): <1µA
- Max refresh rate: ~100 Hz

## Reference Material

- SSD1322 Datasheet
- Application notes from Solomon Systech
- Common implementation examples from embedded projects

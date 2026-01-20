package protocol

import (
	"testing"

	"github.com/flavioheleno/oled-emulator/device"
)

func TestSPIBridgeBasic(t *testing.T) {
	dev := device.NewSSD1322(256, 64)
	bridge := NewSPIBridge(dev)

	// Test DC pin control
	bridge.SetDC(false) // Command mode
	if bridge.dcPin != false {
		t.Error("DC pin should be false (command mode)")
	}

	bridge.SetDC(true) // Data mode
	if bridge.dcPin != true {
		t.Error("DC pin should be true (data mode)")
	}
}

func TestCommandBuilder(t *testing.T) {
	builder := NewCommandBuilder()
	builder.AddCommand(0xFD).AddData(0xB1).AddCommand(0xAE)

	bytes := builder.Build()
	if len(bytes) != 3 {
		t.Errorf("expected 3 bytes, got %d", len(bytes))
	}

	expected := []byte{0xFD, 0xB1, 0xAE}
	for i, b := range bytes {
		if b != expected[i] {
			t.Errorf("byte %d: expected 0x%02X, got 0x%02X", i, expected[i], b)
		}
	}
}

func TestGetCommandInfo(t *testing.T) {
	info, err := GetCommandInfo(0xFD)
	if err != nil {
		t.Fatalf("failed to get command info: %v", err)
	}

	if info.Name != "CommandLock" {
		t.Errorf("expected CommandLock, got %s", info.Name)
	}

	if info.DataBytes != 1 {
		t.Errorf("expected 1 data byte, got %d", info.DataBytes)
	}
}

func TestUnknownCommand(t *testing.T) {
	_, err := GetCommandInfo(0xFF)
	if err == nil {
		t.Error("should return error for unknown command")
	}
}

func TestSSD1322InitSequence(t *testing.T) {
	seq := SSD1322InitSequence()
	if len(seq) == 0 {
		t.Error("init sequence should not be empty")
	}

	// Verify it starts with command unlock
	if seq[0] != 0xFD || seq[1] != 0xB1 {
		t.Error("init sequence should start with command unlock")
	}

	// Verify it ends with display ON
	if seq[len(seq)-1] != 0xAF {
		t.Error("init sequence should end with display ON command")
	}
}

func TestContrastCommand(t *testing.T) {
	cmd := ContrastCommand(0x80)
	if len(cmd) != 2 {
		t.Errorf("expected 2 bytes, got %d", len(cmd))
	}

	if cmd[0] != 0xC1 || cmd[1] != 0x80 {
		t.Errorf("expected [0xC1, 0x80], got [0x%02X, 0x%02X]", cmd[0], cmd[1])
	}
}

func TestInversionCommand(t *testing.T) {
	// Test normal
	cmd := InversionCommand(false)
	if cmd[0] != 0xA6 {
		t.Errorf("expected 0xA6 for normal, got 0x%02X", cmd[0])
	}

	// Test inverted
	cmd = InversionCommand(true)
	if cmd[0] != 0xA7 {
		t.Errorf("expected 0xA7 for inverted, got 0x%02X", cmd[0])
	}
}

func TestPowerCommand(t *testing.T) {
	// Test ON
	cmd := PowerCommand(true)
	if cmd[0] != 0xAF {
		t.Errorf("expected 0xAF for power ON, got 0x%02X", cmd[0])
	}

	// Test OFF
	cmd = PowerCommand(false)
	if cmd[0] != 0xAE {
		t.Errorf("expected 0xAE for power OFF, got 0x%02X", cmd[0])
	}
}

func TestCommandBuilderReset(t *testing.T) {
	builder := NewCommandBuilder()
	builder.AddCommand(0xFD).AddData(0xB1)

	if len(builder.Build()) != 2 {
		t.Error("first build should have 2 bytes")
	}

	builder.Reset()
	if len(builder.Build()) != 0 {
		t.Error("after reset, build should be empty")
	}
}

func TestSPIBridgeStatus(t *testing.T) {
	dev := device.NewSSD1322(256, 64)
	bridge := NewSPIBridge(dev)

	bridge.SetDC(false)
	bridge.SetCS(true)

	status := bridge.GetStatus()
	if status.DCPin != false || status.CSPin != true {
		t.Error("status should reflect pin states")
	}

	if status.CommandMode != true {
		t.Error("should be in command mode")
	}
}

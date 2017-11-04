package braster

import (
	"errors"
	"io"
)

// defined at <Raster Command Reference> P22
var commandInvalidate = []byte{0}
var commandInitialize = []byte{0x1B, 0x40}
var commandStatusInfoRequest = []byte{0x1B, 0x69, 0x53}
var commandSwitchDynamicCommandMode = []byte{0x1B, 0x69, 0x61}
var commandPrintInfoCommand = []byte{0x1B, 0x69, 0x7A}
var commandVariousModeSettings = []byte{0x1B, 0x69, 0x4D}
var commandAdvancedModeSettings = []byte{0x1B, 0x69, 0x4B}
var commandSpecifyMarginAmount = []byte{0x1B, 0x69, 0x64}
var commandSelectCompressionMode = []byte{0x4D}
var commandRasterGrapicsTransfer = []byte{0x67}
var commandZeroRasterGraphics = []byte{0x5A}
var commandPrint = []byte{0x0C}
var commandPrintWithFeeding = []byte{0x1A}

// CommandBuilder used to build and send a command
type CommandBuilder struct {
	io.Writer
}

// Invalidate sends invalidate command,
// if data transmission is to be stopped midway, send the “initialize” command after sending the “invalidate” command for the appropriate number of bytes to return to the receiving state, where the print buffer is cleared.
func (cb CommandBuilder) Invalidate(repeat int) error {
	for i := 0; i < repeat; i++ {
		if _, err := cb.Write(commandInvalidate); err != nil {
			return err
		}
	}
	return nil
}

// Initialize initializes mode settings, also used to cancel printing.
func (cb CommandBuilder) Initialize() error {
	_, err := cb.Write(commandInitialize)
	return err
}

// StatusInformationRequest send a request to the printer for status information.
func (cb CommandBuilder) StatusInformationRequest() error {
	_, err := cb.Write(commandStatusInfoRequest)
	return err
}

// DynamicCommandMode used as parameters of command SwitchDynamicCommandMode
// defined at <Raster Command Reference> P31
type DynamicCommandMode int

const (
	// ESCPMode default mode
	ESCPMode = DynamicCommandMode(0)
	// RasterMode (Be sure to switch to this mode)
	RasterMode = DynamicCommandMode(1)
	// PtouchTemplateMode P-touch Template mode
	PtouchTemplateMode = DynamicCommandMode(2)
)

// SwitchDynamicCommandMode dynamically switches between the printer's command modes. A printer that receives this command operates in the specified command mode until the printer is turned off
// The printer must be switched to raster mode before raster data is sent to it. Therefore, send this command to switch the printer to raster mode.
func (cb CommandBuilder) SwitchDynamicCommandMode(mode DynamicCommandMode) error {
	_, err := cb.Write(commandSwitchDynamicCommandMode)
	if err != nil {
		return err
	}
	_, err = cb.Write([]byte{byte(mode)})
	return err
}

func (cb CommandBuilder) PrintInformationCommand() error {
	panic("not implements")
}

// VariousMode used as parameters of command VariousModeSettings
// defined at <Raster Command Reference> P33
type VariousMode int

const (
	// AutoCut automatically cuts
	AutoCut = VariousMode(1 << 6)
	// MirrorPrinting mirror printing
	MirrorPrinting = VariousMode(1 << 7)
)

// VariousModeSettings sends various mode settings command.
func (cb CommandBuilder) VariousModeSettings(mode VariousMode) error {
	_, err := cb.Write(commandVariousModeSettings)
	if err != nil {
		return err
	}
	_, err = cb.Write([]byte{byte(mode)})
	return err
}

// AdvancedMode used as parameters of command AdvancedModeSettings
// defined at <Raster Command Reference> P33
type AdvancedMode int

const (
	// NoChainPrinting feeding and cutting are performed after the last one is printed.
	NoChainPrinting = AdvancedMode(1 << 3)
	// ChainPrinting feeding and cutting are not performed after the last one is printed.
	ChainPrinting = AdvancedMode(0 << 3)

	// SpecialTape labels are not cut when special tape is installed.
	SpecialTape = AdvancedMode(1 << 4)
	// NoCutting same as SpecialTape
	NoCutting = AdvancedMode(1 << 4)

	// NoBufferClearingWhenPrinting the expansion buffer of the machine is not cleared with the “no buffer clearing when printing” command.
	// If this command is sent when the data of the first label is printed (it is specified between the “initialize” command and the print data), printing is possible only if a print command is sent with the second or later label.
	NoBufferClearingWhenPrinting = AdvancedMode(1 << 7)
)

// AdvancedModeSettings sends advanced mode settings command.
func (cb CommandBuilder) AdvancedModeSettings(mode AdvancedMode) error {
	_, err := cb.Write(commandAdvancedModeSettings)
	if err != nil {
		return err
	}
	_, err = cb.Write([]byte{byte(mode)})
	return err
}

var errMarginAmountUnaccpectable = errors.New("margin amount unaccpectable")

// SpecifyMarginAmount specifies the amount of the margins.
func (cb CommandBuilder) SpecifyMarginAmount(dots int) error {
	if dots < 0 || dots > 0xFFFF {
		return errMarginAmountUnaccpectable
	}

	_, err := cb.Write(commandSpecifyMarginAmount)
	if err != nil {
		return err
	}
	buf := []byte{byte(dots), byte(dots >> 8)}
	_, err = cb.Write(buf[:])
	return err
}

// CompressionMode used as parameters of command SelectCompressionMode
// defined at <Raster Command Reference> P35
type CompressionMode int

const (
	// NoCompression no compression
	NoCompression = CompressionMode(0)
	// TIFF enable TIFF compression
	TIFF = CompressionMode(2)
)

// SelectCompressionMode selects the compression mode.
// Data compression is available only for data in raster graphic transfer
func (cb CommandBuilder) SelectCompressionMode(mode CompressionMode) error {
	_, err := cb.Write(commandSelectCompressionMode)
	if err != nil {
		return err
	}
	_, err = cb.Write([]byte{byte(mode)})
	return err
}

var errDataTooLong = errors.New("data too long")

// RasterGraphicsTransfer transfers the specified number of bytes of data.
// ses <Raster Command Reference> P37
func (cb CommandBuilder) RasterGraphicsTransfer(data []byte) error {
	if len(data) > 0xFFFF {
		return errDataTooLong
	}

	_, err := cb.Write(commandRasterGrapicsTransfer)
	if err != nil {
		return err
	}
	buf := []byte{byte(len(data)), byte(len(data) >> 8)}
	_, err = cb.Write(buf[:])
	if err != nil {
		return err
	}

	_, err = cb.Write(data)
	return err
}

// ZeroRasterGraphics fills raster line with 0 data.
func (cb CommandBuilder) ZeroRasterGraphics() error {
	_, err := cb.Write(commandZeroRasterGraphics)
	return err
}

// Print used as a print command at the end of pages other than the last page when multiple pages are printed.
func (cb CommandBuilder) Print() error {
	_, err := cb.Write(commandPrint)
	return err
}

// PrintWithFeeding used as a print command at the end of the last page.
func (cb CommandBuilder) PrintWithFeeding() error {
	_, err := cb.Write(commandPrintWithFeeding)
	return err
}

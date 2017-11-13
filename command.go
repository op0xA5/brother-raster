package braster

import (
	"encoding/binary"
	"errors"
	"io"
)

// defined at <Raster Command Reference> P22
var commandInvalidate = [16]byte{}
var commandInitialize = []byte{0x1B, 0x40}
var commandStatusInfoRequest = []byte{0x1B, 0x69, 0x53}
var commandSwitchDynamicCommandMode = []byte{0x1B, 0x69, 0x61}
var commandPrintInfoCommand = []byte{0x1B, 0x69, 0x7A}
var commandVariousModeSettings = []byte{0x1B, 0x69, 0x4D}
var commandAdvancedModeSettings = []byte{0x1B, 0x69, 0x4B}
var commandSpecifyMarginAmount = []byte{0x1B, 0x69, 0x64}
var commandSelectCompressionMode = []byte{0x4D}
var commandRasterGrapicsTransfer = []byte{0x47}
var commandZeroRasterGraphics = []byte{0x5A}
var commandPrint = []byte{0x0C}
var commandPrintWithFeeding = []byte{0x1A}

var commandEndian = binary.LittleEndian

// CommandBuilder used to build and send a command
type CommandBuilder struct {
	io.Writer
}

// Invalidate sends invalidate command,
// if data transmission is to be stopped midway, send the “initialize” command after sending the “invalidate” command for the appropriate number of bytes to return to the receiving state, where the print buffer is cleared.
func (cb CommandBuilder) Invalidate(repeat int) error {
	var err error
	for i := 0; i < repeat; i += len(commandInvalidate) {
		if repeat-i > len(commandInvalidate) {
			_, err = cb.Write(commandInvalidate[:])
		} else {
			_, err = cb.Write(commandInvalidate[:repeat-i])
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (cb CommandBuilder) cmd(b []byte) error {
	_, err := cb.Write(b)
	return err
}
func (cb CommandBuilder) cmd8(b []byte, p byte) error {
	buf := make([]byte, len(b)+1)
	n := copy(buf, b)
	buf[n] = p
	_, err := cb.Write(buf)
	return err
}
func (cb CommandBuilder) cmd16(b []byte, p uint16) error {
	buf := make([]byte, len(b)+2)
	n := copy(buf, b)
	commandEndian.PutUint16(buf[n:], p)
	_, err := cb.Write(buf)
	return err
}

// Initialize initializes mode settings, also used to cancel printing.
func (cb CommandBuilder) Initialize() error {
	return cb.cmd(commandInitialize)
}

// StatusInformationRequest send a request to the printer for status information.
func (cb CommandBuilder) StatusInformationRequest() error {
	return cb.cmd(commandStatusInfoRequest)
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
	return cb.cmd8(commandSwitchDynamicCommandMode, byte(mode))
}

// Specifies which values are valid
const (
	PrintInformationKind    = 0x02
	PrintInformationWidth   = 0x04
	PrintInformationLength  = 0x08
	PrintInformationQuality = 0x40 // Not used
	PrintInformationRecover = 0x80 // always on
)

// PrintInformationCommandParameter specifies the print information.
type PrintInformationCommandParameter struct {
	Flag            int
	MediaType       MediaType
	MediaWidth      MediaWidth
	MediaLength     int
	RasterNunmber   int
	NotStartingPage bool
}

var errInvalidMediaType = errors.New("invalid media type")
var errInvalidMediaWidth = errors.New("invalid media width")
var errInvalidMediaLength = errors.New("invalid media length")
var errInvalidRasterNunmber = errors.New("invalid raster number")

// PrintInformationCommand specifies the print information.
func (cb CommandBuilder) PrintInformationCommand(p PrintInformationCommandParameter) error {
	if p.Flag&PrintInformationKind != 0 && !p.MediaType.IsValid() {
		return errInvalidMediaType
	}
	if p.MediaWidth < 0 || p.MediaWidth > 256 {
		return errInvalidMediaWidth
	}
	if p.MediaLength < 0 || p.MediaLength > 256 {
		return errInvalidMediaLength
	}
	if p.RasterNunmber < 0 || p.RasterNunmber > 0xFFFFFFFF {
		return errInvalidRasterNunmber
	}

	buf := make([]byte, len(commandPrintInfoCommand)+10)
	n := copy(buf, commandPrintInfoCommand)
	buf[n+0] = byte(p.Flag | PrintInformationRecover)
	if p.Flag&PrintInformationKind != 0 {
		buf[n+1] = byte(p.MediaType)
	}
	if p.Flag&PrintInformationWidth != 0 {
		buf[n+2] = byte(p.MediaWidth)
	}
	if p.Flag&PrintInformationLength != 0 {
		buf[n+3] = byte(p.MediaLength)
	}
	if p.Flag&(PrintInformationKind|PrintInformationWidth|PrintInformationLength) != 0 {
		commandEndian.PutUint32(buf[n+4:], uint32(p.RasterNunmber))
	}
	if !p.NotStartingPage {
		buf[n+8] = 0
	} else {
		buf[n+8] = 1
	}
	buf[n+9] = 0

	_, err := cb.Write(buf)
	return err
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
	return cb.cmd8(commandVariousModeSettings, byte(mode))
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
	return cb.cmd8(commandAdvancedModeSettings, byte(mode))
}

var errMarginAmountUnaccpectable = errors.New("margin amount unaccpectable")

// SpecifyMarginAmount specifies the amount of the margins.
func (cb CommandBuilder) SpecifyMarginAmount(dots int) error {
	return cb.cmd16(commandSpecifyMarginAmount, uint16(dots))
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
	return cb.cmd8(commandSelectCompressionMode, byte(mode))
}

var errDataTooLong = errors.New("data too long")

// RasterGraphicsTransfer transfers the specified number of bytes of data.
// ses <Raster Command Reference> P37
func (cb CommandBuilder) RasterGraphicsTransfer(data []byte) error {
	if len(data) > 0xFFFF {
		return errDataTooLong
	}

	buf := make([]byte, len(commandRasterGrapicsTransfer)+2+len(data))
	n := copy(buf, commandRasterGrapicsTransfer)
	commandEndian.PutUint16(buf[n:], uint16(len(data)))
	copy(buf[n+2:], data)

	_, err := cb.Write(buf)
	return err
}

// ZeroRasterGraphics fills raster line with 0 data.
func (cb CommandBuilder) ZeroRasterGraphics() error {
	return cb.cmd(commandZeroRasterGraphics)
}

// Print used as a print command at the end of pages other than the last page when multiple pages are printed.
func (cb CommandBuilder) Print() error {
	return cb.cmd(commandPrint)
}

// PrintWithFeeding used as a print command at the end of the last page.
func (cb CommandBuilder) PrintWithFeeding() error {
	return cb.cmd(commandPrintWithFeeding)
}

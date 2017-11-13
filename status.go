package braster

import "errors"
import "io"

type ModelCode byte

const (
	// ModelPTH500 PT-H500
	ModelPTH500 = ModelCode('d')
	// ModelPTE500 PT-E500
	ModelPTE500 = ModelCode('e')
	// ModelPTP700 PT-P700
	ModelPTP700 = ModelCode('g')
)

func (mc ModelCode) String() string {
	switch mc {
	case ModelPTH500:
		return "PT-H500"
	case ModelPTE500:
		return "PT-E500"
	case ModelPTP700:
		return "PT-P700"
	default:
		return "unknown"
	}
}

func (mc ModelCode) DPI() DPI {
	return DPI(180)
}
func (mc ModelCode) TotalDots() int {
	return 128
}

type ErrorInfomation int

const (
	NoError                 = ErrorInfomation(0)
	ErrorNoMedia            = ErrorInfomation(1 << 0)
	ErrorCutterJam          = ErrorInfomation(1 << 2)
	ErrorWeakBatteries      = ErrorInfomation(1 << 3)
	ErrorHighVoltageAdapter = ErrorInfomation(1 << 6)
	ErrorWrongMedia         = ErrorInfomation(1 << 8)
	ErrorCoverOpen          = ErrorInfomation(1 << 12)
	ErrorOverheating        = ErrorInfomation(1 << 13)
)

func (ei ErrorInfomation) String() string {
	switch ei {
	case NoError:
		return "No error"
	case ErrorNoMedia:
		return "No media"
	case ErrorCutterJam:
		return "Cutter jam"
	case ErrorWeakBatteries:
		return "Weak batteries"
	case ErrorHighVoltageAdapter:
		return "High-voltage adapter"
	case ErrorWrongMedia:
		return "Wrong media"
	case ErrorCoverOpen:
		return "Cover open"
	case ErrorOverheating:
		return "Overheating"
	default:
		return "Unknown"
	}
}

type MediaWidth int

const (
	NoTape        = MediaWidth(0)
	MediaWidth3_5 = MediaWidth(4)
	MediaWidth6   = MediaWidth(6)
	MediaWidth9   = MediaWidth(9)
	MediaWidth12  = MediaWidth(12)
	MediaWidth18  = MediaWidth(18)
	MediaWidth24  = MediaWidth(24)
)

func (mw MediaWidth) String() string {
	switch mw {
	case NoTape:
		return "No tape"
	case MediaWidth3_5:
		return "3.5mm"
	case MediaWidth6:
		return "6mm"
	case MediaWidth9:
		return "9mm"
	case MediaWidth12:
		return "12mm"
	case MediaWidth18:
		return "18mm"
	case MediaWidth24:
		return "24mm"
	default:
		return "Unknown"
	}
}

type MediaType int

const (
	NoMedia          = MediaType(0)
	LaminatedTape    = MediaType(0x01)
	NonlaminatedTape = MediaType(0x03)
	HeatShrinkTube   = MediaType(0x11)
	IncompatibleTape = MediaType(0xFF)
)

func (mt MediaType) String() string {
	switch mt {
	case NoMedia:
		return "No media"
	case LaminatedTape:
		return "Laminated tape"
	case NonlaminatedTape:
		return "Non-laminated tape"
	case HeatShrinkTube:
		return "Heat-Shrink Tube"
	case IncompatibleTape:
		return "Incompatible tape"
	}
	return "Unknown"
}
func (mt MediaType) IsValid() bool {
	switch mt {
	case LaminatedTape, NonlaminatedTape, HeatShrinkTube:
		return true
	}
	return false
}

type StatusType int

const (
	StatusReplyToStatusRequest = StatusType(0x00)
	StatusPrintingCompleted    = StatusType(0x01)
	StatusErrorOccurred        = StatusType(0x02)
	StatusExitIFMode           = StatusType(0x03) // Not used
	StatusTurnedOff            = StatusType(0x04)
	StatusNotification         = StatusType(0x05)
	StatusPhaseChange          = StatusType(0x06)
)

func (st StatusType) String() string {
	switch st {
	case StatusReplyToStatusRequest:
		return "Reply to status request "
	case StatusPrintingCompleted:
		return "Printing completed"
	case StatusErrorOccurred:
		return "Error occurred"
	case StatusExitIFMode:
		return "Exit IF mode "
	case StatusTurnedOff:
		return "Turned off"
	case StatusNotification:
		return "Notification"
	case StatusPhaseChange:
		return "Phase change"
	default:
		if 0x07 <= st && st <= 0x20 {
			return "(Not used)"
		}
		return "(Reserved)"
	}
}

type NotificationNumber int

const (
	NotificationNotAvailable = NotificationNumber(0x00)
	NotificationCoverOpen    = NotificationNumber(0x01)
	NotificationCoverClosed  = NotificationNumber(0x02)
)

type StatusInformation struct {
	Model            ModelCode
	ErrorInformation ErrorInfomation
	MediaWidth       MediaWidth
	MediaType        MediaType
	MediaLength      int
	Status           StatusType
	Notification     NotificationNumber
}

func (si StatusInformation) Media() Media {
	return RecognizeMedia(si.MediaType, si.MediaWidth)
}

var errStatusInformationDataTooShort = errors.New("status information data too short")

// ReadStatusInformation reads status information form bytes, bytes length must be 32 bytes.
func ReadStatusInformation(b []byte) (*StatusInformation, error) {
	if len(b) < 32 {
		return nil, errStatusInformationDataTooShort
	}
	si := StatusInformation{}
	si.Model = ModelCode(b[4])
	si.ErrorInformation = ErrorInfomation(b[8]) | (ErrorInfomation(b[9]) << 8)
	si.MediaWidth = MediaWidth(b[10])
	si.MediaType = MediaType(b[11])
	si.MediaLength = int(b[17])
	si.Status = StatusType(b[18])
	si.Notification = NotificationNumber(b[22])
	return &si, nil
}

// QueryStatusInformation query status information form io.
func QueryStatusInformation(rw io.ReadWriter) (*StatusInformation, error) {
	err := CommandBuilder{rw}.StatusInformationRequest()
	if err != nil {
		return nil, err
	}
	buf := [32]byte{}
	_, err = io.ReadFull(rw, buf[:])
	if err != nil {
		return nil, err
	}
	return ReadStatusInformation(buf[:])
}

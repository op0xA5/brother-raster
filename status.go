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

type ErrorInfomation int

const (
	ErrorNoMedia            = ErrorInfomation(1 << 0)
	ErrorCutterJam          = ErrorInfomation(1 << 2)
	ErrorWeakBatteries      = ErrorInfomation(1 << 3)
	ErrorHighVoltageAdapter = ErrorInfomation(1 << 6)
	ErrorWrongMedia         = ErrorInfomation(1 << 8)
	ErrorCoverOpen          = ErrorInfomation(1 << 12)
	ErrorOverheating        = ErrorInfomation(1 << 13)
)

type MediaWidth int

const (
	NoTape        = MediaType(0)
	MediaWidth3_5 = MediaWidth(4)
	MediaWidth6   = MediaWidth(6)
	MediaWidth9   = MediaWidth(9)
	MediaWidth12  = MediaWidth(12)
	MediaWidth18  = MediaWidth(18)
	MediaWidth24  = MediaWidth(24)
)

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

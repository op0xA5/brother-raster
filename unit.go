package braster

type DPI float32

const millimeterInch = 25.4

func (dpi DPI) MillimetreToDots(mm float32) int {
	return int(mm / millimeterInch * float32(dpi))
}
func (dpi DPI) DotsToMillimetre(d int) float32 {
	return float32(d) / float32(dpi) * millimeterInch
}

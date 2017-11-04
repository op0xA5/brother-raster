package braster

// MediaType supported media types defined at <Raster Command Reference> P14
type MediaType int

const (
	// TZeTape3_5 3.5mm TZe tape
	TZeTape3_5 = MediaType(263)
	// TZeTape6 6mm TZe tape
	TZeTape6 = MediaType(257)
	// TZeTape9 9mm TZe tape
	TZeTape9 = MediaType(258)
	// TZeTape12 12mm TZe tape
	TZeTape12 = MediaType(259)
	// TZeTape18 18mm TZe tape
	TZeTape18 = MediaType(260)
	// TZeTape24 24mm TZe tape
	TZeTape24 = MediaType(261)
	// HeatShrinkTube6 6mm Heat-Shrink Tube
	HeatShrinkTube6 = MediaType(415)
	// HeatShrinkTube9 9mm Heat-Shrink Tube
	HeatShrinkTube9 = MediaType(416)
	// HeatShrinkTube12 12mm Heat-Shrink Tube
	HeatShrinkTube12 = MediaType(417)
	// HeatShrinkTube18 18mm Heat-Shrink Tube
	HeatShrinkTube18 = MediaType(418)
	// HeatShrinkTube24 24mm Heat-Shrink Tube
	HeatShrinkTube24 = MediaType(419)
)

type mediaInfoType struct {
	mediaType                   MediaType
	size                        float32
	printAreaDots               int
	leftMarginDots              int
	rightMarginDots             int
	pageMarginDots              int
	bytesRasterGraphicsTransfer int
}

var mediaInfo = []mediaInfoType{
	mediaInfoType{TZeTape3_5, 3.5, 24, 52, 52, 0, 16},
	mediaInfoType{TZeTape6, 6, 32, 48, 48, 5, 16},
	mediaInfoType{TZeTape9, 9, 50, 39, 39, 7, 16},
	mediaInfoType{TZeTape12, 12, 70, 29, 29, 7, 16},
	mediaInfoType{TZeTape18, 18, 112, 8, 8, 8, 16},
	mediaInfoType{TZeTape24, 24, 128, 0, 0, 21, 16},
	mediaInfoType{HeatShrinkTube6, 6, 28, 50, 50, 6, 16},
	mediaInfoType{HeatShrinkTube9, 9, 48, 40, 40, 8, 16},
	mediaInfoType{HeatShrinkTube12, 12, 66, 31, 31, 8, 16},
	mediaInfoType{HeatShrinkTube18, 18, 106, 11, 11, 10, 16},
	mediaInfoType{HeatShrinkTube24, 24, 128, 0, 0, 20, 16},
}

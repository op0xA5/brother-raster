package braster

// Media supported media types defined at <Raster Command Reference> P14
type Media int

const (
	// UnknownMedia unknown
	UnknownMedia = Media(0)
	// MediaTZeTape3_5 3.5mm TZe tape
	MediaTZeTape3_5 = Media(263)
	// MediaTZeTape6 6mm TZe tape
	MediaTZeTape6 = Media(257)
	// MediaTZeTape9 9mm TZe tape
	MediaTZeTape9 = Media(258)
	// MediaTZeTape12 12mm TZe tape
	MediaTZeTape12 = Media(259)
	// MediaTZeTape18 18mm TZe tape
	MediaTZeTape18 = Media(260)
	// MediaTZeTape24 24mm TZe tape
	MediaTZeTape24 = Media(261)
	// MediaHeatShrinkTube6 6mm Heat-Shrink Tube
	MediaHeatShrinkTube6 = Media(415)
	// MediaHeatShrinkTube9 9mm Heat-Shrink Tube
	MediaHeatShrinkTube9 = Media(416)
	// MediaHeatShrinkTube12 12mm Heat-Shrink Tube
	MediaHeatShrinkTube12 = Media(417)
	// MediaHeatShrinkTube18 18mm Heat-Shrink Tube
	MediaHeatShrinkTube18 = Media(418)
	// MediaHeatShrinkTube24 24mm Heat-Shrink Tube
	MediaHeatShrinkTube24 = Media(419)
)

// RecognizeMedia from ReadStatusInformation or QueryStatusInformation results.
func RecognizeMedia(typ MediaType, width MediaWidth) Media {
	switch typ {
	case LaminatedTape, NonlaminatedTape:
		switch width {
		case MediaWidth3_5:
			return MediaTZeTape3_5
		case MediaWidth6:
			return MediaTZeTape6
		case MediaWidth9:
			return MediaTZeTape9
		case MediaWidth12:
			return MediaTZeTape12
		case MediaWidth18:
			return MediaTZeTape18
		case MediaWidth24:
			return MediaTZeTape24
		}
	case HeatShrinkTube:
		switch width {
		case MediaWidth6:
			return MediaHeatShrinkTube6
		case MediaWidth9:
			return MediaHeatShrinkTube9
		case MediaWidth12:
			return MediaHeatShrinkTube12
		case MediaWidth18:
			return MediaHeatShrinkTube18
		case MediaWidth24:
			return MediaHeatShrinkTube24
		}
	}
	return UnknownMedia
}

type mediaInfoType struct {
	media                       Media
	size                        float32
	printAreaDots               int
	leftMarginDots              int
	rightMarginDots             int
	pageMarginDots              int
	bytesRasterGraphicsTransfer int
	minimumMarginDots           int
	maximumMarginDots           int
	minimumMarginNoPrecutDots   int
}

var mediaInfo = []mediaInfoType{
	mediaInfoType{MediaTZeTape3_5, 3.5, 24, 52, 52, 0, 16, 14, 900, 172},
	mediaInfoType{MediaTZeTape6, 6, 32, 48, 48, 5, 16, 14, 900, 172},
	mediaInfoType{MediaTZeTape9, 9, 50, 39, 39, 7, 16, 14, 900, 172},
	mediaInfoType{MediaTZeTape12, 12, 70, 29, 29, 7, 16, 14, 900, 172},
	mediaInfoType{MediaTZeTape18, 18, 112, 8, 8, 8, 16, 14, 900, 172},
	mediaInfoType{MediaTZeTape24, 24, 128, 0, 0, 21, 16, 14, 900, 172},
	mediaInfoType{MediaHeatShrinkTube6, 6, 28, 50, 50, 6, 16, 14, 900, 172},
	mediaInfoType{MediaHeatShrinkTube9, 9, 48, 40, 40, 8, 16, 14, 900, 172},
	mediaInfoType{MediaHeatShrinkTube12, 12, 66, 31, 31, 8, 16, 14, 900, 172},
	mediaInfoType{MediaHeatShrinkTube18, 18, 106, 11, 11, 10, 16, 14, 900, 172},
	mediaInfoType{MediaHeatShrinkTube24, 24, 128, 0, 0, 20, 16, 14, 900, 172},
}

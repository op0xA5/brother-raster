package braster

import "sync"

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

type MediaInfo struct {
	Name              string
	Designation       float32
	PageSize          float32
	PrintArea         float32
	PageMargin        float32
	MinMargin         float32
	MaxMargin         float32
	MinMarginNoPrecut float32
	MinLength         float32
	MaxLength         float32
}

var mediaInfo = map[Media]*MediaInfo{
	MediaTZeTape3_5:       &MediaInfo{"3.5mm TZe tape", 3.5, 3.4, 3.4, 0, 2, 127, 24.3, 4.4, 1000},
	MediaTZeTape6:         &MediaInfo{"6mm TZe tape", 6, 5.9, 4.5, 0.7, 2, 127, 24.3, 4.4, 1000},
	MediaTZeTape9:         &MediaInfo{"9mm TZe tape", 9, 9, 7.1, 0.98, 2, 127, 24.3, 4.4, 1000},
	MediaTZeTape12:        &MediaInfo{"12mm TZe tape", 12, 11.9, 9.9, 0.98, 2, 127, 24.3, 4.4, 1000},
	MediaTZeTape18:        &MediaInfo{"18mm TZe tape", 18, 18.1, 15.8, 1.12, 2, 127, 24.3, 4.4, 1000},
	MediaTZeTape24:        &MediaInfo{"24mm TZe tape", 24, 24, 18.1, 2.96, 2, 127, 24.3, 4.4, 1000},
	MediaHeatShrinkTube6:  &MediaInfo{"6mm Heat-Shrink Tube", 6, 5.6, 3.9, 0.8, 2, 127, 24.3, 4.4, 500},
	MediaHeatShrinkTube9:  &MediaInfo{"9mm Heat-Shrink Tube", 9, 8.7, 6.8, 1.1, 2, 127, 24.3, 4.4, 500},
	MediaHeatShrinkTube12: &MediaInfo{"12mm Heat-Shrink Tube", 12, 11.6, 9.3, 1.1, 2, 127, 24.3, 4.4, 500},
	MediaHeatShrinkTube18: &MediaInfo{"18mm Heat-Shrink Tube", 18, 17.8, 14.9, 1.4, 2, 127, 24.3, 4.4, 500},
	MediaHeatShrinkTube24: &MediaInfo{"24mm Heat-Shrink Tube", 24, 23.7, 18.1, 2.8, 2, 127, 24.3, 4.4, 500},
}
var invalidMedia = &MediaInfo{Name: "Unknown"}
var mediaInfoMu = new(sync.RWMutex)

func (m Media) MediaInfo() *MediaInfo {
	mediaInfoMu.RLock()
	defer mediaInfoMu.RUnlock()

	mi := mediaInfo[m]
	if mi == nil {
		return invalidMedia
	}
	return mi
}

func (m Media) String() string {
	return m.MediaInfo().Name
}

func RegisterMediaInfo(m Media, mi *MediaInfo) {
	mediaInfoMu.Lock()
	defer mediaInfoMu.Unlock()

	mediaInfo[m] = mi
}

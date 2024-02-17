package world

const (
	GroupCategory = "Group Category"
	CategoryOcean = "Category Ocean"
	CategoryLand  = "Category Land"
)

const (
	GroupTemperature     = "Group Temperature"
	TemperatureWarm      = "Temperature Warm"
	TemperatureTemperate = "Temperature Temperate"
	TemperatureCold      = "Temperature Cold"
	TemperatureFreezing  = "Temperature Freezing"
)

const (
	GroupPrecipitation    = "Group Precipitation"
	PrecipitationDry      = "Precipitation Dry"
	PrecipitationSlight   = "Precipitation Slight"
	PrecipitationModerate = "Precipitation Moderate"
	PrecipitationHeavy    = "Precipitation Heavy"
)

const (
	GroupAltitude = "Group Altitude"
	HeightLow     = "Height Low"
	HeightHigh    = "Height High"
	DepthLow      = "Depth Low"
	DepthHigh     = "Depth High"
)

const (
	GroupVariant   = "Group Variant"
	VariantSpecial = "VariantSpecial"
)

const (
	GroupUnknown = "Group Unknown"
)

type (
	Label      string
	LabelGroup string
)

func (l Label) Group() LabelGroup {
	switch l {
	case CategoryOcean, CategoryLand:
		return GroupCategory
	case TemperatureWarm, TemperatureTemperate, TemperatureCold, TemperatureFreezing:
		return GroupTemperature
	case PrecipitationDry, PrecipitationSlight, PrecipitationModerate, PrecipitationHeavy:
		return GroupPrecipitation
	case HeightLow, HeightHigh, DepthLow, DepthHigh:
		return GroupAltitude
	case VariantSpecial:
		return GroupVariant
	default:
		return GroupUnknown
	}
}

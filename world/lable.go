package world

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	GroupCategory LabelGroup = "Group Category"

	CategoryOcean Label = "Category Ocean"
	CategoryRiver Label = "Category River"
	CategoryLand  Label = "Category Land"
)

const (
	GroupMacroTemperature LabelGroup = "Group Macro Temperature"

	MacroTemperatureWarm      Label = "Macro Temperature Warm"
	MacroTemperatureTemperate Label = "Macro Temperature Temperate"
	MacroTemperatureCold      Label = "Macro Temperature Cold"
	MacroTemperatureFreezing  Label = "Macro Temperature Freezing"
)

const (
	GroupTemperature LabelGroup = "Group Temperature"

	Temperature7 Label = "Temperature 7" // Warm
	Temperature6 Label = "Temperature 6" // Warm
	Temperature5 Label = "Temperature 5" // Temperate
	Temperature4 Label = "Temperature 4" // Temperate
	Temperature3 Label = "Temperature 3" // Cold
	Temperature2 Label = "Temperature 2" // Cold
	Temperature1 Label = "Temperature 1" // Freezing
	Temperature0 Label = "Temperature 0" // Freezing
)

const (
	GroupMacroHumidity LabelGroup = "Group Macro Humidity"

	MacroHumidityHigh     Label = "Macro Humidity High"
	MacroHumidityModerate Label = "Macro Humidity Moderate"
	MacroHumidityLow      Label = "Macro Humidity Low"
	MacroHumidityMinimal  Label = "Macro Humidity Minimal"
)

const (
	GroupHumidity LabelGroup = "Group Humidity"

	Humidity7 Label = "Humidity 7" // High
	Humidity6 Label = "Humidity 6" // High
	Humidity5 Label = "Humidity 5" // Moderate
	Humidity4 Label = "Humidity 4" // Moderate
	Humidity3 Label = "Humidity 3" // Low
	Humidity2 Label = "Humidity 2" // Low
	Humidity1 Label = "Humidity 1" // Minimal
	Humidity0 Label = "Humidity 0" // Minimal
)

const (
	GroupMacroAltitude LabelGroup = "Group Macro Altitude"

	MacroHeightLow      Label = "Height Low"
	MacroHeightMedium   Label = "Height Medium"
	MacroHeightHigh     Label = "Height High"
	MacroHeightVeryHigh Label = "Height Very High"
	MacroDepthLow       Label = "Depth Low"
	MacroDepthHigh      Label = "Depth High"
)

const (
	GroupAltitude LabelGroup = "Group Altitude"

	Height15 Label = "Height 15" // Height Very High
	Height14 Label = "Height 14" // Height Very High
	Height13 Label = "Height 13" // Height Very High
	Height12 Label = "Height 12" // Height Very High
	Height11 Label = "Height 11" // Height High
	Height10 Label = "Height 10" // Height High
	Height9  Label = "Height 9"  // Height Medium
	Height8  Label = "Height 8"  // Height Medium
	Height7  Label = "Height 7"  // Height Medium
	Height6  Label = "Height 6"  // Height Medium
	Height5  Label = "Height 5"  // Height Medium
	Height4  Label = "Height 4"  // Height Low
	Height3  Label = "Height 3"  // Height Low
	Height2  Label = "Height 2"  // Height Low
	Height1  Label = "Height 1"  // Height Low
	Height0  Label = "Height 0"  // Height Low
	Depth0   Label = "Depth 0"   // Depth Low
	Depth1   Label = "Depth 1"   // Depth High
)

const (
	GroupVariant LabelGroup = "Group Variant"

	VariantSpecial  Label = "Variant Special"
	VariantCollinar Label = "Variant Collinar"
)

const (
	GroupUnknown LabelGroup = "Group Unknown"
)

type (
	Label      string
	LabelGroup string
)

func (l Label) Group() LabelGroup {
	switch l {
	case CategoryOcean, CategoryRiver, CategoryLand:
		return GroupCategory
	case MacroTemperatureWarm, MacroTemperatureTemperate, MacroTemperatureCold, MacroTemperatureFreezing:
		return GroupMacroTemperature
	case Temperature7, Temperature6, Temperature5, Temperature4, Temperature3, Temperature2, Temperature1, Temperature0:
		return GroupTemperature
	case MacroHumidityMinimal, MacroHumidityLow, MacroHumidityModerate, MacroHumidityHigh:
		return GroupMacroHumidity
	case Humidity7, Humidity6, Humidity5, Humidity4, Humidity3, Humidity2, Humidity1, Humidity0:
		return GroupHumidity
	case MacroHeightLow, MacroHeightMedium, MacroHeightHigh, MacroHeightVeryHigh, MacroDepthLow, MacroDepthHigh:
		return GroupMacroAltitude
	case Height15, Height14, Height13, Height12, Height11, Height10, Height9, Height8, Height7, Height6, Height5, Height4, Height3, Height2, Height1, Height0, Depth0, Depth1:
		return GroupAltitude
	case VariantSpecial, VariantCollinar:
		return GroupVariant
	default:
		return GroupUnknown
	}
}

func (l Label) String() string {
	return string(l)
}

func (l Label) Value() (int, error) {
	switch l.Group() {
	case GroupTemperature, GroupHumidity, GroupAltitude:
		valStrSlice := strings.Split(l.String(), " ")
		return strconv.Atoi(valStrSlice[len(valStrSlice)-1])
	default:
		return 0, fmt.Errorf("this labels does not support label.Value()")
	}
}

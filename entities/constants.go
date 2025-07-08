package entities

const (
	PresetScreen   string = "screen"
	PresetEbook    string = "ebook"
	PresetPrinter  string = "printer"
	PresetPrepress string = "prepress"
)

type ResolutionLimits struct {
	ColorImageMin float64
	ColorImageMax float64
	MonoImageMin  float64
	MonoImageMax  float64
	GrayImageMin  float64
	GrayImageMax  float64
}

var DefaultResolutionLimits = map[string]ResolutionLimits{
	PresetScreen: {
		ColorImageMax: 72,
		MonoImageMax:  300,
		GrayImageMax:  72,
	},
	PresetEbook: {
		ColorImageMax: 150,
		MonoImageMax:  300,
		GrayImageMax:  150,
	},
	PresetPrinter: {
		ColorImageMax: 300,
		MonoImageMax:  1200,
		GrayImageMax:  300,
	},
	PresetPrepress: {
		ColorImageMax: 300,
		MonoImageMax:  1200,
		GrayImageMax:  300,
	},
}

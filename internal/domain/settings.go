package domain

type Settings struct {
	Locale   LocaleSettings
	Theme    ThemeSettings
	Controls ControlsSettings
}

type LocaleSettings struct {
	Language string
	Region   string
}

type ThemeSettings struct {
	ColourScheme string
	Mode         string
}

type ControlsSettings struct {
	Vim               bool
	DoubleClickToEdit bool
}

var DefaultSettings = Settings{
	Locale: LocaleSettings{
		Language: "en",
		Region:   "gb",
	},
	Theme: ThemeSettings{
		ColourScheme: "default",
		Mode:         "auto",
	},
	Controls: ControlsSettings{
		Vim:               false,
		DoubleClickToEdit: false,
	},
}
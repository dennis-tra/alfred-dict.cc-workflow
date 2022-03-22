package workflow

type Language string

const (
	LanguageEnglish    Language = "english"
	LanguageGerman     Language = "german"
	LanguageFrench     Language = "french"
	LanguageSwedish    Language = "swedish"
	LanguageSpanish    Language = "spanish"
	LanguageBulgarian  Language = "bulgarian"
	LanguageRomanian   Language = "romanian"
	LanguageItalian    Language = "italian"
	LanguagePortuguese Language = "portuguese"
	LanguageRussian    Language = "russian"
)

var (
	LangMap = map[string]Language{
		"en":                       LanguageEnglish,
		"eng":                      LanguageEnglish,
		"de":                       LanguageGerman,
		"ger":                      LanguageGerman,
		"fr":                       LanguageFrench,
		"fra":                      LanguageFrench,
		"sv":                       LanguageSwedish,
		"swe":                      LanguageSwedish,
		"es":                       LanguageSpanish,
		"esp":                      LanguageSpanish,
		"bg":                       LanguageBulgarian,
		"bul":                      LanguageBulgarian,
		"ro":                       LanguageRomanian,
		"rom":                      LanguageRomanian,
		"it":                       LanguageItalian,
		"ita":                      LanguageItalian,
		"pt":                       LanguagePortuguese,
		"por":                      LanguagePortuguese,
		"ru":                       LanguageRussian,
		"rus":                      LanguageRussian,
		string(LanguageEnglish):    LanguageEnglish,
		string(LanguageGerman):     LanguageGerman,
		string(LanguageFrench):     LanguageFrench,
		string(LanguageSwedish):    LanguageSwedish,
		string(LanguageSpanish):    LanguageSpanish,
		string(LanguageBulgarian):  LanguageBulgarian,
		string(LanguageRomanian):   LanguageRomanian,
		string(LanguageItalian):    LanguageItalian,
		string(LanguagePortuguese): LanguagePortuguese,
		string(LanguageRussian):    LanguageRussian,
	}

	LangSubdomainMap = map[Language]string{
		LanguageEnglish:    "en",
		LanguageGerman:     "de",
		LanguageFrench:     "fr",
		LanguageSwedish:    "sv",
		LanguageSpanish:    "es",
		LanguageBulgarian:  "bg",
		LanguageRomanian:   "ro",
		LanguageItalian:    "it",
		LanguagePortuguese: "pt",
		LanguageRussian:    "ru",
	}
)

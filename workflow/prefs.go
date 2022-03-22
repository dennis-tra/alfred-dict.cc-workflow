package workflow

import (
	"fmt"
	"strings"
)

// Env is used to capture environment variables.
// These variables indicate the default language pair the user wants to translate between.
type Env struct {
	FromLanguage string `env:"from_language"`
	ToLanguage   string `env:"to_language"`
}

// Preferences tracks the user preferences about the language pairs. The languages
// can either be:
//   1. the default pair (english, german)
//   2. set in the workflow environment
//   3. given in the query string
type Preferences struct {
	FromLanguage Language
	ToLanguage   Language
}

// NewDefaultPreferences initializes a new Preferences struct with the default language pair.
func NewDefaultPreferences() *Preferences {
	return &Preferences{
		FromLanguage: LanguageEnglish,
		ToLanguage:   LanguageGerman,
	}
}

// Apply applies the environment variable to the current preference object.
func (p *Preferences) Apply(env *Env) {
	fromLanguage, found := LangMap[env.FromLanguage]
	if found {
		p.FromLanguage = fromLanguage
	}

	toLanguage, found := LangMap[env.ToLanguage]
	if found {
		p.ToLanguage = toLanguage
	}
}

// Parse parses the query string to check if the user wants to translate between a custom language pair
func (p *Preferences) Parse(args []string) []string {
	if len(args) < 2 {
		return args
	}

	fromLang, found := LangMap[args[0]]
	if !found {
		return args
	}

	toLang, found := LangMap[args[1]]
	if !found {
		return args
	}

	p.FromLanguage = fromLang
	p.ToLanguage = toLang

	return args[2:]
}

// Subdomain constructs the subdomain to use for the query.
func (p *Preferences) Subdomain() string {
	fromSubdomain, found := LangSubdomainMap[p.FromLanguage]
	if !found {
		return "www"
	}

	toSubdomain, found := LangSubdomainMap[p.ToLanguage]
	if !found {
		return "www"
	}

	return fromSubdomain + toSubdomain
}

func (p *Preferences) String() string {
	return fmt.Sprintf("%s <-> %s", strings.Title(string(p.FromLanguage)), strings.Title(string(p.ToLanguage)))
}

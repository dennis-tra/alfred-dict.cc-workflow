package workflow

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPreferences_Parse(t *testing.T) {
	defaultPrefs := NewDefaultPreferences()
	tests := []struct {
		name             string
		args             []string
		wantFromLanguage Language
		wantToLanguage   Language
		wantQuery        []string
	}{
		{
			name:             "empty args",
			args:             []string{},
			wantFromLanguage: defaultPrefs.FromLanguage,
			wantToLanguage:   defaultPrefs.ToLanguage,
			wantQuery:        []string{},
		},
		{
			name:             "mixed abbreviations",
			args:             []string{"it", "ger"},
			wantFromLanguage: LanguageItalian,
			wantToLanguage:   LanguageGerman,
			wantQuery:        []string{},
		},
		{
			name:             "normal query",
			args:             []string{"test"},
			wantFromLanguage: defaultPrefs.FromLanguage,
			wantToLanguage:   defaultPrefs.ToLanguage,
			wantQuery:        []string{"test"},
		},
		{
			name:             "mixed validity query",
			args:             []string{"ro", "invalid", "query"},
			wantFromLanguage: defaultPrefs.FromLanguage,
			wantToLanguage:   defaultPrefs.ToLanguage,
			wantQuery:        []string{"ro", "invalid", "query"},
		},
		{
			name:             "only mapping",
			args:             []string{"ro", "it"},
			wantFromLanguage: LanguageRomanian,
			wantToLanguage:   LanguageItalian,
			wantQuery:        []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewDefaultPreferences()
			query := p.Parse(tt.args)
			assert.Equal(t, p.FromLanguage, tt.wantFromLanguage)
			assert.Equal(t, p.ToLanguage, tt.wantToLanguage)
			require.Equal(t, len(query), len(tt.wantQuery))
			for i := range query {
				assert.Equal(t, query[i], tt.wantQuery[i])
			}
		})
	}
}

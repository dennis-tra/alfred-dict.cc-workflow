package workflow

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	aw "github.com/deanishe/awgo"
)

// ErrNotFound is the error that gets returned if no results were found (alert: captain obvious)
var ErrNotFound = fmt.Errorf("translation not found")

// Dictcc hold workflow information and preferences.
type Dictcc struct {
	Workflow    *aw.Workflow
	Preferences *Preferences
}

// NewDictcc initializes a new Dictcc struct with the given workflow instance and default preferences.
func NewDictcc(wf *aw.Workflow) *Dictcc {
	return &Dictcc{
		Workflow:    wf,
		Preferences: NewDefaultPreferences(),
	}
}

// Run is called from the workflow instance and starts the execution of this workflow.
func (dcc *Dictcc) Run() {
	// Merge Alfred Workflow environment variables with preferences struct
	env := &Env{}
	if err := dcc.Workflow.Config.To(env); err != nil {
		dcc.Workflow.FatalError(err)
	}
	dcc.Preferences.Apply(env)

	query := dcc.Workflow.Args()
	// the query is give as one command line argument to the workflow
	if len(query) == 1 {
		dcc.HandleQuery(strings.Split(query[0], " "))
	} else {
		dcc.Workflow.FatalError(fmt.Errorf("unexpected length of arguments"))
	}
}

// HandleQuery takes the users' arguments and tries to make sense of them
func (dcc *Dictcc) HandleQuery(args []string) {
	// Try to parse language arguments from the user and merge them with preferences struct
	// args will be stripped by the language arguments and will only contain the translation query
	args = dcc.Preferences.Parse(args)

	// Join args to query string
	query := strings.Join(args, " ")
	query = strings.TrimSpace(query)

	// Check for empty query
	if query == "" {
		dcc.Workflow.NewItem("Translate...").Subtitle(dcc.Preferences.String())
		dcc.Workflow.SendFeedback()
		return
	}

	// Query dict.cc for results
	body, err := dcc.queryDictcc(query)
	if err != nil {
		dcc.Workflow.FatalError(err)
	}

	// Extract results from body for language 1
	resultsLang1, err := getResults(1, body)
	if err != nil {
		dcc.handleError(err, query)
	}

	// Extract results from body for language 2
	resultsLang2, err := getResults(2, body)
	if err != nil {
		dcc.handleError(err, query)
	}

	// Handle case where the result sets have different lengths (should not happen)
	maxIdx := len(resultsLang1)
	if len(resultsLang2) < maxIdx {
		maxIdx = len(resultsLang2)
	}

	// filtered results slices
	filResLang1 := []string{}
	filResLang2 := []string{}
	for i := 0; i < maxIdx; i++ {
		if resultsLang1[i] == "" || resultsLang2[i] == "" {
			continue
		}
		filResLang1 = append(filResLang1, resultsLang1[i])
		filResLang2 = append(filResLang2, resultsLang2[i])
	}
	maxIdx = len(filResLang1)
	resultsLang1 = filResLang1
	resultsLang2 = filResLang2

	// Apply heuristic to identify which results occurred more often
	occurrencesLang1 := 0
	occurrencesLang2 := 0
	for i := 0; i < maxIdx; i++ {
		if strings.ToLower(resultsLang1[i]) == strings.ToLower(query) {
			occurrencesLang1 += 1
		}
		if strings.ToLower(resultsLang2[i]) == strings.ToLower(query) {
			occurrencesLang2 += 1
		}
	}

	if occurrencesLang2 < occurrencesLang1 {
		dcc.prepareResults(query, resultsLang1, resultsLang2, maxIdx)
	} else {
		dcc.prepareResults(query, resultsLang2, resultsLang1, maxIdx)
	}

	dcc.Workflow.SendFeedback()
}

// handleError displays a user-friendly Not-Found message and the actual error otherwise.
func (dcc *Dictcc) handleError(err error, query string) {
	if !errors.Is(err, ErrNotFound) {
		dcc.Workflow.FatalError(err)
	}

	dcc.Workflow.NewItem(fmt.Sprintf("%q not found", query)).Subtitle(dcc.Preferences.String())
	dcc.Workflow.SendFeedback()
}

// prepareResults configures the result items.
func (dcc *Dictcc) prepareResults(query string, fromResults []string, toResults []string, maxIdx int) {
	for i := 0; i < maxIdx; i++ {
		it := dcc.Workflow.
			NewItem(toResults[i]).
			Subtitle(fromResults[i]).
			Valid(true).
			Arg(toResults[i])

		dcc.addMod(it.Cmd(), query)
		dcc.addMod(it.Alt(), query)
	}
}

// addMod adds alternative metadata for the given modifier.
func (dcc *Dictcc) addMod(mod *aw.Modifier, query string) {
	u := dcc.dictccURL(query)
	mod.Subtitle(fmt.Sprintf("Open dict.cc for %q in the browser...", query)).
		Valid(true).
		Arg(u.String())
}

// queryDictcc does an HTTP GET to dict.cc (with varying subdomains depending on the language pair) and
// parses the HTML body to a string.
func (dcc *Dictcc) queryDictcc(query string) (string, error) {
	// Generate URL
	u := dcc.dictccURL(query)

	// Actually query dictcc
	res, err := http.Get(u.String())
	if err != nil {
		return "", err
	}

	// Read complete HTML content
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// dictccURL constructs the URL that we request.
func (dcc *Dictcc) dictccURL(query string) url.URL {
	q := url.Values{}
	q.Set("s", query)

	return url.URL{
		Scheme:   "https",
		Host:     dcc.Preferences.Subdomain() + ".dict.cc",
		RawQuery: q.Encode(),
	}
}

// getResults extracts the translation results from the website. There are two arrays that contain the results:
//   1. `c1Arr`
//   2. `c2Arr`
// The function finds these lines extracts the javascript array content and parses the content as a CSV line.
// The last step is done to handle cases where the results also contain unescaped "," which would make splitting
// the string by "," harder.
func getResults(lang int, body string) ([]string, error) {
	re := regexp.MustCompile(`var c` + strconv.Itoa(lang) + `Arr = new Array\((.*)\);`)
	matches := re.FindStringSubmatch(body)
	if matches == nil || len(matches) != 2 {
		return nil, ErrNotFound
	}

	rows, err := csv.NewReader(strings.NewReader(matches[1])).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read csv line")
	}

	if len(rows) != 1 {
		return nil, fmt.Errorf("is not one csv line")
	}

	results := make([]string, len(rows[0]))
	for i, result := range rows[0] {
		results[i] = strings.ReplaceAll(result, `\'`, "'")
	}

	return results, nil
}

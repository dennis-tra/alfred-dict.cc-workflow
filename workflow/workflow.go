package workflow

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
	if len(args) == 0 {
		dcc.Workflow.NewItem("Translate...").Subtitle(dcc.Preferences.String())
		dcc.Workflow.SendFeedback()
		return
	}

	// Join args to query string
	query := strings.Join(args, " ")

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

	// Apply heuristic to identify which results occurred more often
	occurrencesLang1 := 0
	occurrencesLang2 := 0
	for i := 0; i < len(resultsLang1); i++ {
		if strings.ToLower(resultsLang1[i]) == strings.ToLower(query) {
			occurrencesLang1 += 1
		}
		if strings.ToLower(resultsLang2[i]) == strings.ToLower(query) {
			occurrencesLang2 += 1
		}
	}

	if occurrencesLang2 < occurrencesLang1 {
		dcc.sendResults(resultsLang1, resultsLang2)
	} else {
		dcc.sendResults(resultsLang2, resultsLang1)
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

// sendResults sends the translation results back to Alfred.
func (dcc *Dictcc) sendResults(fromResults []string, toResults []string) {
	maximumIdx := len(toResults)
	if len(fromResults) < maximumIdx {
		maximumIdx = len(fromResults)
	}
	for i := 0; i < maximumIdx; i++ {
		dcc.Workflow.
			NewItem(toResults[i]).
			Subtitle(fromResults[i]).
			Valid(true).
			Arg(toResults[i])
	}
}

// queryDictcc does an HTTP GET to dict.cc (with varying subdomains depending on the language pair) and
// parses the HTML body to a string.
func (dcc *Dictcc) queryDictcc(query string) (string, error) {
	log.Println("Escaping query string", query)
	q := url.Values{}
	q.Set("s", query)

	u := &url.URL{
		Scheme:   "https",
		Host:     dcc.Preferences.Subdomain() + ".dict.cc",
		RawQuery: q.Encode(),
	}

	log.Println("HTTP GET " + u.String())
	res, err := http.Get(u.String())
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
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

	results := []string{}
	for _, word := range rows[0] {
		if word == "" {
			continue
		}
		results = append(results, word)
	}
	return results, nil
}

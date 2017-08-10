package conversation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Intent struct {
	Name string

	matcherString string
	Matcher       *regexp.Regexp
}

func (i Intent) Process(input string) map[string]string {
	input = removeUselessLinkingWord(input)
	input = removeUselessPronouns(input)
	input = removeDiacritics(input)

	var m []string
	if m = i.Matcher.FindStringSubmatch(input); m == nil {
		return nil
	}

	out := make(map[string]string)

	for i, key := range i.Matcher.SubexpNames()[1:] {
		if strings.TrimSpace(m[i+1]) != "" {
			out[key] = strings.TrimSpace(m[i+1])
		}
	}

	return out
}

var PARTITIVE_ARTICLES []string = []string{"du ", "de la ", "d\\'", "de l\\'", "de ", "des "}
var DEFINITIVE_ARTICLES []string = []string{"le ", "la ", "les ", "l\\'"}
var SUBJECT_PRONOUN []string = []string{"tu ", "vous ", "t\\'"}
var FIRST_GROUP_TERM []string = []string{"e", "es", "ez", "iez", "er", "ais"}
var SECOND_GROUP_TERM []string = []string{"s", "ts", "tes", "r", "t"}

func buildFirstGroupTerminaison(input string) string {
	term := strings.Join(FIRST_GROUP_TERM, "|")

	repl := regexp.MustCompile("1G\\((?P<verb>\\w+)\\)")

	matches := repl.FindAllStringSubmatch(input, -1)
	if matches == nil {
		return input
	}

	for _, m := range matches {
		verb := m[1]
		stripVerb := verb[:len(verb)-2]
		input = strings.Replace(input, fmt.Sprintf("1G(%s)", verb), fmt.Sprintf("%s(?:%s)", stripVerb, term), 1)
	}

	return input
}

func buildSecondGroupTerminaison(input string) string {
	term := strings.Join(SECOND_GROUP_TERM, "|")

	repl := regexp.MustCompile("2G\\((?P<verb>\\w+)\\)")

	matches := repl.FindAllStringSubmatch(input, -1)
	if matches == nil {
		return input
	}

	for _, m := range matches {
		verb := m[1]
		stripVerb := verb[:len(verb)-2]
		input = strings.Replace(input, fmt.Sprintf("2G(%s)", verb), fmt.Sprintf("%s(?:%s)", stripVerb, term), 1)
	}

	return input
}

func buildPartitiveArticles(input string) string {
	return strings.Replace(input, "DE ", fmt.Sprintf("(?:%s)", strings.Join(PARTITIVE_ARTICLES, "|")), -1)
}

func buildDefinitiveArticles(input string) string {
	var s string = input

	s = strings.Replace(s, "LE? ", fmt.Sprintf("(?:%s)?", strings.Join(DEFINITIVE_ARTICLES, "|")), -1)
	s = strings.Replace(s, "LE ", fmt.Sprintf("(?:%s)", strings.Join(DEFINITIVE_ARTICLES, "|")), -1)

	return s
}

func buildSubjectPronoun(input string) string {
	var s string = input

	s = strings.Replace(s, "TU? ", fmt.Sprintf("(?P<pronoun>%s)?", strings.Join(SUBJECT_PRONOUN, "|")), -1)
	s = strings.Replace(s, "TU ", fmt.Sprintf("(?P<pronoun>%s)", strings.Join(SUBJECT_PRONOUN, "|")), -1)

	return s
}

func buildSubject(input string) string {
	return strings.Replace(input, "SUBJECT", "(?P<subject>\\w+)", -1)
}

func addInsensitive(input string) string {
	return fmt.Sprintf("(?i)%s", input)
}

func removeUselessLinkingWord(input string) string {
	var s string = input

	s = strings.Replace(s, " donc ", " ", -1)

	return s
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func removeDiacritics(input string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, input)

	return result
}

func removeUselessPronouns(input string) string {
	var s string = input

	s = strings.Replace(s, " me ", " ", -1)
	s = strings.Replace(s, " lui ", " ", -1)
	s = strings.Replace(s, " moi ", " ", -1)

	return s
}

func NewIntent(name string, matcherString string) Intent {
	intent := Intent{Name: name}

	intent.matcherString = buildPartitiveArticles(matcherString)
	intent.matcherString = buildDefinitiveArticles(intent.matcherString)
	intent.matcherString = buildSubjectPronoun(intent.matcherString)
	intent.matcherString = buildSubject(intent.matcherString)
	intent.matcherString = addInsensitive(intent.matcherString)
	intent.matcherString = buildFirstGroupTerminaison(intent.matcherString)
	intent.matcherString = buildSecondGroupTerminaison(intent.matcherString)
	intent.matcherString = removeDiacritics(intent.matcherString)

	intent.Matcher = regexp.MustCompile(intent.matcherString)

	return intent
}

func NewIntentCollection(intentMatchers [][2]string) []Intent {
	out := make([]Intent, len(intentMatchers))

	for i, intentMatcher := range intentMatchers {
		out[i] = NewIntent(intentMatcher[0], intentMatcher[1])
	}

	return out
}

func FindIntent(collection []Intent, input string) (string, map[string]string) {
	for _, intent := range collection {
		if result := intent.Process(input); result != nil {
			return intent.Name, result
		}
	}

	return "", nil
}

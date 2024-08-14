package linguist

import (
	"bufio"
	"bytes"
	"gopkg.in/yaml.v3"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	extensions   = map[string][]string{}
	filenames    = map[string][]string{}
	interpreters = map[string][]string{}
	colors       = map[string]string{}
	types        = map[string]string{}
	group        = map[string]string{}

	vendorRE        *regexp.Regexp
	typeRE          *regexp.Regexp
	doxRE           *regexp.Regexp
	shebangRE       = regexp.MustCompile(`^#!\s*(\S+)(?:\s+(\S+))?.*`)
	scriptVersionRE = regexp.MustCompile(`((?:\d+\.?)+)`)
)

func init() {

	var reg []string
	bytes := []byte(files["data/excluded.yml"])
	if err := yaml.Unmarshal(bytes, &reg); err != nil {
		log.Fatal(err)
		return
	}
	typeRE = regexp.MustCompile(strings.Join(reg, "|"))

	var regexps []string
	bytes = []byte(files["data/vendor.yml"])
	if err := yaml.Unmarshal(bytes, &regexps); err != nil {
		log.Fatal(err)
		return
	}
	vendorRE = regexp.MustCompile(strings.Join(regexps, "|"))

	var moreregex []string
	bytes = []byte(files["data/documentation.yml"])
	if err := yaml.Unmarshal(bytes, &moreregex); err != nil {
		log.Fatal(err)
		return
	}
	doxRE = regexp.MustCompile(strings.Join(moreregex, "|"))

	type language struct {
		Extensions   []string `yaml:"extensions,omitempty"`
		Filenames    []string `yaml:"filenames,omitempty"`
		Interpreters []string `yaml:"interpreters,omitempty"`
		Color        string   `yaml:"color,omitempty"`
		Type         string   `yaml:"type,omitempty"`
		Group        string   `yaml:"group,omitempty"`
	}
	languages := map[string]*language{}
	bytes = []byte(files["data/languages.yml"])
	if err := yaml.Unmarshal(bytes, languages); err != nil {
		log.Fatal(err)
	}

	for n, l := range languages {
		for _, e := range l.Extensions {
			extensions[e] = append(extensions[e], n)
		}
		for _, f := range l.Filenames {
			filenames[f] = append(filenames[f], n)
		}
		for _, i := range l.Interpreters {
			interpreters[i] = append(interpreters[i], n)
		}
		group[n] = l.Group
		types[n] = l.Type
		colors[n] = l.Color
	}
}

func LanguageColor(language string) string {
	if c, ok := colors[language]; ok {
		return c
	}
	return ""
}

func LanguageGroup(x string) string {
	if c, ok := group[x]; ok {
		return c
	}
	return ""
}

func LanguageType(x string) (string, bool) {
	if c, ok := types[x]; ok {
		if isTypeAllowed(c) {
			return c, true
		} else {
			return "", false
		}
	}
	return "", false
}

func LanguageByFilename(filename string) string {
	if l := filenames[filename]; len(l) == 1 {
		return l[0]
	}
	ext := filepath.Ext(filename)
	if ext != "" {
		if l := extensions[ext]; len(l) == 1 {
			return l[0]
		}
	}
	return ""
}

func LanguageHints(filename string) (hints []string) {
	if l, ok := filenames[filename]; ok {
		hints = append(hints, l...)
	}
	if ext := filepath.Ext(filename); ext != "" {
		if l, ok := extensions[ext]; ok {
			hints = append(hints, l...)
		}
	}
	return hints
}

func LanguageByContents(contents []byte, hints []string) string {
	interpreter := detectInterpreter(contents)
	if interpreter != "" {
		if l := interpreters[interpreter]; len(l) == 1 {
			return l[0]
		}
	}
	return Analyse(contents, hints)
}

func ShouldIgnoreFilename(filename string) bool {
	vendored := IsVendored(filename)
	documentation := IsDocumentation(filename)
	return vendored || documentation
	// return IsVendored(filename) || IsDocumentation(filename)
}

func ShouldIgnoreContents(contents []byte) bool {
	return IsBinary(contents)
}

func isTypeAllowed(path string) bool {
	return typeRE.MatchString(path)
}

func IsVendored(path string) bool {
	return vendorRE.MatchString(path)
}

func IsDocumentation(path string) bool {
	return doxRE.MatchString(path)
}
func detectInterpreter(contents []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(contents))
	scanner.Scan()
	line := scanner.Text()
	m := shebangRE.FindStringSubmatch(line)
	if m == nil || len(m) != 3 {
		return ""
	}
	base := filepath.Base(m[1])
	if base == "env" && m[2] != "" {
		base = m[2]
	}
	// Strip suffixed version number.
	return scriptVersionRE.ReplaceAllString(base, "")
}
func IsBinary(contents []byte) bool {
	// NOTE(tso): preliminary testing on this method of checking for binary
	// contents were promising, having fed a document consisting of all
	// utf-8 codepoints from 0000 to FFFF with satisfactory results. Thanks
	// to robpike.io/cmd/unicode:
	// ```
	// unicode -c $(seq 0 65535 | xargs printf "%04x ") | tr -d '\n' > unicode_test
	// ```
	//
	// However, the intentional presence of character escape codes to throw
	// this function off is entirely possible, as is, potentially, a binary
	// file consisting entirely of the 4 exceptions to the rule for the first
	// 512 bytes. It is also possible that more character escape codes need
	// to be added.
	//
	// Further analysis and real world testing of this is required.
	for n, b := range contents {
		if n >= 512 {
			break
		}
		if b < 32 {
			switch b {
			case 0:
				fallthrough
			case 9:
				fallthrough
			case 10:
				fallthrough
			case 13:
				continue
			default:
				return true
			}
		}
	}
	return false
}

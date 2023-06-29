package inflect

import (
	"encoding/json"
	"strings"
	"text/template"
	"unicode"

	"github.com/go-openapi/inflect"
)

var acronyms = make(map[string]struct{})
var rules = ruleset()
var Singular = rules.Singularize
var Tableize = inflect.Tableize

var TemplateFunctions = template.FuncMap{
	"lower":                 strings.ToLower,
	"upper":                 strings.ToUpper,
	"pascal":                Pascal,
	"plural":                plural,
	"singular":              rules.Singularize,
	"snake":                 Snake,
	"join":                  strings.Join,
	"camel":                 camel,
	"firstWordUpperCamel":   firstWordUpperCamel,
	"firstWordLowerCamel":   firstWordLowerCamel,
	"firstLetterLower":      firstLetterLower,
	"sub":                   sub,
	"add":                   add,
	"addStr":                addStr,
	"getTableNameOfHasMany": getTableNameOfHasMany,
	"getTableNameOfThrough": getTableNameOfThrough,

	"goType":          goType,
	"goChangeSetType": goChangeSetType,
	"goKeywordFix":    goKeywordFix,
	"goModelType":     goModelType,

	"tsType":          tsType,
	"tsChangeSetType": tsChangeSetType,
}

type Database struct {
	Name    string
	Adapter string
	Models  []*Model
}

// 为了避免循环依赖，没有直接引用schema.Model
type Model struct {
	Name       string
	TableName  string
	Timestamps bool
	BelongsTo  []*BelongsTo
	HasMany    []*HasMany
	HasOne     []*HasOne
}

type HasOne struct {
	Name       string
	ModelName  string
	Through    string
	ForeignKey string
}

type HasMany struct {
	Name       string
	ModelName  string
	Through    string
	ForeignKey string
	Source     string
}
type BelongsTo struct {
	Name        string
	ModelName   string
	ForeignKey  string
	ForeignType string
	PrimaryKey  string
	Dependent   string
	Optional    bool
	Required    bool
	Default     bool
}

func sub(a int, b int) int {
	return a - b
}

func add(a ...int) int {
	var res = 0
	for i := 0; i < len(a); i++ {
		res = res + a[i]
	}
	return res
}

func addStr(a ...string) string {
	str := ""
	for i := 0; i < len(a); i++ {
		str = str + a[i]
	}
	return str
}

func getTableNameOfThrough(through string, models interface{}) string {
	s := Singular(firstWordUpperCamel(through))
	marshaller, err := json.Marshal(models)
	if err != nil {
		return ""
	}
	var ms []*Model
	err = json.Unmarshal(marshaller, &ms)
	if err != nil {
		return ""
	}
	m := make(map[string]Model)
	for i := 0; i < len(ms); i++ {
		m[ms[i].Name] = *ms[i]
	}
	return m[s].TableName
}

func getTableNameOfHasMany(hasMany interface{}, model interface{}) string {
	marshaller, err := json.Marshal(model)
	if err != nil {
		return ""
	}
	var ms []*Model
	err = json.Unmarshal(marshaller, &ms)
	if err != nil {
		return ""
	}

	_hasMany, err := json.Marshal(hasMany)
	if err != nil {
		return ""
	}
	var h HasMany
	err = json.Unmarshal(_hasMany, &h)
	if err != nil {
		return ""
	}
	m := make(map[string]Model)
	for i := 0; i < len(ms); i++ {
		m[ms[i].Name] = *ms[i]
	}
	return m[h.ModelName].TableName
}

func camel(s string) string {
	words := strings.FieldsFunc(s, isSeparator)
	if len(words) == 1 {
		w := strings.ToLower(words[0])
		return w
	}
	return strings.ToLower(words[0]) + pascalWords(words[1:])
}

func firstLetterLower(s string) string {
	s = strings.Replace(s, " ", "", -1)
	if len(s) <= 0 {
		return ""
	}
	return strings.ToLower(s[0:1])
}

func firstWordUpperCamel(s string) string {
	words := strings.FieldsFunc(s, isSeparator)
	if len(words) == 1 {
		return pascalWords(words[0:])
	}
	return pascalWords(words[0:])
}
func firstWordLowerCamel(s string) string {
	s = firstWordUpperCamel(s)
	if s == "ID" {
		return "id"
	}
	if len(s) > 0 {
		return strings.ToLower(string(s[0])) + s[1:]
	}
	return s
}

func Snake(s string) string {
	var (
		j int
		b strings.Builder
	)
	for i := 0; i < len(s); i++ {
		r := rune(s[i])
		// Put '_' if it is not a start or end of a word, current letter is uppercase,
		// and previous is lowercase (cases like: "UserInfo"), or next letter is also
		// a lowercase and previous letter is not "_".
		if i > 0 && i < len(s)-1 && unicode.IsUpper(r) {
			if unicode.IsLower(rune(s[i-1])) ||
				j != i-1 && unicode.IsLower(rune(s[i+1])) && unicode.IsLetter(rune(s[i-1])) {
				j = i
				b.WriteString("_")
			}
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

func isSeparator(r rune) bool {
	return r == '_' || r == '-' || unicode.IsSpace(r)
}

func ruleset() *inflect.Ruleset {
	rules := inflect.NewDefaultRuleset()
	// Add common initialism from golint and more.
	for _, w := range []string{
		"ACL", "API", "ASCII", "AWS", "CPU", "CSS", "DNS", "EOF", "GB", "GUID",
		"HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "KB", "LHS", "MAC", "MB",
		"QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SQL", "SSH", "SSO", "TCP",
		"TLS", "TTL", "UDP", "UI", "UID", "URI", "URL", "UTF8", "UUID", "VM",
		"XML", "XMPP", "XSRF", "XSS",
	} {
		acronyms[w] = struct{}{}
		rules.AddAcronym(w)
	}
	return rules
}

func pascalWords(words []string) string {
	for i, w := range words {
		upper := strings.ToUpper(w)
		if _, ok := acronyms[upper]; ok {
			words[i] = upper
		} else {
			words[i] = rules.Capitalize(w)
		}
	}
	return strings.Join(words, "")
}

func Pascal(s string) string {
	words := strings.FieldsFunc(s, isSeparator)
	return pascalWords(words)
}

func plural(name string) string {
	p := rules.Pluralize(name)
	if p == name {
		p += "Slice"
	}
	return p
}

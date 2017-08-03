package generator

import (
	"regexp"
	"time"
)

type Replacer func(string) (string, bool)

type Generator struct {
	replacers []Replacer
}

func NewGenerator() *Generator {
	return &Generator{
		replacers: []Replacer{},
	}
}

func (g *Generator) AddReplacer(replacer Replacer) {
	g.replacers = append(g.replacers, replacer)
}

func defaultReplacer(m string) (string, bool) {
	switch m {
	case "%F":
		return time.Now().Format("2006/01/02"), true
	case "%T":
		return time.Now().Format("15:04:05"), true
	case "%%":
		return "%", true
	}
	return "", false
}

var placeholderRegexp = regexp.MustCompile(`%\w`)

func (g *Generator) Execute(tmpl string) string {
	return placeholderRegexp.ReplaceAllStringFunc(tmpl, func(m string) string {
		for _, replacer := range append(g.replacers, defaultReplacer) {
			r, ok := replacer(m)
			if ok {
				return r
			}
		}
		return m
	})
}

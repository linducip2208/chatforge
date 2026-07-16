package msgtemplate

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
)

var spintaxRe = regexp.MustCompile(`\{([^{}]*\|[^{}]*)\}`)

// Vars holds the values used to replace {placeholders}.
type Vars struct {
	Name    string
	Phone   string
	Message string
}

// Render applies variable replacement then spintax to a template string.
// This is an ADDITIVE layer — it does not change auto-reply matching.
func Render(tpl string, v Vars) string {
	out := replaceVars(tpl, v)
	out = spintax(out)
	return out
}

func replaceVars(s string, v Vars) string {
	now := time.Now()
	repl := strings.NewReplacer(
		"{name}", v.Name,
		"{phone}", strings.TrimPrefix(v.Phone, "+"),
		"{message}", v.Message,
		"{time}", now.Format("15:04"),
		"{date}", now.Format("02 Jan 2006"),
		"{datetime}", now.Format("02 Jan 2006 15:04"),
	)
	return repl.Replace(s)
}

// spintax: {a|b|c} -> random choice, innermost-first.
func spintax(s string) string {
	guard := 0
	for spintaxRe.MatchString(s) && guard < 100 {
		s = spintaxRe.ReplaceAllStringFunc(s, func(m string) string {
			inner := m[1 : len(m)-1]
			choices := strings.Split(inner, "|")
			return strings.TrimSpace(choices[rand.Intn(len(choices))])
		})
		guard++
	}
	return s
}

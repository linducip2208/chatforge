package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// Lang is a single language pack.
type Lang struct {
	Code    string            // file code, e.g. "id", "en"
	ISO     string            // flag-icon iso, e.g. "id", "us"
	Name    string            // display name
	Flag    string            // flag-icon code
	strings map[string]string // key -> text
}

var (
	langs   = map[string]*Lang{}
	order   []string
	def     = "id"
	loadMu  sync.RWMutex
)

// Load reads all *.json files from dir into memory.
func Load(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	loadMu.Lock()
	defer loadMu.Unlock()
	langs = map[string]*Lang{}
	order = nil
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		code := e.Name()[:len(e.Name())-len(".json")]
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return err
		}
		var m map[string]string
		if err := json.Unmarshal(raw, &m); err != nil {
			return fmt.Errorf("lang %s: %w", code, err)
		}
		l := &Lang{Code: code, ISO: m["iso"], Name: m["name"], Flag: m["flag"], strings: m}
		if l.ISO == "" {
			l.ISO = code
		}
		if l.Name == "" {
			l.Name = code
		}
		if l.Flag == "" {
			l.Flag = l.ISO
		}
		langs[code] = l
		order = append(order, code)
	}
	sort.Strings(order)
	if _, ok := langs[def]; !ok && len(order) > 0 {
		def = order[0]
	}
	if len(langs) == 0 {
		return fmt.Errorf("no language files found in %s", dir)
	}
	return nil
}

// Default returns the default language code.
func Default() string { return def }

// Has reports whether a language code exists.
func Has(code string) bool {
	loadMu.RLock()
	defer loadMu.RUnlock()
	_, ok := langs[code]
	return ok
}

// T translates a key for a language (falls back to key).
func T(code, key string) string {
	loadMu.RLock()
	defer loadMu.RUnlock()
	l := langs[code]
	if l == nil {
		l = langs[def]
	}
	if l != nil {
		if v, ok := l.strings[key]; ok {
			return v
		}
	}
	return key
}

// Translator returns a bound T function for use in templates.
func Translator(code string) func(string) string {
	if !Has(code) {
		code = def
	}
	return func(key string) string { return T(code, key) }
}

// List returns all languages (for the switcher), in order.
func List() []*Lang {
	loadMu.RLock()
	defer loadMu.RUnlock()
	out := make([]*Lang, 0, len(order))
	for _, c := range order {
		out = append(out, langs[c])
	}
	return out
}

// Get returns a language pack.
func Get(code string) *Lang {
	loadMu.RLock()
	defer loadMu.RUnlock()
	if l, ok := langs[code]; ok {
		return l
	}
	return langs[def]
}

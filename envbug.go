package dbgo

// This file provides cached, concurrency-safe lookups of boolean values
// derived from environment variables (or values set programmatically with
// SetFlag). Each variable's raw value is read at most once and its parsed
// boolean is cached for the life of the process.

import (
	"os"
	"sync"
)

// envCache caches the raw string value for each looked-up name.
var envCache map[string]string

// envCache0 caches the parsed boolean for each looked-up name.
var envCache0 map[string]bool

// envLock guards envCache and envCache0.
var envLock sync.RWMutex

func init() {
	envCache = make(map[string]string)
	envCache0 = make(map[string]bool)
}

// SetFlag sets a named flag's raw string value and immediately parses and
// caches it, overriding any environment variable of the same name for
// subsequent ChkEnv calls. It is safe for concurrent use.
func SetFlag(name, sValue string) {
	envLock.Lock()
	defer envLock.Unlock()
	envCache[name] = sValue
	envCache0[name] = ParseBool(sValue)
}

// ChkEnv reports whether the named environment variable parses as a true value
// (see ParseBool). The os.Getenv lookup is performed at most once per name; the
// raw value and parsed result are cached. A value set via SetFlag takes
// precedence over the actual environment.
func ChkEnv(envVar string) bool {
	envLock.RLock()
	// Fast path: a previously parsed answer exists.
	if vv, ok := envCache0[envVar]; ok {
		envLock.RUnlock()
		return vv
	}
	// Use a cached raw value (from SetFlag or a prior lookup) if present...
	v, hasRaw := envCache[envVar]
	envLock.RUnlock()

	if !hasRaw {
		// ...otherwise fall back to the actual environment.
		v = os.Getenv(envVar)
	}
	b := ParseBool(v)

	envLock.Lock()
	envCache[envVar] = v
	envCache0[envVar] = b
	envLock.Unlock()
	return b
}

// trueValues is the set of strings that ParseBool treats as true.
var trueValues map[string]bool

func init() {
	trueValues = make(map[string]bool)
	trueValues["t"] = true
	trueValues["T"] = true
	trueValues["yes"] = true
	trueValues["Yes"] = true
	trueValues["YES"] = true
	trueValues["1"] = true
	trueValues["true"] = true
	trueValues["True"] = true
	trueValues["TRUE"] = true
	trueValues["on"] = true
	trueValues["On"] = true
	trueValues["ON"] = true
}

// ParseBool reports whether s is one of the recognized "true" spellings:
// t, T, yes, Yes, YES, 1, true, True, TRUE, on, On, ON. Any other value
// (including the empty string) is false.
func ParseBool(s string) bool {
	_, b := trueValues[s]
	return b
}

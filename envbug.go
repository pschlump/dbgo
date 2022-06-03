package dbgo

import (
	"os"
	"sync"
)

var envCache map[string]string
var envCache0 map[string]bool
var envlock sync.RWMutex

func init() {
	envCache = make(map[string]string)
	envCache0 = make(map[string]bool)
}

func SetFlag(name, sValue string) {
	envlock.Lock()
	envCache[name] = sValue
	envCache0[name] = ParseBool(sValue)
	envlock.Unlock()
}

func ChkEnv(envVar string) bool {
	envlock.RLock()
	vv, ok := envCache0[envVar]
	if ok {
		envlock.RUnlock()
		return vv
	}
	v, ok := envCache[envVar]
	if ok {
		envlock.RUnlock()
		return ParseBool(v)
	}
	envlock.RUnlock()

	v = os.Getenv(envVar)
	envlock.Lock()
	envCache[envVar] = v
	envlock.Unlock()
	x := ChkEnv(envVar)
	envlock.Lock()
	envCache0[envVar] = x
	envlock.Unlock()
	return x
}

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

func ParseBool(s string) (b bool) {
	_, b = trueValues[s]
	return
	//if InArray(s, []string{"t", "T", "yes", "Yes", "YES", "1", "true", "True", "TRUE", "on", "On", "ON"}) {
	//	return true
	//}
	//return false
}

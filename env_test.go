package dbgo

import "testing"

func TestParseBool(t *testing.T) {
	truthy := []string{"t", "T", "yes", "Yes", "YES", "1", "true", "True", "TRUE", "on", "On", "ON"}
	for _, s := range truthy {
		if !ParseBool(s) {
			t.Errorf("ParseBool(%q) = false, want true", s)
		}
	}
	falsy := []string{"", "no", "0", "false", "off", "2", "random", "y", "yeah", "enable"}
	for _, s := range falsy {
		if ParseBool(s) {
			t.Errorf("ParseBool(%q) = true, want false", s)
		}
	}
}

func TestSetFlagAndChkEnv(t *testing.T) {
	// Use a unique name to avoid colliding with other tests or the examples.
	const name = "DBGO_TEST_FLAG_UNIQUE_42"

	SetFlag(name, "yes")
	if !ChkEnv(name) {
		t.Errorf("ChkEnv(%q) after SetFlag yes = false, want true", name)
	}

	SetFlag(name, "nope")
	if ChkEnv(name) {
		t.Errorf("ChkEnv(%q) after SetFlag nope = true, want false", name)
	}
}

func TestChkEnvCachesEnvVar(t *testing.T) {
	// A unique, almost-certainly-unset variable: ChkEnv returns false and
	// caches that result so subsequent os.Getenv calls are avoided.
	const name = "DBGO_TEST_ENV_CACHE_UNIQUE_7"

	if ChkEnv(name) {
		t.Errorf("ChkEnv(unset %q) = true, want false", name)
	}
	// Second and subsequent calls hit the cache and stay false.
	for i := 0; i < 100; i++ {
		if ChkEnv(name) {
			t.Errorf("ChkEnv(cached unset) = true on iter %d, want false", i)
		}
	}
}

func TestChkEnvReadsRealEnv(t *testing.T) {
	// A never-before-seen name exercises the os.Getenv code path (the cache
	// has no prior entry). t.Setenv sets it for this process only and
	// restores it on cleanup.
	const name = "DBGO_TEST_REAL_ENV_UNIQUE_99"

	t.Setenv(name, "YES")
	if !ChkEnv(name) {
		t.Errorf("ChkEnv(%q) with env=YES = false, want true", name)
	}

	// Because the value is now cached, changing the env var does not change
	// the answer for the rest of this process (documented caching behavior).
	t.Setenv(name, "no")
	if !ChkEnv(name) {
		t.Errorf("ChkEnv(%q) changed to false after cache, want cached true", name)
	}
}

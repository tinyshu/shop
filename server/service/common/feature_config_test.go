package common

import "testing"

func TestParseFeatureBool(t *testing.T) {
	cases := []struct {
		in      string
		enabled bool
		ok      bool
	}{
		{"1", true, true},
		{"true", true, true},
		{"TRUE", true, true},
		{"on", true, true},
		{"0", false, true},
		{"false", false, true},
		{"off", false, true},
		{" 1 ", true, true},
		{"courier", false, false},
		{"", false, false},
		{"yes", false, false},
	}
	for _, c := range cases {
		en, ok := parseFeatureBool(c.in)
		if en != c.enabled || ok != c.ok {
			t.Fatalf("parseFeatureBool(%q)=(%v,%v) want (%v,%v)", c.in, en, ok, c.enabled, c.ok)
		}
	}
}

func TestFeatureDefaultsWithoutDB(t *testing.T) {
	// global.DB 未初始化时必须返回默认值，不 panic
	if got := FeatureEnabled(KeyUserAudit, false); got != false {
		t.Fatalf("FeatureEnabled without DB: got %v want false", got)
	}
	if got := FeatureEnabled(KeyUserAudit, true); got != true {
		t.Fatalf("FeatureEnabled without DB: got %v want true", got)
	}
	if got := FeatureString(KeyCourierMode, "courier"); got != "courier" {
		t.Fatalf("FeatureString without DB: got %q want courier", got)
	}
}

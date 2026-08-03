package main

import "testing"

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"0.1.0", "0.1.0", 0},
		{"0.1.0", "0.2.0", -1},
		{"0.2.0", "0.1.0", 1},
		{"1.0.0", "0.9.9", 1},
		{"1.2", "1.2.0", 0},
		{"1.2.0", "1.2", 0},
		{"1.2.1", "1.2", 1},
		{"0.0.1", "0.0.2", -1},
	}
	for _, c := range cases {
		got, err := compareVersions(c.a, c.b)
		if err != nil {
			t.Errorf("compareVersions(%q, %q) unexpected error: %v", c.a, c.b, err)
			continue
		}
		if got != c.want {
			t.Errorf("compareVersions(%q, %q) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestCompareVersionsInvalid(t *testing.T) {
	if _, err := compareVersions("dev", "0.1.0"); err == nil {
		t.Error("compareVersions(\"dev\", \"0.1.0\") expected an error, got nil")
	}
	if _, err := compareVersions("0.1.0", "not-a-version"); err == nil {
		t.Error("compareVersions(\"0.1.0\", \"not-a-version\") expected an error, got nil")
	}
}

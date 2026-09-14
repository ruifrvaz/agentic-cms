package main

import "testing"

func TestParseInitArgs(t *testing.T) {
	cases := []struct {
		name       string
		args       []string
		wantDir    string
		wantType   string
		wantErrNil bool
	}{
		{"no args", nil, ".", "", true},
		{"dir only", []string{"myproject"}, "myproject", "", true},
		{"type only, dir defaults", []string{"--type", "candidate-interview"}, ".", "candidate-interview", true},
		{"dir then type", []string{"myproject", "--type", "candidate-interview"}, "myproject", "candidate-interview", true},
		{"type then dir", []string{"--type", "candidate-interview", "myproject"}, "myproject", "candidate-interview", true},
		{"type list", []string{"--type", "list"}, ".", "list", true},
		{"dangling --type", []string{"--type"}, "", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir, typeName, err := parseInitArgs(c.args)
			if (err == nil) != c.wantErrNil {
				t.Fatalf("parseInitArgs(%v) error = %v, want nil-ness %v", c.args, err, c.wantErrNil)
			}
			if err != nil {
				return
			}
			if dir != c.wantDir || typeName != c.wantType {
				t.Errorf("parseInitArgs(%v) = (%q, %q), want (%q, %q)", c.args, dir, typeName, c.wantDir, c.wantType)
			}
		})
	}
}

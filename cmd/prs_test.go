package cmd

import (
	"testing"
)

func TestIsDependabotAuthor(t *testing.T) {
	cases := []struct {
		login string
		want  bool
	}{
		{"app/dependabot", true},
		{"dependabot[bot]", true},
		{"octocat", false},
		{"dependabot", false},
		{"", false},
	}

	for _, tc := range cases {
		t.Run(tc.login, func(t *testing.T) {
			if got := isDependabotAuthor(tc.login); got != tc.want {
				t.Errorf("isDependabotAuthor(%q) = %v, want %v", tc.login, got, tc.want)
			}
		})
	}
}

func TestAllChecksPassedOrSkipped(t *testing.T) {
	cases := []struct {
		name   string
		checks []statusCheck
		want   bool
	}{
		{
			name:   "no checks",
			checks: nil,
			want:   true,
		},
		{
			name:   "all success",
			checks: []statusCheck{{Conclusion: "SUCCESS"}, {Conclusion: "SUCCESS"}},
			want:   true,
		},
		{
			name:   "all skipped",
			checks: []statusCheck{{Conclusion: "SKIPPED"}, {Conclusion: "SKIPPED"}},
			want:   true,
		},
		{
			name:   "mixed success and skipped",
			checks: []statusCheck{{Conclusion: "SUCCESS"}, {Conclusion: "SKIPPED"}},
			want:   true,
		},
		{
			name:   "one failure",
			checks: []statusCheck{{Conclusion: "SUCCESS"}, {Conclusion: "FAILURE"}},
			want:   false,
		},
		{
			name:   "pending",
			checks: []statusCheck{{Conclusion: ""}},
			want:   false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := allChecksPassedOrSkipped(tc.checks); got != tc.want {
				t.Errorf("allChecksPassedOrSkipped(%v) = %v, want %v", tc.checks, got, tc.want)
			}
		})
	}
}

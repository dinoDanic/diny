package commit

import (
	"strings"
	"testing"

	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/groq"
)

func TestNormalizePlanPath(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"features/auth/login.tsx", "features/auth/login.tsx"},
		{"a/features/auth/login.tsx", "features/auth/login.tsx"},
		{"b/features/auth/login.tsx", "features/auth/login.tsx"},
		{"i/features/auth/login.tsx", "features/auth/login.tsx"},
		{"c/features/auth/login.tsx", "features/auth/login.tsx"},
		{"w/features/auth/login.tsx", "features/auth/login.tsx"},
		{"./features/auth/login.tsx", "features/auth/login.tsx"},
		{"  a/features/auth/login.tsx  ", "features/auth/login.tsx"},
		{"a/b/features/auth/login.tsx", "features/auth/login.tsx"},
		{"./a/features/auth/login.tsx", "features/auth/login.tsx"},
		{"z/features/auth/login.tsx", "z/features/auth/login.tsx"}, // unknown single-letter prefix is left alone
		{"ab/features/auth/login.tsx", "ab/features/auth/login.tsx"},
	}
	for _, c := range cases {
		got := normalizePlanPath(c.in)
		if got != c.want {
			t.Errorf("normalizePlanPath(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestValidatePlan_StripsBogusPrefix(t *testing.T) {
	staged := []git.StagedFile{
		{Status: "A", Path: "features/auth/login.tsx"},
		{Status: "A", Path: "features/auth/register.tsx"},
	}
	plan := []groq.SplitGroup{
		{
			Order:   1,
			Type:    "feat",
			Message: "add auth",
			Files:   []string{"i/features/auth/login.tsx", "a/features/auth/register.tsx"},
		},
	}
	if err := ValidatePlan(plan, staged); err != nil {
		t.Fatalf("ValidatePlan returned error: %v", err)
	}
	if plan[0].Files[0] != "features/auth/login.tsx" {
		t.Errorf("Files[0] not rewritten: %q", plan[0].Files[0])
	}
	if plan[0].Files[1] != "features/auth/register.tsx" {
		t.Errorf("Files[1] not rewritten: %q", plan[0].Files[1])
	}
}

func TestValidatePlan_UnresolvablePathStillErrors(t *testing.T) {
	staged := []git.StagedFile{
		{Status: "A", Path: "features/auth/login.tsx"},
	}
	plan := []groq.SplitGroup{
		{
			Order:   1,
			Type:    "feat",
			Message: "add auth",
			Files:   []string{"features/auth/nonexistent.tsx"},
		},
	}
	err := ValidatePlan(plan, staged)
	if err == nil {
		t.Fatal("expected error for unresolvable path")
	}
	if !strings.Contains(err.Error(), "not staged") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidatePlan_DuplicateAcrossGroupsAfterNormalize(t *testing.T) {
	staged := []git.StagedFile{
		{Status: "A", Path: "features/auth/login.tsx"},
	}
	plan := []groq.SplitGroup{
		{Order: 1, Type: "feat", Message: "a", Files: []string{"features/auth/login.tsx"}},
		{Order: 2, Type: "feat", Message: "b", Files: []string{"a/features/auth/login.tsx"}},
	}
	err := ValidatePlan(plan, staged)
	if err == nil {
		t.Fatal("expected duplicate error")
	}
	if !strings.Contains(err.Error(), "groups") {
		t.Errorf("unexpected error message: %v", err)
	}
}

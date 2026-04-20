package commit

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/groq"
)

// SplitGroup is an alias to the transport type for ergonomic use across the TUI.
type SplitGroup = groq.SplitGroup

// SplitRequestExtras is an alias for transport-layer extras.
type SplitRequestExtras = groq.RequestExtras

// CreateSplitPlan asks the backend to group the staged diff into multiple commits.
func CreateSplitPlan(gitDiff string, cfg *config.Config, extras *SplitRequestExtras) ([]SplitGroup, error) {
	return groq.CreateSplitPlanWithGroq(gitDiff, cfg, extras)
}

// normalizePlanPath strips bogus prefixes that LLMs sometimes add to file
// paths in split plans (e.g. "a/", "b/", "i/", "c/" copied from git diff
// headers, or leading "./"). The returned path can be compared against the
// real staged-file set.
func normalizePlanPath(p string) string {
	p = strings.TrimSpace(p)
	for {
		switch {
		case strings.HasPrefix(p, "./"):
			p = p[2:]
		case len(p) >= 2 && p[1] == '/' && strings.ContainsRune("abciw", rune(p[0])):
			p = p[2:]
		default:
			return p
		}
	}
}

// ValidatePlan checks that every staged file appears in exactly one group and
// no group references files that are not currently staged. File paths returned
// by the model are normalized first (see normalizePlanPath) and the plan is
// mutated in place to use the canonical staged path, so downstream `git add`
// calls receive valid paths. Returns the (possibly rewritten) plan.
func ValidatePlan(plan []SplitGroup, staged []git.StagedFile) error {
	if len(plan) == 0 {
		return fmt.Errorf("plan has no groups")
	}

	stagedSet := make(map[string]struct{}, len(staged))
	for _, f := range staged {
		stagedSet[f.Path] = struct{}{}
	}

	seen := make(map[string]int, len(staged))
	for gi := range plan {
		g := &plan[gi]
		for fi, f := range g.Files {
			canonical := f
			if _, ok := stagedSet[canonical]; !ok {
				canonical = normalizePlanPath(f)
			}
			if _, ok := stagedSet[canonical]; !ok {
				return fmt.Errorf("plan references %q which is not staged", f)
			}
			if canonical != f {
				g.Files[fi] = canonical
			}
			if groupIdx, dup := seen[canonical]; dup {
				return fmt.Errorf("plan assigns %q to groups %d and %d", canonical, groupIdx+1, g.Order)
			}
			seen[canonical] = g.Order - 1
		}
	}

	var missing []string
	for _, f := range staged {
		if _, ok := seen[f.Path]; !ok {
			missing = append(missing, f.Path)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("plan does not assign %d staged file(s): %s", len(missing), strings.Join(missing, ", "))
	}

	return nil
}

// NormalizePlan sorts groups by Order and renumbers them 1..N.
func NormalizePlan(plan []SplitGroup) []SplitGroup {
	out := make([]SplitGroup, len(plan))
	copy(out, plan)
	sort.SliceStable(out, func(i, j int) bool { return out[i].Order < out[j].Order })
	for i := range out {
		out[i].Order = i + 1
	}
	return out
}

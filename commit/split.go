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

// ValidatePlan checks that every staged file appears in exactly one group and
// no group references files that are not currently staged.
func ValidatePlan(plan []SplitGroup, staged []git.StagedFile) error {
	if len(plan) == 0 {
		return fmt.Errorf("plan has no groups")
	}

	stagedSet := make(map[string]struct{}, len(staged))
	for _, f := range staged {
		stagedSet[f.Path] = struct{}{}
	}

	seen := make(map[string]int, len(staged))
	for _, g := range plan {
		for _, f := range g.Files {
			if _, ok := stagedSet[f]; !ok {
				return fmt.Errorf("plan references %q which is not staged", f)
			}
			if groupIdx, dup := seen[f]; dup {
				return fmt.Errorf("plan assigns %q to groups %d and %d", f, groupIdx+1, g.Order)
			}
			seen[f] = g.Order - 1
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

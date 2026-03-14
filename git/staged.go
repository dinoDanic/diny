package git

import (
	"os/exec"
	"strings"
)

type StagedFile struct {
	Status string // "A", "M", "D", "R"
	Path   string
}

func GetStagedFiles() ([]StagedFile, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-status")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	raw := strings.TrimSpace(string(output))
	if raw == "" {
		return nil, nil
	}

	var files []StagedFile
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) < 2 {
			continue
		}
		status := parts[0]
		path := parts[1]

		// Renames show as R100\told\tnew — normalize to "R" with new path
		if strings.HasPrefix(status, "R") {
			status = "R"
			renameParts := strings.SplitN(path, "\t", 2)
			if len(renameParts) == 2 {
				path = renameParts[1]
			}
		}

		files = append(files, StagedFile{Status: status, Path: path})
	}

	return files, nil
}

// GetUnstagedFiles returns modified, deleted, and untracked files not yet staged.
func GetUnstagedFiles() ([]StagedFile, error) {
	var files []StagedFile
	seen := map[string]bool{}

	// 1. Modified/deleted unstaged changes (working tree vs index)
	diffOut, err := exec.Command("git", "diff", "--name-status").Output()
	if err == nil {
		for _, line := range strings.Split(strings.TrimSpace(string(diffOut)), "\n") {
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, "\t", 2)
			if len(parts) < 2 {
				continue
			}
			status, path := parts[0], parts[1]
			if seen[path] {
				continue
			}
			seen[path] = true
			switch status {
			case "M":
				files = append(files, StagedFile{Status: "M", Path: path})
			case "D":
				files = append(files, StagedFile{Status: "D", Path: path})
			default:
				files = append(files, StagedFile{Status: status, Path: path})
			}
		}
	}

	// 2. Untracked (new) files
	lsOut, err := exec.Command("git", "ls-files", "--others", "--exclude-standard").Output()
	if err == nil {
		for _, path := range strings.Split(strings.TrimSpace(string(lsOut)), "\n") {
			if path == "" || seen[path] {
				continue
			}
			seen[path] = true
			files = append(files, StagedFile{Status: "?", Path: path})
		}
	}

	return files, nil
}

func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

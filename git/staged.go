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
	cmd := exec.Command("git", "status", "--porcelain")
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
		if len(line) < 4 {
			continue
		}
		xy := line[:2]
		path := line[3:]

		if xy == "??" {
			files = append(files, StagedFile{Status: "?", Path: path})
			continue
		}
		unstaged := rune(xy[1])
		if unstaged == ' ' || unstaged == '!' {
			continue
		}
		switch unstaged {
		case 'M':
			files = append(files, StagedFile{Status: "M", Path: path})
		case 'D':
			files = append(files, StagedFile{Status: "D", Path: path})
		default:
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

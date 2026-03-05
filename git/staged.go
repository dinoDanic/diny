package git

import (
	"os/exec"
	"strings"
)

type StagedFile struct {
	Status string
	Path   string
}

func GetStagedFiles() ([]StagedFile, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-status")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var files []StagedFile
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			continue
		}
		files = append(files, StagedFile{
			Status: parts[0],
			Path:   parts[1],
		})
	}

	return files, nil
}

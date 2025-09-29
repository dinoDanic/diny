/*
Copyright © 2025 NAME HERE dino.danic@gmail.com
*/
package git

import (
	"os/exec"
	"strings"
)

func GetGitName() (string, error) {
	cmd := exec.Command("git", "config", "user.name")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

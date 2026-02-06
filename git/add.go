package git

import "os/exec"

func AddAll() error {
	cmd := exec.Command("git", "add", "-A")
	return cmd.Run()
}

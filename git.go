package u

import (
	"fmt"
	"os/exec"
	"strings"
)

func GitPullMust(dir string) {
	cmd := exec.Command("git", "pull")
	if dir != "" {
		cmd.Dir = dir
	}
	RunCmdMust(cmd)
}

func GitStatusMust(dir string) string {
	cmd := exec.Command("git", "status")
	if dir != "" {
		cmd.Dir = dir
	}
	return RunCmdMust(cmd)
}

func IsGitClean(dir string) bool {
	s := GitStatusMust(dir)
	expected := []string{
		"On branch master",
		"Your branch is up to date with 'origin/master'.",
		"nothing to commit, working tree clean",
	}
	for _, exp := range expected {
		if !strings.Contains(s, exp) {
			Logf("Git repo in '%s' not clean.\nDidn't find '%s' in output of git status:\n%s\n", dir, exp, s)
			return false
		}
	}
	return true
}

func EnsureGitClean(dir string) {
	if !IsGitClean(dir) {
		Must(fmt.Errorf("git repo in '%s' is not clean", dir))
	}
}

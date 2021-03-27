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
	expected1 := []string{
		"On branch master",
		"Your branch is up to date with 'origin/master'.",
		"nothing to commit, working tree clean",
	}
	expected2 := []string{
		"On branch main",
		"Your branch is up to date with 'origin/main'.",
		"nothing to commit, working tree clean",
	}
	{
		hasAll := true
		for _, exp := range expected1 {
			if !strings.Contains(s, exp) {
				//Logf("Git repo in '%s' not clean.\nDidn't find '%s' in output of git status:\n%s\n", dir, exp, s)
				hasAll = false
			}
		}
		if hasAll {
			return true
		}
	}
	{
		hasAll := true
		for _, exp := range expected2 {
			if !strings.Contains(s, exp) {
				hasAll = false
			}
		}
		if hasAll {
			return true
		}
	}
	Logf("Git repo in '%s' not clean.\nGit status:\n%s\n", dir, s)
	return false
}

func EnsureGitClean(dir string) {
	if !IsGitClean(dir) {
		Must(fmt.Errorf("git repo in '%s' is not clean", dir))
	}
}

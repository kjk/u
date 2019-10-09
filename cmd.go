package u

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func RunCmdMust(cmd *exec.Cmd) string {
	fmt.Printf("> %s\n", cmd)
	canCapture := (cmd.Stdout == nil) && (cmd.Stderr == nil)
	if canCapture {
		out, err := cmd.CombinedOutput()
		if err == nil {
			if len(out) > 0 {
				fmt.Printf("Output:\n%s\n", string(out))
			}
			return string(out)
		}
		fmt.Printf("cmd '%s' failed with '%s'. Output:\n%s\n", cmd, err, string(out))
		Must(err)
		return string(out)
	}
	err := cmd.Run()
	if err == nil {
		return ""
	}
	fmt.Printf("cmd '%s' failed with '%s'\n", cmd, err)
	Must(err)
	return ""
}

func OpenNotepadWithFileMust(path string) {
	cmd := exec.Command("notepad.exe", path)
	err := cmd.Start()
	Must(err)
}

func OpenCodeDiffMust(path1, path2 string) {
	if runtime.GOOS == "darwin" {
		path1 = strings.Replace(path1, ".\\", "./", -1)
		path2 = strings.Replace(path2, ".\\", "./", -1)
	}
	cmd := exec.Command("code", "--new-window", "--diff", path1, path2)
	fmt.Printf("> %s\n", cmd)
	err := cmd.Start()
	Must(err)
}

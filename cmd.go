package u

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// FmtCmdShort formats exec.Cmd in a short way
func FmtCmdShort(cmd exec.Cmd) string {
	cmd.Path = filepath.Base(cmd.Path)
	return cmd.String()
}

// RunCmdLoggedMust runs a command and returns its stdout
// Shows output as it happens
func RunCmdLoggedMust(cmd *exec.Cmd) string {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return RunCmdMust(cmd)
}

// RunCmdMust runs a command and returns its stdout
func RunCmdMust(cmd *exec.Cmd) string {
	fmt.Printf("> %s\n", FmtCmdShort(*cmd))
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
	fmt.Printf("> %s\n", FmtCmdShort(*cmd))
	err := cmd.Start()
	Must(err)
}

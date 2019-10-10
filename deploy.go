package u

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

var (
	ServerIPAddress  string
	IdentityFilePath string
)

func SshInteractive(identityFile string, user string) {
	cmd := exec.Command("ssh", "-i", identityFile, user)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	RunCmdMust(cmd)
}

func LoginAsRoot() {
	user := fmt.Sprintf("root@%s", ServerIPAddress)
	SshInteractive(IdentityFilePath, user)
}

// "-o StrictHostKeyChecking=no" is for the benefit of CI which start
// fresh environment
func ScpCopy(localSrcPath string, serverDstPath string) {
	cmd := exec.Command("scp", "-o", "StrictHostKeyChecking=no", "-i", IdentityFilePath, localSrcPath, serverDstPath)
	RunCmdMust(cmd)
}

// "-o StrictHostKeyChecking=no" is for the benefit of CI which start
// fresh environment
func SshExec(user string, script string) {
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "-i", IdentityFilePath, user)
	r := bytes.NewBufferString(script)
	cmd.Stdin = r
	RunCmdMust(cmd)
}

func MakeExecScript(name string) string {
	script := fmt.Sprintf(`
chmod ug+x ./%s
./%s
rm ./%s
	`, name, name, name)
	return script
}

func CopyAndExecServerScript(scriptName, user string) {
	serverAndUser := fmt.Sprintf("%s@%s", user, ServerIPAddress)
	serverPath := "/root/" + scriptName
	if user != "root" {
		serverPath = "/home/" + user + "/" + scriptName
	}
	{
		serverDstPath := fmt.Sprintf("%s:%s", serverAndUser, serverPath)
		ScpCopy(scriptName, serverDstPath)
	}
	{
		script := MakeExecScript(scriptName)
		SshExec(serverAndUser, script)
	}
}

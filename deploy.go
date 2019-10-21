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

func SshInteractive(user string) {
	PanicIf(IdentityFilePath == "", "No identity file")
	cmd := exec.Command("ssh", "-i", IdentityFilePath, user)
	cmd.Stdin = os.Stdin
	RunCmdLoggedMust(cmd)
}

func LoginAsRoot() {
	user := fmt.Sprintf("root@%s", ServerIPAddress)
	SshInteractive(user)
}

// "-o StrictHostKeyChecking=no" is for the benefit of CI which start
// fresh environment
func ScpCopy(localSrcPath string, serverDstPath string) {
	PanicIf(IdentityFilePath == "", "No identity file")
	cmd := exec.Command("scp", "-o", "StrictHostKeyChecking=no", "-i", IdentityFilePath, localSrcPath, serverDstPath)
	RunCmdLoggedMust(cmd)
}

// "-o StrictHostKeyChecking=no" is for the benefit of CI which start
// fresh environment
func SshExec(user string, script string) {
	PanicIf(IdentityFilePath == "", "No identity file")
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "-i", IdentityFilePath, user)
	r := bytes.NewBufferString(script)
	cmd.Stdin = r
	RunCmdLoggedMust(cmd)
}

func MakeExecScript(name string) string {
	script := fmt.Sprintf(`
chmod ug+x ./%s
./%s
rm ./%s
	`, name, name, name)
	return script
}

func CopyAndExecServerScript(scriptPath, user string) {
	PanicIf(IdentityFilePath == "", "No identity file")
	PanicIf(!FileExists(scriptPath), "script file '%s' doesn't exist", scriptPath)
	serverAndUser := fmt.Sprintf("%s@%s", user, ServerIPAddress)
	serverPath := "/root/" + scriptPath
	if user != "root" {
		serverPath = "/home/" + user + "/" + scriptPath
	}
	{
		serverDstPath := fmt.Sprintf("%s:%s", serverAndUser, serverPath)
		ScpCopy(scriptPath, serverDstPath)
	}
	{
		script := MakeExecScript(scriptPath)
		SshExec(serverAndUser, script)
	}
}

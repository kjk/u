package u

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	ServerIPAddress  string
	IdentityFilePath string
)

func SshInteractive(user string) {
	panicIfServerInfoNotSet()
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
	panicIfServerInfoNotSet()
	cmd := exec.Command("scp", "-o", "StrictHostKeyChecking=no", "-i", IdentityFilePath, localSrcPath, serverDstPath)
	RunCmdLoggedMust(cmd)
}

// "-o StrictHostKeyChecking=no" is for the benefit of CI which start
// fresh environment
func SshExec(user string, script string) {
	panicIfServerInfoNotSet()
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

func panicIfServerInfoNotSet() {
	PanicIf(IdentityFilePath == "", "IdentityFilePath not set")
	PanicIf(!FileExists(IdentityFilePath), "IdentityFilePath '%s' doesn't exist", IdentityFilePath)
	PanicIf(ServerIPAddress == "", "ServerIPAddress not set")
}

// CopyAndExecServerScript copies a given script to the server and executes
// it under a given user name
func CopyAndExecServerScript(scriptPath, user string) {
	panicIfServerInfoNotSet()
	PanicIf(!FileExists(scriptPath), "script file '%s' doesn't exist", scriptPath)

	serverAndUser := fmt.Sprintf("%s@%s", user, ServerIPAddress)
	serverPath := "/root/" + filepath.Base(scriptPath)
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

// This service is used to backup files to cloud object storage.
//
// Features:
// 	- This service will begin backing up daily at 2:00AM local time.
// 	- Support cloud object storage, such as AWS S3 and OpenStack.
// 	- Multiple backup destinations can be configured.
// 	- Backup files is gz or tgz format.
// 	- Backup files is encrypted.
package main

import (
	_ "embed"
	"github.com/tappoy/env"

	"os"
	"strings"
	"net"
	)

// Usage is the help message.
//go:embed Usage.txt
var Usage string

// VaultPrefix is the prefix of the vault keys.
var VaultPrefix string

// PidFile is the path of the pid file.
const PidFile = "/var/run/backupd.pid"

// SocketFile is the path of the socket file.
const SocketFile = "/var/run/backupd.sock"

// Log directory
const LogDir = "/var/log/backupd"

// RunHelpMessage prints the simple error message.
func RunHelpMessage() {
	env.Errf("Run %s help\n", env.Args[0])
}

func help() {
	env.Outf(Usage)
}

func info() {
	_, err := os.Stat(PidFile)
	if err != nil {
		env.Outf("Service is not running.\n")
		return
	}

	pid, err := os.ReadFile(PidFile)
	var pidStr string
	if err != nil {
		env.Outf("Failed to read pid file.\n")
		env.Outf("Cannot access the service.\n")
		return
	}

	pidStr = string(pid)
	env.Outf("Service is running. pid=%s\n", pidStr)
}

func run0(command string) bool {
	switch command {
	case "help":
		help()
		return true
	case "info":
		info()
		return true
	default:
		return false
	}
}

func main() {
	// get command
	if len(env.Args) < 2 {
		env.Errf("No command specified\n")
		RunHelpMessage()
		env.Exit(1)
	}
	command := env.Args[1]

	// run0 command not need to check env
	if run0(command) {
		env.Exit(0)
	}

	// set LogDir
	env.SetLogger(LogDir)

	// run command
	env.Exit(run(command))
}

// run
func run(command string) int {
	switch command {
	case "start":
		return start()
	}

	// check pid file
	_, err := os.Stat(PidFile)
	if err != nil {
		env.Errf("Service is not running.\n")
		return 1
	}

	// check socket file
	_, err = os.Stat(SocketFile)
	if err != nil {
		env.Errf("Socket file is not found.\n")
		return 1
	}

	switch command {
	case "stop", "vault", "backup", "get", "remove", "list", "du":
		return sendCommand()
	default:
		env.Errf("Unknown command: %s\n", command)
		RunHelpMessage()
		return 1
	}
}

func sendCommand() int {
	// shift args
	args := env.Args[1:]

	// join args
	cmd := strings.Join(args, " ")

	env.EDebug("sendCommand: %s\n", cmd)

	// send command to socket
	conn, err := net.Dial("unix", SocketFile)
	if err != nil {
		env.Errf("Failed to connect to socket: %s\n", err)
		return 1
	}
	defer conn.Close()

	_, err = conn.Write([]byte(cmd))
	if err != nil {
		env.Errf("Failed to write to socket: %s\n", err)
		return 1
	}

	return 0
}

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
	)

// Usage is the help message.
//go:embed Usage.txt
var Usage string

// VaultPrefix is the prefix of the vault keys.
var VaultPrefix string

// PidFile is the path of the pid file.
const PidFile = "/var/run/backupd.pid"

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

	VaultPrefix = env.Getenv("VAULT_PREFIX", "")
	if VaultPrefix == "" {
		env.Errf("VAULT_PREFIX is not set\n")
		RunHelpMessage()
		env.Exit(1)
	}

	// run command
	env.Exit(run(command))
}

// run
func run(command string) int {
	switch command {
	case "help":
		help()
		return 0
	default:
		env.Errf("Unknown command: %s\n", command)
		RunHelpMessage()
		return 1
	}
}

// Upload(client, archive, prefix, basename, hash)
	// if IsSameHash(client, object, hash)
		// return
	// upload object

// IsSameHash(client, prefix, basename, hash) bool
	// if object exist
		// get the object's hash
		// compare hash
		// if hash is the same
			// continue

// LogModeBackup(clients, src)
	// make log archive
		// gz crypt
	// defer rm log archive
	// make hash
	// make log prefix
	// make log base name
	// loop clients
		// Upload(client, archive, prefix, basename, hash)

// DataModeBackup(clients, src)
	// make data archive
	// defer rm data archive
		// tar gz crypt
	// make hash
	// make data container name
	// make data base name
	// loop clients
		// Upload(client, archive, prefix, basename, hash)

// Run()
	// loop
		// if it is not time to backup
			// sleep 1 min
			// continue
		// Load srcConfig file
		// new vaults
		// new clients
		// get hostname
		// loop srcConfig
			// check src
			// if mode is data
				// DataModeBackup(vaults, src)
			// if mode is log
				// if src is dir
					// get files
					// loop files
						// LogModeBackup(clients, src)
				// if src is file
					// LogModeBackup(clients, src)

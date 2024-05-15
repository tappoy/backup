// This service is used to backup files to OpenStack object storage.
//
// Features:
// 	- This service will begin backing up daily at 2:00AM local time.
// 	- This service use to vault to get credencial of objectstrage. See https://github.com/tappoy/vault-cli
//
// Container Name:
// 	- backup/$HOST/NAME/data/YYYY-MM-DD.tgz
// 	- backup/$HOST/NAME/log/PATH/FILENAME.gz
//
// Environment Variables:
// 	- HOST: The hostname of the machine running this service.
package main

import (
	_ "embed"
	"github.com/tappoy/env"
	)

// Usage is the help message.
//go:embed Usage.txt
var Usage string

// RunHelpMessage prints the simple error message.
func RunHelpMessage() {
	env.Errf("Run %s help\n", env.Args[0])
}

func main() {
	RunHelpMessage()
}

// Upload(client, archive, container, object, hash)
	// if IsSameHash(client, container, object, hash)
		// return
	// upload object

// IsSameHash(client, container, object, hash) bool
	// if object exist
		// get the object's hash
		// compare hash
		// if hash is the same
			// continue

// LogModeBackup()
	// make log archive
		// gz crypt
	// defer rm log archive
	// make hash
	// make log container name
	// make log object name
	// loop vault
		// Upload(client, archive, container, object, hash)

// DataModeBackup()
	// make data archive
	// defer rm data archive
		// tar gz crypt
	// make hash
	// make data container name
	// make data object name
	// loop vault
		// Upload(client, archive, container, object, hash)

// Run()
	// loop
		// if it is not time to backup
			// sleep 1 min
			// continue
		// Load config file
		// new vaults
		// get hostname
		// Backup(hostname, vaults)
		// loop config
			// check src
			// if mode is data
				// DataModeBackup()
			// if mode is log
				// if src is dir
					// get files
					// loop files
						// LogModeBackup()
				// if src is file
					// LogModeBackup()

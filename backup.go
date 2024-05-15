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

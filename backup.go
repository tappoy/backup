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

package main

import "github.com/tappoy/env"

var vaultPrefix string
var vaultPassword string

func vault(args []string) {
	env.EInfo("vault command received")
	if len(args) != 2 {
		env.ENotice("recevie invalid vault command, %v", args)
		return
	}

	vaultPrefix = args[0]
	vaultPassword = args[1]
	env.EInfo("set vaultPrefix=(secret), vaultPassword=(secret)")
}

package main

import "github.com/tappoy/env"

var vaultDir string
var vaultPassword string

func vault(args []string) {
	env.EInfo("vault command received")
	if len(args) != 2 {
		env.ENotice("recevie invalid vault command, %v", args)
		return
	}

	vaultDir = args[0]
	vaultPassword = args[1]
	env.EInfo("set vaultDir=%s, vaultPassword=(secret)", vaultDir)
}

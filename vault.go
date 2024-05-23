package main

import "github.com/tappoy/env"

var _vaultPrefix string
var _vaultPassword string

func vault(args []string) {
	env.EInfo("vault command received")
	if len(args) != 2 {
		env.ENotice("recevie invalid vault command, %v", args)
		return
	}

	_vaultPrefix = args[0]
	_vaultPassword = args[1]
	env.EInfo("set _vaultPrefix=(secret), _vaultPassword=(secret)")
}

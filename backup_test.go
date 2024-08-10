package main

import (
	"github.com/tappoy/env"
	vaultLib "github.com/tappoy/vault"

	"os"
	"testing"
)

const testVaultDir = "tmp/vault"
const testVaultPrefix = "TEST"
const testVaultPassword = "password"

func makeVault(t *testing.T) *vaultLib.Vault {
	_vaultDir = testVaultDir
	_vaultPrefix = testVaultPrefix
	_vaultPassword = testVaultPassword

	// remove test vault directory
	os.RemoveAll(_vaultDir)

	// create test vault directory
	err := os.MkdirAll(_vaultDir, 0755)
	if err != nil {
		t.Fatalf("failed to create test vault directory: %v", err)
	}

	v, err := vaultLib.NewVault(_vaultPassword, _vaultDir)
	if err != nil {
		t.Fatalf("failed to create vault: %v", err)
	}

	// init
	err = v.Init()
	if err != nil {
		t.Fatalf("failed to init vault: %v", err)
	}

	// read mock1 config
	mock1, err := os.ReadFile("test-data/config/mock1.config")
	if err != nil {
		t.Fatalf("failed to read mock config: %v", err)
	}

	// read mock2 config
	mock2, err := os.ReadFile("test-data/config/mock2.config")
	if err != nil {
		t.Fatalf("failed to read mock config: %v", err)
	}

	// read local config
	local, err := os.ReadFile("test-data/config/local.config")
	if err != nil {
		t.Fatalf("failed to read mock config: %v", err)
	}

	// set client
	err = v.Set(_vaultPrefix+"_1", string(mock1))
	if err != nil {
		t.Fatalf("failed to set client: %v", err)
	}

	// set client
	err = v.Set(_vaultPrefix+"_2", string(mock2))
	if err != nil {
		t.Fatalf("failed to set client: %v", err)
	}

	// set client
	err = v.Set(_vaultPrefix+"_3", string(local))
	if err != nil {
		t.Fatalf("failed to set client: %v", err)
	}

	return v
}

func TestBackupByCommand(t *testing.T) {
	host, err := os.Hostname()
	if err != nil {
		t.Fatalf("failed to get hostname: %v", err)
	}
	t.Logf("host %v\n", host)

	v := makeVault(t)
	test1, err := v.Get(_vaultPrefix + "_1")
	if err != nil {
		t.Fatalf("failed to get client: %v", err)
	}
	t.Logf("vault TEST_1:\n%v", test1)

	test2, err := v.Get(_vaultPrefix + "_2")
	if err != nil {
		t.Fatalf("failed to get client: %v", err)
	}
	t.Logf("vault TEST_2:\n%v", test2)

	// set CONFIG
	env.DummyEnv["CONFIG"] = "test-data/config/backupd.config"

	// make log dir
	err = os.MkdirAll("tmp/log", 0755)
	if err != nil {
		t.Fatalf("failed to create log directory: %v", err)
	}
	env.SetLogger("tmp/log")

	// backup
	backupByCommand([]string{})
}

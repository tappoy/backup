package main

import (
	"github.com/tappoy/crypto"
	"github.com/tappoy/env"
	"github.com/tappoy/storage/v2"
	"github.com/tappoy/storage/v2/types"
	vaultLib "github.com/tappoy/vault"

	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func backupByCommand(args []string) {
	env.EInfo("backup command received")

	// check vault directory
	if _vaultDir == "" {
		env.EError("vault directory is not set")
		return
	}

	// check vault prefix
	if _vaultPrefix == "" {
		env.EError("vault prefix is not set")
		return
	}

	// check vault password
	if _vaultPassword == "" {
		env.EError("vault password is not set")
		return
	}

	// get config file path
	configFile := env.Getenv("CONFIG", DefaultConfigFile)

	srcConfig, err := LoadSrcConfigFile(configFile)
	if err != nil {
		env.EError("failed to load src config file: %s, %v", configFile, err)
		return
	}

	v, err := vaultLib.NewVault(_vaultPassword, _vaultDir)
	if err != nil {
		env.EError("failed to create vaults: %v", err)
		return
	}

	// get clinets
	clients := make([]types.Client, 0)
	for n := 1; true; n++ {
		key := fmt.Sprintf("%s_%d", _vaultPrefix, n)
		clientConfig, err := v.Get(key)
		if err != nil {
			env.EInfo("not found. n:%d", n)
			break
		}

		client, err := storage.NewClientFromString(clientConfig)
		if err != nil {
			env.ENotice("failed to create client: %v", err)
			break
		}
		clients = append(clients, client)
	}

	// src loop
	for _, sc := range srcConfig {
		err = backup(sc, clients)
		if err != nil {
			env.ENotice("failed to backup: %v, %v", err, sc.Name)
		}
	}
}

func makePrefix(sc SrcConfig) string {
	host := env.Hostname("UNKNOWN_HOST")
	return fmt.Sprintf("%s/%s/%s/%s/", host, sc.Name, sc.Mode, sc.Source)
}

func dataModeBackup(sc SrcConfig, clients []types.Client) error {
	// create tmp file
	tmpFile, err := os.CreateTemp("", "backupd")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	// tar gz crypt
	err = crypto.TarGzEncrypto(sc.Source, tmpFile, _vaultPassword)
	if err != nil {
		return err
	}

	// prefix
	prefix := makePrefix(sc)

	// basename
	basename := fmt.Sprintf("%s.tgzc", time.Now().Format("2006-01-02"))

	// body
	body, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return err
	}

	// objectname
	objectname := filepath.Clean(prefix + basename)

	upload(objectname, body, clients)
	return nil
}

func logModeBackupFile(sc SrcConfig, file string, clients []types.Client) error {
	// create tmp file
	tmpFile, err := os.CreateTemp("", "backupd")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	// open file
	fr, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fr.Close()

	// gz crypt
	err = crypto.GzEncrypto(fr, tmpFile, _vaultPassword)
	if err != nil {
		return err
	}

	// prefix
	prefix := makePrefix(sc)

	// basename
	basename := fmt.Sprintf("%s.gzc", filepath.Base(file))

	// body
	body, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return err
	}

	// upload
	upload(prefix+basename, body, clients)
	return nil
}

func logModeBackup(sc SrcConfig, clients []types.Client) error {
	// get src stat
	stat, err := os.Stat(sc.Source)
	if err != nil {
		return err
	}

	files := make([]string, 0)

	// if src is dir
	if stat.IsDir() {
		err = filepath.Walk(sc.Source, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// check dir
			if info.IsDir() {
				return nil
			}

			files = append(files, path)
			return nil
		})
	} else {
		files = append(files, sc.Source)
	}

	// loop
	for _, file := range files {
		err = logModeBackupFile(sc, file, clients)
		if err != nil {
			env.ENotice("failed to backup: %v, %v", err, file)
		}
	}
	return nil
}

func upload(name string, body []byte, clients []types.Client) {
	for _, client := range clients {
		// head
		info, err := client.Head(name)
		if err != nil {
			if err != types.ErrNotFound {
				// Head: ng, !ErrNotFound
				env.ENotice("%v: failed to head: %v", client, err)
				continue
			}
			// Head: ng, ErrNotFound
		}

		if err != types.ErrNotFound {
			// Head: ok, nil
			// md5
			md5 := crypto.Md5(string(body))

			// check hash
			if info.Hash == md5 {
				env.EInfo("%v: same hash: %s", client, name)
				continue
			}
		} else {
			// Head: ng, ErrNotFound
			// do nothing
		}

		// reader
		reader := bytes.NewReader(body)

		// upload
		err = client.Put(name, reader)
		if err != nil {
			env.ENotice("%v: failed to put: %v", client, err)
			continue
		}
		env.EInfo("%v: uploaded: %s", client, name)
	}
}

func backup(sc SrcConfig, clients []types.Client) error {
	switch sc.Mode {
	case DataMode:
		env.EInfo("data mode backup: %s", sc.Source)
		return dataModeBackup(sc, clients)
	case LogMode:
		env.EInfo("log mode backup: %s", sc.Source)
		return logModeBackup(sc, clients)
	default:
		return fmt.Errorf("invalid mode: %s", sc.Mode)
	}
}

func backupBySchedule(c chan error, cmd chan string, sigterm chan os.Signal) {
	for {
		select {
		case <-sigterm:
			env.EInfo("backup interrupted by signal")
			c <- ErrInterrupted
			return
		case line := <-cmd:
			// execute command
			env.EInfo("command received while backup: %s", line)
			if line == "stop" {
				env.EInfo("backup interrupted by stop command")
				c <- ErrInterrupted
				return
			} else {
				env.EInfo("command ignored while backup: %s", line)
			}
		default:
			// do something
			time.Sleep(1 * time.Second)
			env.Errf(".")
		}
	}
}

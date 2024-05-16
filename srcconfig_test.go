package main

import (
	"github.com/tappoy/tsv"

	"errors"
	"testing"
	"reflect"
	"os"
	)

func TestLoadSrcConfigFile(t *testing.T) {
	config, err := LoadSrcConfigFile("test-data/config/backupd.config")
	if err != nil {
		t.Errorf("err: %v", err)
	}

	want := []SrcConfig{
		{"test-vault", DataMode, "test-data/srv/vault"},
		{"test-vault", LogMode, "test-data/var/log/vault"},
		{"test-config", DataMode, "test-data/config"},
	}

	if !reflect.DeepEqual(config, want) {
		t.Errorf("got %v, want %v", config, want)
	}

}

func TestLoadSrcConfigFileNoExist(t *testing.T) {
	_, err := LoadSrcConfigFile("test-data/config/no-exist.config")
	if err == nil {
		t.Errorf("err: %v", err)
	}
}

func TestLoadSrcConfigFileError(t *testing.T) {
	// invalid number of fields
	f, err := os.Create("tmp/invalid-number-of-fields.config")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer f.Close()

	f.WriteString("test-vault\tdata\ttest-data/srv/vault\texfield\n")

	_, err = LoadSrcConfigFile("tmp/invalid-number-of-fields.config")
	want := SrcConfigError{
		TsvLine: tsv.TsvLine{
			LineNo: 1,
			Line: "test-vault\tdata\ttest-data/srv/vault\texfield",
			Fields: []string{"test-vault", "data", "test-data/srv/vault", "exfield"},
		},
		Err: errors.New("invalid number of fields. it must be 3"),
	}

	if !reflect.DeepEqual(err, want) {
		t.Errorf("got %v, want %v", err, want)
	}
}

func TestLoadSrcConfigFileError2(t *testing.T) {
	// empty name
	fileName := "tmp/empty-name.config"
	f, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer f.Close()

	f.WriteString("\tdata\ttest-data/srv/vault\n")

	_, err = LoadSrcConfigFile(fileName)
	want := SrcConfigError{
		TsvLine: tsv.TsvLine{
			LineNo: 1,
			Line: "\tdata\ttest-data/srv/vault",
			Fields: []string{"", "data", "test-data/srv/vault"},
		},
		Err: errors.New("name is empty"),
	}

	if !reflect.DeepEqual(err, want) {
		t.Errorf("got %v, want %v", err, want)
	}
}

func TestLoadSrcConfigFileError3(t *testing.T) {
	// invalid mode
	fileName := "tmp/invalid-mode.config"
	f, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer f.Close()

	f.WriteString("test-vault\tinvalid-mode\ttest-data/srv/vault\n")

	_, err = LoadSrcConfigFile(fileName)
	want := SrcConfigError{
		TsvLine: tsv.TsvLine{
			LineNo: 1,
			Line: "test-vault\tinvalid-mode\ttest-data/srv/vault",
			Fields: []string{"test-vault", "invalid-mode", "test-data/srv/vault"},
		},
		Err: errors.New("invalid mode. it must be data or log"),
	}

	if !reflect.DeepEqual(err, want) {
		t.Errorf("got %v, want %v", err, want)
	}
}

func TestLoadSrcConfigFileError4(t *testing.T) {
	// source directory is empty
	fileName := "tmp/source-directory-is-empty.config"
	f, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer f.Close()

	f.WriteString("test-vault\tdata\t\n")

	_, err = LoadSrcConfigFile(fileName)
	want := SrcConfigError{
		TsvLine: tsv.TsvLine{
			LineNo: 1,
			Line: "test-vault\tdata\t",
			Fields: []string{"test-vault", "data", ""},
		},
		Err: errors.New("source directory is empty"),
	}

	if !reflect.DeepEqual(err, want) {
		t.Errorf("got %v, want %v", err, want)
	}

	wantMessage := "line 1: test-vault\tdata\t, source directory is empty"
	if err.Error() != wantMessage {
		t.Errorf("got %v, want %v", err.Error(), wantMessage)
	}
}

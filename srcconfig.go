package main

import (
	"github.com/tappoy/tsv"

	"errors"
	"fmt"
	)

// BackupMode is the kind of backup mode.
type BackupMode string

var (
	// Tar gz the target directory and keep 14 recent backups.
	// The backup file name format is YYYY-MM-DD.tgz.
	// If the latest backups has the same hash, that backup is skipped.
	DataMode BackupMode = "data"

	// Recursively gzip compresses files in the target directory
	// and keeps all backups. This backup never expires.
	// If a file with the same name already exists in the container,
	// it is skipped.
	LogMode BackupMode = "log"
)

// SrcConfig represents backup settings defined by a yaml file.
type SrcConfig struct {
	// Name is the name of the setting.
	// It uses the container name in the bucket.
	Name string

	// Mode is the mode of backup.
	Mode BackupMode

	// Source is the source directory path to be backed up.
	Source string
}

// SrcConfigError represents an error in the config file.
type SrcConfigError struct {
	TsvLine tsv.TsvLine
	Err error
}

func (e SrcConfigError) Error() string {
	return fmt.Sprintf("line %d: %s, %s", e.TsvLine.LineNo, e.TsvLine.Line, e.Err)
}

func configError(msg string, line tsv.TsvLine) SrcConfigError {
	return SrcConfigError{
		line,
		errors.New(msg),
	}
}

// LoadSrcConfigFile loads backup settings from config file.
func LoadSrcConfigFile(path string) ([]SrcConfig, error) {
	tsvLines, err := tsv.ReadTsvFile(path)
	if err != nil {
		return nil, err
	}

	configs := []SrcConfig{}

	for _, tsvLine := range tsvLines {
		if len(tsvLine.Fields) != 3 {
			return nil, configError("invalid number of fields. it must be 3", tsvLine)
		}
		name := tsvLine.Fields[0]
		mode := BackupMode(tsvLine.Fields[1])
		src := tsvLine.Fields[2]
		if len(name) == 0 {
			return nil, configError("name is empty", tsvLine)
		}
		if mode != DataMode && mode != LogMode {
			return nil, configError("invalid mode. it must be data or log", tsvLine)
		}
		if len(src) == 0 {
			return nil, configError("source directory is empty", tsvLine)
		}

		configs = append(configs, SrcConfig{
			Name: name,
			Mode: mode,
			Source: src,
		})
	}

	return configs, nil
}

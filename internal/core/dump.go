package core

import (
	"os"
)

type DumpFile struct {
	Path string
}

func NewDumpFile(path string) DumpFile {
	return DumpFile{
		Path: path,
	}
}

func (d *DumpFile) Exists() bool {
	_, err := os.Stat(d.Path)
	return err == nil
}

func (d *DumpFile) FileName() (string, error) {
	i, err := os.Stat(d.Path)
	if err != nil {
		return "", err
	}
	n := i.Name()
	return n, nil
}

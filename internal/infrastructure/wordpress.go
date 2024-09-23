package infrastructure

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type WordpressImpl struct {
	contentDirPath string
	logger         *log.Logger
}

// NewWordPressImpl creates a new WordPressImpl from the given content directory path
func NewWordPressImpl(p string) *WordpressImpl {
	l := log.New(os.Stdout, "WordPressImpl", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	return &WordpressImpl{contentDirPath: p, logger: l}
}

// CreateBackup creates a backup of the WordPress content directory
func (w *WordpressImpl) CreateBackup(name string) (string, error) {
	destFile := w.getArchiveFilePath(name)

	w.logger.Printf("Creating backup %s\n...", destFile)

	file, err := os.Create(destFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	err = filepath.Walk(w.contentDirPath, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		header.Name = filepath.Join(filepath.Base(w.contentDirPath), file)

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !fi.Mode().IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			defer data.Close()

			_, err = io.Copy(tw, data)
			if err != nil {
				return err
			}
		}

		return err
	})

	if err != nil {
		w.CleanupArchive(name)
		return "", err
	}

	w.logger.Printf("Content backup file %s created successfully\n", destFile)

	return destFile, nil
}

// CleanupArchive removes the archive local file forthe given backup name
func (w *WordpressImpl) CleanupArchive(backupName string) error {
	return os.Remove(w.getArchiveFilePath(backupName))
}

func (w *WordpressImpl) getArchiveFilePath(backupName string) string {
	return fmt.Sprintf("%s.tar.gz", backupName)
}

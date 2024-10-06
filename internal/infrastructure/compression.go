package infrastructure

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CompressFolder(src, dest string) error {
	logger.Printf("Starting compress source folder %s to %s\n...", src, dest)

	destDir := filepath.Dir(dest)
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		err = os.MkdirAll(destDir, 0755)
		if err != nil {
			return err
		}
		err = os.Chmod(destDir, 0755)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	err = filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, file)
		if err != nil {
			return err
		}
		header.Name = relPath

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
		cleanupArchive(dest)
		return err
	}

	logger.Printf("folder %s successfully copressed into archive file %s\n", src, dest)

	return nil
}

func UncompressFolder(src, dest string) error {
	logger.Printf("Starting uncompress source archive %s to %s\n...", src, dest)

	if ok := fileExists(src); !ok {
		return fmt.Errorf("file %s does not exists", src)
	}

	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			f.Close()
		}
	}
	return nil
}

func cleanupArchive(path string) error {
	return os.Remove(path)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

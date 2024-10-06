package infrastructure

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type MariaDbImpl struct {
	username      string
	password      string
	containerName string
	logger        *log.Logger
}

// NewMariaDbImpl creates a new MariaDbImpl from the given username, password and container name
func NewMariaDbImpl(u, p, c string) *MariaDbImpl {
	l := log.New(os.Stdout, "MariaDbImpl", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	return &MariaDbImpl{u, p, c, l}
}

// BackupDatabaseFromDocker creates a backup of the database from the given container
func (db *MariaDbImpl) BackupDatabaseFromDocker(backupPath string) error {
	cmd := exec.Command("docker", "exec", db.containerName, "mariadb-dump", "--all-databases", fmt.Sprintf("-u%s", db.username), fmt.Sprintf("-p%s", db.password))
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		return err
	}

	err = os.WriteFile(backupPath, output, 0644)
	if err != nil {
		return err
	}

	return nil
}

// RestoreDatabaseFromDocker restores the database from the given dump file
func (db *MariaDbImpl) RestoreDatabaseFromDocker(dumpFilePath string) error {
	_, err := os.Stat(dumpFilePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("dump file %s does not exists", dumpFilePath)
	}

	restoreCmd := fmt.Sprintf("exec mariadb -u%s -p%s", db.username, db.password)
	dumpFile, err := os.Open(dumpFilePath)
	if err != nil {
		return err
	}
	cmd := exec.Command("docker", "exec", "-i", db.containerName, "sh", "-c", restoreCmd)
	cmd.Stdin = dumpFile
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

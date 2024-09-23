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
func (db *MariaDbImpl) BackupDatabaseFromDocker(backupName string) (string, error) {
	dumpFilePath := fmt.Sprintf("%s.sql", backupName)
	cmd := exec.Command("docker", "exec", db.containerName, "mariadb-dump", "--all-databases", fmt.Sprintf("-u%s", db.username), fmt.Sprintf("-p%s", db.password))
	db.logger.Printf("Running command: %s\n", cmd.String())
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	err = os.WriteFile(dumpFilePath, output, 0644)
	if err != nil {
		return "", err
	}

	return dumpFilePath, nil
}

// RestoreDatabaseFromDocker restores the database from the given dump file
func (db *MariaDbImpl) RestoreDatabaseFromDocker(dumpFilePath string) error {
	_, err := os.Stat(dumpFilePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("dump file %s does not exists", dumpFilePath)
	}

	restoreCmd := fmt.Sprintf("mysql -u%s -p%s", db.username, db.password)
	cmd := exec.Command("docker", "exec", "-i", db.containerName, "sh", "-c", fmt.Sprintf("'%s'", restoreCmd), "<", dumpFilePath)
	cmd.Stdin = os.Stdin
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

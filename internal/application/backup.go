package application

import (
	"errors"

	"github.com/thomas-marquis/wordpress-simple-backup/internal/core"
)

type BackupApplication struct {
	repo core.Repository
}

func NewBackupApplication(repo core.Repository) *BackupApplication {
	return &BackupApplication{
		repo: repo,
	}
}

func (a *BackupApplication) Save(name string) error {
	existingBackup, err := a.repo.GetExistingBackup(name)
	if err != nil {
		return errors.New("an error occurred when getting previous backup")
	}

	newVersion, err := a.repo.CreateNewBackupVersion(name)
	if err != nil {
		return errors.New("an error occurred when creating new backup version")
	}

	if !newVersion.AreLocalBackupsAllReady() {
		return errors.New("backup files are not ready")
	}

	existingBackup.AddVersion(*newVersion)

	if err := a.repo.SaveBackup(existingBackup); err != nil {
		return errors.New("an error occurred when saving backup")
	}
	return nil
}

func (a *BackupApplication) Restore(backupName string, versionID int) error {
	// Get backup

	// Get version if exists

	// Restore version
	return errors.New("not implemented")
}

func (a *BackupApplication) List() ([]core.Version, error) {
	return []core.Version{}, errors.New("not implemented")
}

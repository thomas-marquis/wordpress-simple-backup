package application

import (
	"errors"
	"fmt"

	"github.com/thomas-marquis/wordpress-simple-backup/internal/core"
)

func (a *BackupApplication) Save() error {
	logger.Printf("Starting backup process...")

	logger.Printf("Creating content dump...")
	contentDump, err := a.repo.CreateContentDump()
	if err != nil {
		return errors.New("an error occurred when creating content dump")
	}
	if !contentDump.Exists() {
		return fmt.Errorf("the dump file %s does not exist", contentDump.Path)
	}
	defer func() {
		if err := a.repo.ClearDump(contentDump); err != nil {
			fmt.Printf("an error occurred when clearing content dump: %v", err)
		}
	}()
	logger.Printf("Content dump created")

	logger.Printf("Creating db dump...")
	dbDump, err := a.repo.CreateDbDump()
	if err != nil {
		return errors.New("an error occurred when creating db dump")
	}
	if !dbDump.Exists() {
		return fmt.Errorf("the dump file %s does not exist", dbDump.Path)
	}
	defer func() {
		if err := a.repo.ClearDump(dbDump); err != nil {
			fmt.Printf("an error occurred when clearing db dump: %v", err)
		}
	}()
	logger.Printf("Db dump created")

	versions, err := a.repo.ListVersions()
	if err != nil {
		return fmt.Errorf("an error occurred when listing versions: %v", err)
	}
	logger.Printf("Existing versions: %s", versionListToString(versions))

	var higherID core.VersionID
	for _, v := range versions {
		if v.ID > higherID {
			higherID = v.ID
		}
	}
	newVersion := core.NewVersion(higherID+1, dbDump, contentDump)
	logger.Printf("Creating new version: %s", newVersion.String())

	logger.Printf("Saving new version...")
	if err = a.repo.SaveVersion(newVersion); err != nil {
		return errors.New("an error occurred when saving version")
	}
	logger.Printf("New version saved")

	versions = append(versions, newVersion)

	if len(versions) > a.versionsToKeep {
		nbVersionsToRemove := len(versions) - a.versionsToKeep
		logger.Printf("Too many versions are stored, removing the %d oldest ones...", nbVersionsToRemove)
		for i := 0; i < nbVersionsToRemove; i++ {
			if err = a.repo.RemoveVersion(versions[i].ID); err != nil {
				return fmt.Errorf("an error occurred when removing version: %s", err)
			}
		}
		logger.Printf("Old versions removed")
	}

	logger.Printf("Backup process completed")

	return nil
}

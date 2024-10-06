package application

import (
	"errors"
	"fmt"

	"github.com/thomas-marquis/wordpress-simple-backup/internal/core"
)

func (a *BackupApplication) Save() error {
	contentDump, err := a.repo.CreateContentDump()
	if err != nil {
		return errors.New("an error occurred when creating content dump")
	}
	if !contentDump.Exists() {
		return fmt.Errorf("the dump file %s does not exist", contentDump.Path)
	}
	defer func() {
		err := a.repo.ClearDump(contentDump)
		if err != nil {
			fmt.Printf("an error occurred when clearing content dump: %v", err)
		}
	}()

	dbDump, err := a.repo.CreateDbDump()
	if err != nil {
		return errors.New("an error occurred when creating db dump")
	}
	if !dbDump.Exists() {
		return fmt.Errorf("the dump file %s does not exist", dbDump.Path)
	}
	defer func() {
		err := a.repo.ClearDump(dbDump)
		if err != nil {
			fmt.Printf("an error occurred when clearing db dump: %v", err)
		}
	}()

	versions, err := a.repo.ListVersions()
	if err != nil {
		return errors.New("an error occurred when listing versions")
	}

	var higherID core.VersionID
	for _, v := range versions {
		if v.ID > higherID {
			higherID = v.ID
		}
	}
	newVersion := core.NewVersion(higherID+1, dbDump, contentDump)

	err = a.repo.SaveVersion(newVersion)
	if err != nil {
		return errors.New("an error occurred when saving version")
	}

	versions = append(versions, newVersion)

	if len(versions) > a.versionsToKeep {
		nbVersionsToRemove := len(versions) - a.versionsToKeep
		for i := 0; i < nbVersionsToRemove; i++ {
			err = a.repo.RemoveVersion(versions[i].ID)
			if err != nil {
				return errors.New("an error occurred when removing version")
			}
		}
	}

	return nil
}

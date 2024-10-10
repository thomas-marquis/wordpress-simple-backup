package application

import (
	"fmt"

	"github.com/thomas-marquis/wordpress-simple-backup/internal/core"
)

func (a *BackupApplication) List() ([]*core.Version, error) {
	v, err := a.repo.ListVersions()
	if err != nil {
		return nil, fmt.Errorf("an error occurred when listing versions: %S", err)
	}
	return v, nil
}

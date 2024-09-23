package core_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/core"
)

func Test_AreLocalBackupsAllReady_ShouldReturnTrueWhenBothFilesExist(t *testing.T) {
	// Given
	dbFile, err := os.CreateTemp("", "test_db.sql")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dbFile.Name())

	contentFile, err := os.CreateTemp("", "test_content.tar.gz")
	if err != nil {
		t.Fatal(err)
	}

	v := core.NewVersion(1)
	v.WithLocalDbBackupFile(dbFile.Name())
	v.WithLocalContentBackupFile(contentFile.Name())

	// When
	result := v.AreLocalBackupsAllReady()

	// Then
	assert.True(t, result)
}

func Test_AreLocalBackupsAllReady_ShouldReturnFalseWhenDbFileIsMissing(t *testing.T) {
	// Given
	contentFile, err := os.CreateTemp("", "test_content.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(contentFile.Name())

	v := core.NewVersion(1)
	v.WithLocalContentBackupFile(contentFile.Name())

	// When
	result := v.AreLocalBackupsAllReady()

	// Then
	assert.False(t, result)
}

func Test_AreLocalBackupsAllReady_ShouldReturnFalseWhenContentFileIsMissing(t *testing.T) {
	// Given
	dbFile, err := os.CreateTemp("", "test_db.sql")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dbFile.Name())

	v := core.NewVersion(1)
	v.WithLocalDbBackupFile(dbFile.Name())

	// When
	result := v.AreLocalBackupsAllReady()

	// Then
	assert.False(t, result)
}

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
	defer os.Remove(contentFile.Name())

	content := core.NewDumpFile(contentFile.Name())
	db := core.NewDumpFile(dbFile.Name())

	v := core.NewVersion(1, db, content)

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

	content := core.NewDumpFile(contentFile.Name())
	db := core.NewDumpFile("fakepath")

	v := core.NewVersion(1, db, content)

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

	content := core.NewDumpFile("fakepath")
	db := core.NewDumpFile(dbFile.Name())

	v := core.NewVersion(1, db, content)

	// When
	result := v.AreLocalBackupsAllReady()

	// Then
	assert.False(t, result)
}

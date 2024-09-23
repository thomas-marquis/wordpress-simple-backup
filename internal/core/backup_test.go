package core_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/core"
)

func Test_AddVersion_ShouldAddVersionToBackup(t *testing.T) {
	// Given
	b := core.NewBackup("test", 1)
	version := &core.Version{}

	// When
	b.AddVersion(*version)

	// Then
	assert.Equal(t, 1, len(b.Versions))
}

func Test_AddVersion_ShouldRemoveOldVersions(t *testing.T) {
	// Given
	b := core.NewBackup("test", 1)
	version1 := &core.Version{}
	version1.CreatedAt = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	version2 := &core.Version{}
	version2.CreatedAt = time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)

	// When
	b.AddVersion(*version1)
	b.AddVersion(*version2)

	// Then
	assert.Equal(t, 1, len(b.Versions))
	assert.Equal(t, *version2, b.Versions[0])
}

package core_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/core"
)

func Test_Exists_ShouldReturnTrue(t *testing.T) {
	// Given
	file, err := os.CreateTemp("", "db-*.sql")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	dump := core.NewDumpFile(file.Name())

	// When
	exists := dump.Exists()

	// Then
	assert.True(t, exists)
}

func Test_Exists_ShouldReturnFalse(t *testing.T) {
	// Given
	dump := core.NewDumpFile("fakepath")

	// When
	exists := dump.Exists()

	// Then
	assert.False(t, exists)
}

func Test_FileName_ShouldReturnFileName(t *testing.T) {
	// Given
	file, err := os.CreateTemp("", "db-toto.sql")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	i, err := os.Stat(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	expectedName := i.Name()
	dump := core.NewDumpFile(file.Name())

	// When
	name, err := dump.FileName()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedName, name)
}

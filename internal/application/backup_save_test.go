package application_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/application"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/core"
	mocks_core "github.com/thomas-marquis/wordpress-simple-backup/mocks"
	"go.uber.org/mock/gomock"
)

func getFakeBackupFiles(t *testing.T) (*os.File, *os.File, func()) {
	dbFile, err := os.CreateTemp("", "db-*.sql")
	if err != nil {
		t.Fatal(err)
	}

	contentFile, err := os.CreateTemp("", "content-*.tar.gz")
	if err != nil {
		t.Fatal(err)
	}

	return dbFile, contentFile, func() {
		os.Remove(dbFile.Name())
		os.Remove(contentFile.Name())
	}
}

func Test_Save_ShouldSaveWhenNewBackupIsCreatedAndExistingAlreadyExists(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	repo := mocks_core.NewMockRepository(ctrl)
	app := application.NewBackupApplication(repo, 100)

	dbFile, contentFile, cleanup := getFakeBackupFiles(t)
	defer cleanup()

	versions := []*core.Version{
		{ID: 1},
		{ID: 2},
	}
	newVersion := &core.Version{ID: 3}

	repo.EXPECT().CreateContentDump().Return(core.NewDumpFile(contentFile.Name()), nil).Times(1)
	repo.EXPECT().CreateDbDump().Return(core.NewDumpFile(dbFile.Name()), nil).Times(1)
	repo.EXPECT().ListVersions().Return(versions, nil).Times(1)
	repo.EXPECT().SaveVersion(gomock.Any()).Do(func(v *core.Version) error {
		assert.Equal(t, newVersion.ID, v.ID)
		assert.Equal(t, dbFile.Name(), v.DbBackupFile().Path)
		assert.Equal(t, contentFile.Name(), v.ContentBackupFile().Path)
		return nil
	}).Return(nil).Times(1)
	repo.EXPECT().ClearDump(core.NewDumpFile(contentFile.Name())).Return(nil).Times(1)
	repo.EXPECT().ClearDump(core.NewDumpFile(dbFile.Name())).Return(nil).Times(1)

	// When
	err := app.Save()

	// Then
	assert.NoError(t, err)
}

func Test_Save_ShouldRemoveOldestVersionWhenLimitWasReached(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	repo := mocks_core.NewMockRepository(ctrl)
	app := application.NewBackupApplication(repo, 3)

	dbFile, contentFile, cleanup := getFakeBackupFiles(t)
	defer cleanup()

	versions := []*core.Version{
		{ID: 1},
		{ID: 2},
		{ID: 3},
	}
	newVersion := &core.Version{ID: 4}

	repo.EXPECT().CreateContentDump().Return(core.NewDumpFile(contentFile.Name()), nil).Times(1)
	repo.EXPECT().CreateDbDump().Return(core.NewDumpFile(dbFile.Name()), nil).Times(1)
	repo.EXPECT().ListVersions().Return(versions, nil).Times(1)
	repo.EXPECT().SaveVersion(gomock.Any()).Do(func(v *core.Version) error {
		assert.Equal(t, newVersion.ID, v.ID)
		assert.Equal(t, dbFile.Name(), v.DbBackupFile().Path)
		assert.Equal(t, contentFile.Name(), v.ContentBackupFile().Path)
		return nil
	}).Return(nil).Times(1)
	repo.EXPECT().RemoveVersion(core.VersionID(1)).Return(nil).Times(1)
	repo.EXPECT().ClearDump(gomock.Any()).Return(nil).AnyTimes()

	// When
	err := app.Save()

	// Then
	assert.NoError(t, err)
}

func Test_Save_ShouldRemoveManyOlderVersionIfNeeded(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	repo := mocks_core.NewMockRepository(ctrl)
	app := application.NewBackupApplication(repo, 3)

	dbFile, contentFile, cleanup := getFakeBackupFiles(t)
	defer cleanup()

	versions := []*core.Version{
		{ID: 1},
		{ID: 2},
		{ID: 3},
		{ID: 4},
	}

	repo.EXPECT().CreateContentDump().Return(core.NewDumpFile(contentFile.Name()), nil).Times(1)
	repo.EXPECT().CreateDbDump().Return(core.NewDumpFile(dbFile.Name()), nil).Times(1)
	repo.EXPECT().ListVersions().Return(versions, nil).Times(1)
	repo.EXPECT().SaveVersion(gomock.Any()).Return(nil).Times(1)

	repo.EXPECT().RemoveVersion(gomock.Any()).DoAndReturn(func(id core.VersionID) error {
		if id != core.VersionID(1) && id != core.VersionID(2) {
			t.Errorf("Unexpected id: %v", id)
		}
		return nil
	}).Times(2)
	repo.EXPECT().ClearDump(gomock.Any()).Return(nil).AnyTimes()

	// When
	err := app.Save()

	// Then
	assert.NoError(t, err)
}

func Test_Save_ShouldClearDumpsWhenAnErrorOccurredDuringListVersions(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	repo := mocks_core.NewMockRepository(ctrl)
	app := application.NewBackupApplication(repo, 3)

	dbFile, contentFile, cleanup := getFakeBackupFiles(t)
	defer cleanup()

	repo.EXPECT().CreateContentDump().Return(core.NewDumpFile(contentFile.Name()), nil).Times(1)
	repo.EXPECT().CreateDbDump().Return(core.NewDumpFile(dbFile.Name()), nil).Times(1)
	repo.EXPECT().ListVersions().Return(nil, assert.AnError).Times(1)
	repo.EXPECT().ClearDump(core.NewDumpFile(contentFile.Name())).Return(nil).Times(1)
	repo.EXPECT().ClearDump(core.NewDumpFile(dbFile.Name())).Return(nil).Times(1)

	// When
	err := app.Save()

	// Then
	assert.Error(t, err)
}

func Test_Save_ShouldNotSaveNewVersionWhenOneDumpPathIsWrong(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	repo := mocks_core.NewMockRepository(ctrl)
	app := application.NewBackupApplication(repo, 3)

	_, contentFile, cleanup := getFakeBackupFiles(t)
	defer cleanup()

	repo.EXPECT().CreateContentDump().Return(core.NewDumpFile(contentFile.Name()), nil).Times(1)
	repo.EXPECT().CreateDbDump().Return(core.NewDumpFile("wrongfilepath"), nil).Times(1)
	repo.EXPECT().ClearDump(core.NewDumpFile(contentFile.Name())).Return(nil).Times(1)
	repo.EXPECT().SaveVersion(gomock.Any()).Times(0)
	repo.EXPECT().RemoveVersion(gomock.Any()).Times(0)

	// When
	err := app.Save()

	// Then
	assert.Error(t, err)
	assert.Equal(t, "the dump file wrongfilepath does not exist", err.Error())
}

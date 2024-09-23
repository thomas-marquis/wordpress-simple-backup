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

func makeFakeCreateVersionImpl(v *core.Version, dbFile *os.File, contentFile *os.File) func(name string) (*core.Version, error) {
	return func(name string) (*core.Version, error) {
		if dbFile != nil {
			v.WithLocalDbBackupFile(dbFile.Name())
		}
		if contentFile != nil {
			v.WithLocalContentBackupFile(contentFile.Name())
		}
		return v, nil
	}
}

func Test_Save_ShouldSaveBackupProcess(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	repo := mocks_core.NewMockRepository(ctrl)
	app := application.NewBackupApplication(repo)
	backupName := "test"

	dbFile, contentFile, cleanup := getFakeBackupFiles(t)
	defer cleanup()

	// Get existing backup from repository
	existingBackup := &core.Backup{}
	existingBackup.VersionsToKeep = 1
	repo.EXPECT().GetExistingBackup(backupName).Return(existingBackup, nil).Times(1)

	// Create new backup version
	newVersion := &core.Version{ID: 1}
	repo.EXPECT().CreateNewBackupVersion(backupName).
		DoAndReturn(makeFakeCreateVersionImpl(newVersion, dbFile, contentFile)).
		Times(1)

	// Save new backup version to repository
	repo.EXPECT().SaveBackup(existingBackup).Do(func(bu *core.Backup) {
		assert.Equal(t, 1, len(bu.Versions))
		assert.Equal(t, 1, bu.Versions[0].ID)
	}).Return(nil).Times(1)

	// When
	err := app.Save(backupName)

	// Then
	assert.NoError(t, err)
}

func Test_Save_ShouldReturnErrorWhenGetExistingBackupFails(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	repo := mocks_core.NewMockRepository(ctrl)
	app := application.NewBackupApplication(repo)
	backupName := "test"

	repo.EXPECT().GetExistingBackup(backupName).Return(nil, assert.AnError).Times(1)
	repo.EXPECT().SaveBackup(gomock.Any()).Times(0)

	// When
	err := app.Save(backupName)

	// Then
	assert.Error(t, err)
	assert.EqualError(t, err, "an error occurred when getting previous backup")
}

func Test_Save_ShouldReturnErrorWhenCreateBackupVersionFails(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	repo := mocks_core.NewMockRepository(ctrl)
	app := application.NewBackupApplication(repo)
	backupName := "test"

	repo.EXPECT().GetExistingBackup(gomock.Any()).Return(&core.Backup{}, nil).AnyTimes()

	repo.EXPECT().CreateNewBackupVersion(backupName).Return(nil, assert.AnError).Times(1)
	repo.EXPECT().SaveBackup(gomock.Any()).Times(0)

	// When
	err := app.Save(backupName)

	// Then
	assert.Error(t, err)
	assert.EqualError(t, err, "an error occurred when creating new backup version")
}

func Test_Save_ShouldReturnErrorWhenSaveBackupFails(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	repo := mocks_core.NewMockRepository(ctrl)
	app := application.NewBackupApplication(repo)
	backupName := "test"

	dbFile, contentFile, cleanup := getFakeBackupFiles(t)
	defer cleanup()

	repo.EXPECT().GetExistingBackup(gomock.Any()).Return(&core.Backup{}, nil).AnyTimes()
	repo.EXPECT().CreateNewBackupVersion(gomock.Any()).
		DoAndReturn(makeFakeCreateVersionImpl(&core.Version{}, dbFile, contentFile)).
		AnyTimes()

	repo.EXPECT().SaveBackup(gomock.Any()).Return(assert.AnError).Times(1)

	// When
	err := app.Save(backupName)

	// Then
	assert.Error(t, err)
	assert.EqualError(t, err, "an error occurred when saving backup")
}

func Test_Save_ShouldReturnErrorWhenBackupFilesNotReadyBeforeSaving(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	repo := mocks_core.NewMockRepository(ctrl)
	app := application.NewBackupApplication(repo)
	backupName := "test"

	dbFile, _, cleanup := getFakeBackupFiles(t)
	defer cleanup()

	repo.EXPECT().GetExistingBackup(backupName).Return(&core.Backup{}, nil).Times(1)
	repo.EXPECT().CreateNewBackupVersion(backupName).
		DoAndReturn(makeFakeCreateVersionImpl(&core.Version{}, dbFile, nil)).
		Times(1)

	repo.EXPECT().SaveBackup(gomock.Any()).Times(0)

	// When
	err := app.Save(backupName)

	// Then
	assert.Error(t, err)
	assert.EqualError(t, err, "backup files are not ready")
}

func Test_Restore_ShouldExectuteNominalRestorationProcess(t *testing.T) {
	// Given

	// When

	// Then
}

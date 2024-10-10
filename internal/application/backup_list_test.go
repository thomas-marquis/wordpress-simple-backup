package application_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/application"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/core"
	mocks_core "github.com/thomas-marquis/wordpress-simple-backup/mocks"
	"go.uber.org/mock/gomock"
)

func Test_List_ShouldListExistingVersions(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	repo := mocks_core.NewMockRepository(ctrl)
	app := application.NewBackupApplication(repo, 100)

	versions := []*core.Version{
		{ID: 1, CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
		{ID: 2, CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)},
	}

	repo.EXPECT().ListVersions().Return(versions, nil).Times(1)

	// When
	res, err := app.List()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, versions, res)
}

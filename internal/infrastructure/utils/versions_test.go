package infrautils_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/core"
	infrautils "github.com/thomas-marquis/wordpress-simple-backup/internal/infrastructure/utils"
)

func Test_ParseVersionDirName(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected core.Version
	}{
		{
			Input:    "mysiteBackups/1x2024-12-31_23-59-49",
			Expected: core.Version{ID: 1, CreatedAt: time.Date(2024, 12, 31, 23, 59, 49, 0, time.UTC)},
		},
		{
			Input:    "mysiteBackups/1x2024-12-31_23-59-49/",
			Expected: core.Version{ID: 1, CreatedAt: time.Date(2024, 12, 31, 23, 59, 49, 0, time.UTC)},
		},
		{
			Input:    "mysiteBackups/1234x2024-12-31_23-59-49",
			Expected: core.Version{ID: 1234, CreatedAt: time.Date(2024, 12, 31, 23, 59, 49, 0, time.UTC)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Input, func(t *testing.T) {
			v, err := infrautils.ParseVersionDirName(tc.Input)
			assert.NoError(t, err)
			assert.Equal(t, tc.Expected, *v)
		})
	}
}

func Test_FormatVersionDirName(t *testing.T) {
	testCases := []struct {
		Input    core.Version
		Expected string
	}{
		{
			Input:    core.Version{ID: 1, CreatedAt: time.Date(2024, 12, 31, 23, 59, 49, 0, time.UTC)},
			Expected: "1x2024-12-31_23-59-49",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Expected, func(t *testing.T) {
			s := infrautils.FormatVersionDirName(&tc.Input)
			assert.Equal(t, tc.Expected, s)
		})
	}
}

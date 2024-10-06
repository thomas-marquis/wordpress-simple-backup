package infrautils

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/thomas-marquis/wordpress-simple-backup/internal/core"
)

const (
	timeFormat = "2006-01-02_15-04-05"
)

func ParseVersionDirName(key string) (*core.Version, error) {
	key = strings.TrimSuffix(key, "/")
	parts := strings.Split(key, "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid version dir name")
	}
	versionPart := parts[len(parts)-1]
	versionParts := strings.Split(versionPart, "x")
	if len(versionParts) != 2 {
		return nil, errors.New("invalid version dir name")
	}
	id, err := strconv.Atoi(versionParts[0])
	if err != nil {
		return nil, err
	}

	createdAt, err := time.Parse(timeFormat, versionParts[1])
	if err != nil {
		return nil, err
	}

	return &core.Version{
		ID:        core.VersionID(id),
		CreatedAt: createdAt,
	}, nil
}

func FormatVersionDirName(v *core.Version) string {
	return v.ID.String() + "x" + v.CreatedAt.Format(timeFormat)
}

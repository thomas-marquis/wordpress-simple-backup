package core

import (
	"strconv"
	"time"
)

type VersionID int

func (v VersionID) String() string {
	return strconv.Itoa(int(v))
}

type Version struct {
	ID        VersionID
	CreatedAt time.Time

	dbDump      DumpFile
	contentDump DumpFile
}

func NewVersion(id VersionID, db, content DumpFile) *Version {
	return &Version{
		ID:          id,
		CreatedAt:   time.Now(),
		dbDump:      db,
		contentDump: content,
	}
}

func (v *Version) Age() time.Duration {
	return time.Since(v.CreatedAt)
}

func (v *Version) DbBackupFile() DumpFile {
	return v.dbDump
}

func (v *Version) ContentBackupFile() DumpFile {
	return v.contentDump
}

func (v *Version) AreLocalBackupsAllReady() bool {
	return v.dbDump.Exists() && v.contentDump.Exists()
}

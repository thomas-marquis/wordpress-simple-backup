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

func (v *Version) RemoveDumps() error {
	err := v.dbDump.Remove()
	if err != nil {
		return err
	}
	err = v.contentDump.Remove()
	if err != nil {
		return err
	}
	v.dbDump = DumpFile{}
	v.contentDump = DumpFile{}
	return nil
}

func (v *Version) String() string {
	return "Version " + v.ID.String() + " created at " + v.CreatedAt.Format(time.RFC3339)
}

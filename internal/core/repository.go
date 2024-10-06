package core

type Repository interface {
	CreateContentDump() (DumpFile, error)
	CreateDbDump() (DumpFile, error)
	ClearDump(d DumpFile) error
	ListVersions() ([]*Version, error)
	SaveVersion(v *Version) error
	RemoveVersion(vid VersionID) error
}

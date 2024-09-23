package core

type Backup struct {
	Versions       []Version
	Name           string
	VersionsToKeep int
}

func NewBackup(name string, versionsToKeep int) *Backup {
	return &Backup{
		Name:           name,
		VersionsToKeep: versionsToKeep,
	}
}

func (b *Backup) LastVersion() *Version {
	if len(b.Versions) == 0 {
		return nil
	}
	return &b.Versions[len(b.Versions)-1]
}

func (b *Backup) AddVersion(v Version) {
	b.Versions = append(b.Versions, v)
	if len(b.Versions) > b.VersionsToKeep {
		b.Versions = b.Versions[1:]
	}
}

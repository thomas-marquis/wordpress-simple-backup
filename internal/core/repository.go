package core

type Repository interface {
	GetExistingBackup(name string) (*Backup, error)

	CreateNewBackupVersion(name string) (*Version, error)

	SaveBackup(b *Backup) error

	RestoreToVersion(b *Backup, versionID int) error
}

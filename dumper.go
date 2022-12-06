package repbak

// Dumper defines an interface for backing up a database.
type Dumper interface {
	// Dump does a backup of the database
	Dump() Stat

	// Stop stops the database backup if one is running
	Stop()
}

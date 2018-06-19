package migration

// Migrator interface for database migrations
type Migrator interface {
	Migrate() error
}

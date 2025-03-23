package repository

// RepositoryError represents a repository-specific error
type RepositoryError string

// Error returns the error message
func (e RepositoryError) Error() string {
	return string(e)
}

// Common repository errors
var (
	// ErrMissingMySQLConfig is returned when MySQL configuration is missing
	ErrMissingMySQLConfig = RepositoryError("MySQL configuration is missing")

	// ErrMissingPostgreSQLConfig is returned when PostgreSQL configuration is missing
	ErrMissingPostgreSQLConfig = RepositoryError("PostgreSQL configuration is missing")

	// ErrMissingRedisConfig is returned when Redis configuration is missing
	ErrMissingRedisConfig = RepositoryError("Redis configuration is missing")

	// ErrInvalidTransaction is returned when attempting to use an invalid transaction
	ErrInvalidTransaction = RepositoryError("invalid transaction")

	// ErrInvalidSession is returned when attempting to use an invalid session
	ErrInvalidSession = RepositoryError("invalid session")

	// ErrUnsupportedStoreType is returned when using an unsupported store type
	ErrUnsupportedStoreType = RepositoryError("unsupported store type")
)

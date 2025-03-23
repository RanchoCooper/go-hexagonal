package repo

// Common errors for domain repositories
type RepoError string

func (e RepoError) Error() string {
	return string(e)
}

// Common repository errors
var (
	// ErrNotFound is returned when a requested entity is not found
	ErrNotFound = RepoError("entity not found")
)

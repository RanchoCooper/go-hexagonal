// Package mysql provides MySQL database implementation for repositories
package mysql

import (
	"go-hexagonal/adapter/repository/mysql/entity"
)

// NewExampleRepo creates a new instance of Example repository
func NewExampleRepo() *entity.EntityExample {
	return entity.NewExample()
}

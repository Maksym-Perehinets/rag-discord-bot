package database

import "gorm.io/gorm"

type Service interface {
	// Health checks the health of the database connection by pinging it.
	// It returns a map with keys indicating various health statistics.
	Health() map[string]string
	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close()

	// DB exposes the GORM DB instance for application use.
	DB() *gorm.DB
}

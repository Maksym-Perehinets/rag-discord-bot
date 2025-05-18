package main

import "github.com/Maksym-Perehinets/shared/database"

func main() {
	dbService := database.New()

	// Initialize the database connection
	defer dbService.Close()

	// Perform health check
	health := dbService.Health()
	// Print health status
	for key, value := range health {
		println(key, ":", value)
	}
}

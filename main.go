package main

import (
	"fmt"
	"os"
	"quiz3-go/database"
	"quiz3-go/routers"
)

func main() {
	// Inisialisasi database
	database.InitDB()
	database.Migrate()

	// Setup router
	r := routers.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("🚀 Server running on port " + port)
	r.Run(":" + port)
}

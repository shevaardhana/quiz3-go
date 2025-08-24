package main

import (
	"quiz3-go/database"
	"quiz3-go/routers"
)

func main() {
	database.InitDB()
	r := routers.SetupRouter()
	r.Run(":8080")
}
